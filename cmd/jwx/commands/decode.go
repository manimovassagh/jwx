package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/manimovassagh/jwx/internal/clipboard"
	"github.com/manimovassagh/jwx/internal/display"
	"github.com/manimovassagh/jwx/internal/jwt"
)

var jsonOutput bool
var clipboardFlag bool

var decodeCmd = &cobra.Command{
	Use:   "decode [token]",
	Short: "Decode and display a JWT token",
	Long: `Decode a JWT token and display its header, payload, and signature
with beautiful colorized output.

Pass the token as an argument, read from clipboard, or pipe via stdin:
  jwx decode eyJhbGci...
  jwx decode --clipboard
  echo "eyJhbGci..." | jwx decode`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDecode,
}

func init() {
	decodeCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output as JSON (for piping)")
}

func runDecode(cmd *cobra.Command, args []string) error {
	var tokenStr string

	switch {
	case clipboardFlag:
		text, err := clipboard.Read()
		if err != nil {
			return fmt.Errorf("failed to read from clipboard: %w", err)
		}
		tokenStr = text
	case len(args) > 0:
		tokenStr = args[0]
	default:
		if isatty.IsTerminal(os.Stdin.Fd()) && !isatty.IsCygwinTerminal(os.Stdin.Fd()) {
			return fmt.Errorf("no token provided\n\nUsage: jwx decode <token>\n       jwx decode --clipboard\n       echo <token> | jwx decode")
		}
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			tokenStr = strings.TrimSpace(scanner.Text())
		}
		if tokenStr == "" {
			return fmt.Errorf("no token provided via stdin")
		}
	}

	token, err := jwt.Decode(tokenStr)
	if err != nil {
		return fmt.Errorf("failed to decode token: %w", err)
	}

	if jsonOutput {
		out, err := display.RenderJSON(token)
		if err != nil {
			return fmt.Errorf("failed to render JSON: %w", err)
		}
		fmt.Println(out)
	} else {
		fmt.Println(display.Render(token))
	}

	// Return exit code 2 if expired (output already shown above)
	if token.IsExpired {
		return &ExitError{Code: 2, Err: fmt.Errorf("token is expired")}
	}

	return nil
}
