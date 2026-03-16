# Contributing to jwx

Thank you for your interest in contributing to **jwx**. This guide covers everything you need to get started.

## Prerequisites

Before you begin, make sure you have the following installed:

- **Go 1.23+** — verify with `go version`
- **golangci-lint** — install via `brew install golangci-lint` (macOS) or `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- **git** — verify with `git --version`

## Development Setup

```bash
# Clone the repository
git clone https://github.com/manimovassagh/jwx.git
cd jwx

# Install pre-commit hooks (required)
make setup

# Build the binary
make build

# Run tests
make test
```

**Requirements:** Go 1.23+, [golangci-lint](https://golangci-lint.run/welcome/install-local/)

## Pre-commit Hooks

Running `make setup` installs a pre-commit hook that enforces the following on every commit:

1. **gofmt** - all files must be formatted
2. **go vet** - static analysis must pass
3. **golangci-lint** - linter must pass with no issues
4. **Tests + coverage** - all tests must pass with a minimum of **80% coverage**

If any check fails, the commit is rejected with a clear message explaining what to fix.

## Code Style

- Format all code with `gofmt` (or run `make fmt`)
- Run `golangci-lint run` before pushing (or run `make lint`)
- Run `make check` to execute formatting, linting, and tests in one step
- Follow standard Go conventions: [Effective Go](https://go.dev/doc/effective_go), [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)

## Testing Requirements

- **Minimum 80% code coverage** (enforced by the pre-commit hook)
- Run tests with the race detector: `go test -race ./...`
- Place test files alongside the code they test (`*_test.go`)
- Place test fixtures in the `testdata/` directory
- All new features and bug fixes must include tests

## Pull Request Process

1. **One PR at a time.** Do not open a new PR until the previous one is merged.
2. **Always create PRs as drafts first** so maintainers can monitor progress.
3. Wait for CI to pass, then mark the PR as ready for review.
4. Keep PRs focused -- one logical change per PR.
5. Write a clear PR title and description explaining *why* the change is needed.

**Flow:** branch -> commit -> push -> draft PR -> CI passes -> ready for review -> merge

## Architecture Overview

```
jwx/
├── cmd/jwx/              # Entry point
│   ├── main.go           # Main function with version/build info
│   └── commands/         # CLI command definitions (Cobra)
│       ├── root.go       # Root command and global flags
│       ├── decode.go     # jwx decode
│       ├── sign.go       # jwx sign
│       └── version.go    # jwx version
├── internal/
│   ├── jwt/              # Core JWT parsing, signing, and validation logic
│   └── display/          # Terminal output formatting (boxes, colors, JSON)
├── plugin/               # Claude Code plugin definitions
│   ├── commands/         # Plugin command specs
│   └── skills/           # Plugin skill specs
├── hooks/                # Git hooks (pre-commit)
├── testdata/             # Test fixtures (keys, sample tokens)
└── docs/                 # Documentation site
```

**Key packages:**

- `cmd/jwx/commands` -- Defines CLI commands using [Cobra](https://github.com/spf13/cobra). Each command lives in its own file.
- `internal/jwt` -- Pure JWT logic: decoding, signing, verification. No terminal output.
- `internal/display` -- Presentation layer: colorized boxes, JSON formatting, human-readable timestamps.
- `plugin/` -- Claude Code plugin integration (commands and skills).

## How to Add a New Command

1. **Create the command file** in `cmd/jwx/commands/` (e.g., `verify.go`).

2. **Define the Cobra command:**
   ```go
   package commands

   import "github.com/spf13/cobra"

   var verifyCmd = &cobra.Command{
       Use:   "verify <token>",
       Short: "Verify a JWT signature",
       RunE: func(cmd *cobra.Command, args []string) error {
           // Implementation here
           return nil
       },
   }
   ```

3. **Register it** in `root.go` by adding `rootCmd.AddCommand(verifyCmd)` to the `init()` function.

4. **Add core logic** in `internal/jwt/` -- keep command files thin, with business logic in the internal packages.

5. **Add tests** in `cmd/jwx/commands/` and `internal/jwt/` to maintain 80% coverage.

6. **Update the README** with usage examples for the new command.

## Issue Labels

| Label | Description |
|-------|-------------|
| `bug` | Something is broken or behaving incorrectly |
| `feature` | New functionality or enhancement request |
| `security` | Security vulnerability or hardening task |

Security findings should always be filed as GitHub issues with reproduction steps, affected files, and a recommended fix.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
