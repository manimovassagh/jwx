package commands

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/manimovassagh/jwx/internal/clipboard"
)

const testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

func TestDecodeCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "valid token",
			args:    []string{testToken},
			wantErr: false,
		},
		{
			name:    "invalid token",
			args:    []string{"not-a-jwt"},
			wantErr: true,
		},
		{
			name:    "empty args no stdin",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := decodeCmd
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.RunE(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeCmd error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeCmd_JSON(t *testing.T) {
	jsonOutput = true
	defer func() { jsonOutput = false }()

	cmd := decodeCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.RunE(cmd, []string{testToken})
	if err != nil {
		t.Fatalf("decodeCmd --json error = %v", err)
	}
}

func TestSignCmd_MissingAlg(t *testing.T) {
	cmd := signCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	signAlg = ""
	signSecret = "test"
	err := cmd.RunE(cmd, []string{`{"sub":"1234"}`})
	if err == nil {
		t.Error("expected error for empty algorithm")
	}
}

func TestSignCmd_HMAC(t *testing.T) {
	signAlg = "HS256"
	signSecret = "test-secret"
	signKey = ""
	signFrom = ""

	cmd := signCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.RunE(cmd, []string{`{"sub":"1234"}`})
	if err != nil {
		t.Fatalf("signCmd error = %v", err)
	}
}

func TestSignCmd_NoClaims(t *testing.T) {
	signAlg = "HS256"
	signSecret = "test"
	signFrom = ""

	cmd := signCmd
	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Error("expected error for no claims")
	}
}

func TestLooksLikeJWT(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0.signature", true},
		{"  eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0.sig  ", true},
		{"not-a-jwt", false},
		{"eyJ.only-two-parts", false},
		{"", false},
		{"hello.world.test", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := looksLikeJWT(tt.input)
			if got != tt.want {
				t.Errorf("looksLikeJWT(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestVersionCmd(t *testing.T) {
	SetVersion("1.0.0", "abc123", "2026-01-01")

	if buildVersion != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", buildVersion)
	}
	if buildCommit != "abc123" {
		t.Errorf("expected commit abc123, got %s", buildCommit)
	}
}

func TestRootCmd_JWTArg(t *testing.T) {
	cmd := rootCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.RunE(cmd, []string{testToken})
	if err != nil {
		t.Fatalf("rootCmd with JWT arg error = %v", err)
	}
}

func TestDecodeCmd_Clipboard(t *testing.T) {
	// Mock clipboard to return a valid token
	origRead := clipboard.Read
	clipboard.Read = func() (string, error) {
		return testToken, nil
	}
	defer func() { clipboard.Read = origRead }()

	clipboardFlag = true
	defer func() { clipboardFlag = false }()

	err := runDecode(decodeCmd, []string{})
	if err != nil {
		t.Fatalf("decode with clipboard error = %v", err)
	}
}

func TestClipboardFlagExists(t *testing.T) {
	// Verify the --clipboard / -c flag is registered on decodeCmd
	f := decodeCmd.Flags().Lookup("clipboard")
	if f == nil {
		t.Fatal("decode command missing --clipboard flag")
	}
	if f.Shorthand != "c" {
		t.Errorf("expected shorthand 'c', got %q", f.Shorthand)
	}

	// Verify the flag is also on rootCmd
	rf := rootCmd.Flags().Lookup("clipboard")
	if rf == nil {
		t.Fatal("root command missing --clipboard flag")
	}
	if rf.Shorthand != "c" {
		t.Errorf("expected root shorthand 'c', got %q", rf.Shorthand)
	}
}

func TestClipboardCommandFunc(t *testing.T) {
	name, args, err := clipboard.CommandFunc()

	switch runtime.GOOS {
	case "darwin":
		if err != nil {
			t.Fatalf("unexpected error on darwin: %v", err)
		}
		if name != "pbpaste" {
			t.Errorf("expected pbpaste, got %s", name)
		}
		if len(args) != 0 {
			t.Errorf("expected no args for pbpaste, got %v", args)
		}
	case "linux":
		if err != nil {
			t.Logf("no clipboard tool on this linux: %v", err)
		} else if name != "xclip" && name != "xsel" {
			t.Errorf("expected xclip or xsel, got %s", name)
		}
	default:
		if err == nil {
			t.Errorf("expected error on unsupported OS %s", runtime.GOOS)
		}
	}
}

func TestDecodeCmd_ClipboardError(t *testing.T) {
	origRead := clipboard.Read
	defer func() { clipboard.Read = origRead }()

	clipboard.Read = func() (string, error) {
		return "", fmt.Errorf("clipboard unavailable")
	}

	clipboardFlag = true
	defer func() { clipboardFlag = false }()

	err := runDecode(decodeCmd, []string{})
	if err == nil {
		t.Fatal("expected error when clipboard fails")
	}
}

func TestDecodeCmd_ClipboardPriorityOverArg(t *testing.T) {
	origRead := clipboard.Read
	defer func() { clipboard.Read = origRead }()

	clipboard.Read = func() (string, error) {
		return testToken, nil
	}

	clipboardFlag = true
	defer func() { clipboardFlag = false }()

	// Pass an invalid token as arg — clipboard should be used instead
	err := runDecode(decodeCmd, []string{"not-a-jwt"})
	if err != nil {
		t.Fatalf("clipboard should take priority over arg, got error: %v", err)
	}
}

func TestRootCmd_ClipboardMocked(t *testing.T) {
	origRead := clipboard.Read
	defer func() { clipboard.Read = origRead }()

	clipboard.Read = func() (string, error) {
		return testToken, nil
	}

	clipboardFlag = true
	defer func() { clipboardFlag = false }()

	cmd := rootCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Fatalf("rootCmd --clipboard error = %v", err)
	}
}

func TestRootCmd_NonJWTArg(t *testing.T) {
	cmd := rootCmd
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.RunE(cmd, []string{"not-a-jwt"})
	// Should show help (no error)
	if err != nil {
		t.Fatalf("rootCmd with non-JWT arg error = %v", err)
	}
}

func TestExecute(t *testing.T) {
	// Execute with no args should show help and not error
	err := Execute()
	if err != nil {
		t.Fatalf("Execute() with no args error = %v", err)
	}
}

func TestLooksLikeJWT_EdgeCases(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"eyJhbGciOiJIUzI1NiJ9.eyJ0ZXN0IjoiMSJ9.sig", true},
		{"eyJ.a.b", true},
		{"eyJ.a.b.c", false},      // 4 parts
		{"  eyJ.a.b  \n", true},   // whitespace
		{"bearer eyJ.a.b", false}, // bearer prefix breaks it
		{"EYJ.a.b", false},        // wrong case
	}

	for _, tt := range tests {
		got := looksLikeJWT(tt.input)
		if got != tt.want {
			t.Errorf("looksLikeJWT(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestSignCmd_FromFile(t *testing.T) {
	// Create temp claims file
	f, err := os.CreateTemp("", "claims-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(f.Name()) }()

	if _, err := f.WriteString(`{"sub":"file-test"}`); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	signAlg = "HS256"
	signSecret = "test"
	signFrom = f.Name()
	defer func() { signFrom = "" }()

	err = runSign(signCmd, []string{})
	if err != nil {
		t.Fatalf("sign from file error = %v", err)
	}
}
