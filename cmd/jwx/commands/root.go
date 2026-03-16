package commands

import (
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jwx [token]",
	Short: "A beautiful CLI for working with JWTs",
	Long: `jwx is a fast, beautiful command-line tool for decoding, signing,
verifying, and auditing JSON Web Tokens.

Just paste a token and go:
  jwx eyJhbGciOiJIUzI1NiIs...`,
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If the first arg looks like a JWT, decode it
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
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(decodeCmd)
	rootCmd.AddCommand(signCmd)
	rootCmd.AddCommand(versionCmd)

	// Add --json flag to root too so `jwx --json <token>` works
	rootCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output as JSON")
}
