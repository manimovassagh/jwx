package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/manimovassagh/jwx/cmd/jwx/commands"
)

// Set via -ldflags at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	commands.SetVersion(version, commit, date)
	if err := commands.Execute(); err != nil {
		var exitErr *commands.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.Code)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
