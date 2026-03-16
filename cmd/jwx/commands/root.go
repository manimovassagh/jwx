package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jwx",
	Short: "A beautiful CLI for working with JWTs",
	Long:  `jwx is a fast, beautiful command-line tool for decoding, signing, verifying, and auditing JSON Web Tokens.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(decodeCmd)
	rootCmd.AddCommand(versionCmd)
}
