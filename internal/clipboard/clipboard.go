package clipboard

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// CommandFunc returns the command name and args for reading the clipboard
// on the current OS. Exported for testing.
func CommandFunc() (name string, args []string, err error) {
	switch runtime.GOOS {
	case "darwin":
		return "pbpaste", nil, nil
	case "linux":
		// Prefer xclip, fall back to xsel
		if _, err := exec.LookPath("xclip"); err == nil {
			return "xclip", []string{"-selection", "clipboard", "-o"}, nil
		}
		if _, err := exec.LookPath("xsel"); err == nil {
			return "xsel", []string{"--clipboard", "--output"}, nil
		}
		return "", nil, fmt.Errorf("no clipboard tool found: install xclip or xsel")
	default:
		return "", nil, fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}
}

// Read reads text from the system clipboard.
var Read = func() (string, error) {
	name, args, err := CommandFunc()
	if err != nil {
		return "", err
	}

	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", fmt.Errorf("failed to read clipboard: %w", err)
	}

	text := strings.TrimSpace(string(out))
	if text == "" {
		return "", fmt.Errorf("clipboard is empty")
	}
	return text, nil
}
