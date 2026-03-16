package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/manimovassagh/jwx/internal/display"
)

var noColor bool

var rootCmd = &cobra.Command{
	Use:   "jwx [token]",
	Short: "A beautiful CLI for working with JWTs",
	Long: `jwx is a fast, beautiful command-line tool for decoding, signing,
verifying, and auditing JSON Web Tokens.

Just paste a token and go:
  jwx eyJhbGciOiJIUzI1NiIs...`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if noColor || os.Getenv("NO_COLOR") != "" {
			display.NoColor = true
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if clipboardFlag {
			return runDecode(cmd, args)
		}
		if len(args) > 0 && looksLikeJWT(args[0]) {
			return runDecode(cmd, args)
		}
		return cmd.Help()
	},
}

func looksLikeJWT(s string) bool {
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "eyJ") && strings.Count(s, ".") == 2
}

func Execute() error {
	// Intercept cobra's unknown command error:
	// If the "unknown command" looks like a JWT, decode it.
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	err := rootCmd.Execute()
	if err != nil {
		// Check if the error is about an unknown command that's actually a JWT
		errMsg := err.Error()
		if strings.Contains(errMsg, "unknown command") {
			// Extract the token from os.Args
			for _, arg := range os.Args[1:] {
				if looksLikeJWT(arg) {
					return runDecode(rootCmd, []string{arg})
				}
			}
		}
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(decodeCmd)
	rootCmd.AddCommand(signCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output as JSON")
	rootCmd.Flags().BoolVarP(&clipboardFlag, "clipboard", "c", false, "Read JWT from system clipboard")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output (also respects NO_COLOR env var)")
}
