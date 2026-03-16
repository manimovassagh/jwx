package clipboard

import (
	"fmt"
	"runtime"
	"testing"
)

func TestCommandFunc_Darwin(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("skipping: test only runs on macOS")
	}
	name, args, err := CommandFunc()
	if err != nil {
		t.Fatalf("CommandFunc() unexpected error: %v", err)
	}
	if name != "pbpaste" {
		t.Errorf("expected command name %q, got %q", "pbpaste", name)
	}
	if args != nil {
		t.Errorf("expected nil args on darwin, got %v", args)
	}
}

func TestCommandFunc_ReturnsValidCommandName(t *testing.T) {
	name, _, err := CommandFunc()
	if runtime.GOOS == "windows" {
		// Unsupported, expect error
		if err == nil {
			t.Fatal("expected error on unsupported OS")
		}
		return
	}
	// On darwin/linux we expect a valid command name
	if err != nil {
		t.Skipf("skipping: clipboard tool not available: %v", err)
	}
	if name == "" {
		t.Error("expected non-empty command name")
	}
}

func TestRead_MockSuccess(t *testing.T) {
	original := Read
	t.Cleanup(func() { Read = original })

	Read = func() (string, error) {
		return "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0.signature", nil
	}

	token, err := Read()
	if err != nil {
		t.Fatalf("Read() unexpected error: %v", err)
	}
	if token != "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0.signature" {
		t.Errorf("unexpected token: %s", token)
	}
}

func TestRead_MockError(t *testing.T) {
	original := Read
	t.Cleanup(func() { Read = original })

	Read = func() (string, error) {
		return "", fmt.Errorf("clipboard is empty")
	}

	_, err := Read()
	if err == nil {
		t.Fatal("expected error from mocked Read")
	}
	if err.Error() != "clipboard is empty" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
