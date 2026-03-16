package commands

import (
	"bytes"
	"testing"
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
