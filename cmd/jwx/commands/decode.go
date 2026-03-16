package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/mani-sh-reddy/jwx/internal/display"
	"github.com/mani-sh-reddy/jwx/internal/jwt"
)

var jsonOutput bool

var decodeCmd = &cobra.Command{
	Use:   "decode [token]",
	Short: "Decode and display a JWT token",
	Long: `Decode a JWT token and display its header, payload, and signature
with beautiful colorized output.

Pass the token as an argument or pipe it via stdin:
  jwx decode eyJhbGci...
  echo "eyJhbGci..." | jwx decode
  pbpaste | jwx decode`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDecode,
}

func init() {
	decodeCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output as JSON (for piping)")
}

func runDecode(cmd *cobra.Command, args []string) error {
	var tokenStr string

	if len(args) > 0 {
		tokenStr = args[0]
	} else {
		// Check if stdin has data
		if isatty.IsTerminal(os.Stdin.Fd()) && isatty.IsCygwinTerminal(os.Stdin.Fd()) == false {
			return fmt.Errorf("no token provided\n\nUsage: jwx decode <token>\n       echo <token> | jwx decode")
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

	// Exit code 2 if expired (still show output)
	if token.IsExpired {
		os.Exit(2)
	}

	return nil
}
