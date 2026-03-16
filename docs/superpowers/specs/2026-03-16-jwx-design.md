# jwx — Beautiful JWT CLI Tool + Claude Code Plugin

**Date:** 2026-03-16
**Status:** Approved
**Goal:** Build a star-worthy JWT CLI that developers love, with a Claude Code plugin for extra distribution.

## Overview

jwx is a Go-based CLI tool for decoding, encoding, signing, verifying, and auditing JWTs. It differentiates from competitors (jwt-cli/Rust ~1.4k stars, golang-jwt/cmd, sgaunet/jwt-cli) through stunning terminal output, pipe-friendly design, and a Claude Code plugin.

## Architecture

```
jwx/
├── cmd/jwx/              # Standalone CLI (cobra-based)
│   └── main.go
├── internal/
│   ├── jwt/              # Core JWT engine (decode, encode, sign, verify)
│   ├── display/          # Pretty-printing (lipgloss boxes, colors, badges)
│   ├── keys/             # Key generation & loading (PEM, JWK, JWKS)
│   └── audit/            # Security checks (alg:none, weak keys, expiry)
├── plugin/               # Claude Code plugin wrapper
│   ├── plugin.json
│   ├── commands/
│   │   └── jwt.md        # /jwt slash command
│   ├── skills/
│   │   └── jwt-decode.md # Auto-triggers on JWT tokens
│   └── agents/
│       └── jwt-auditor.md
├── go.mod
├── go.sum
├── Makefile
├── LICENSE
└── README.md
```

### Key Libraries
- **cobra** — CLI framework
- **lipgloss** + **termenv** — terminal styling (from Charm)
- **golang-jwt/jwt/v5** — JWT parsing engine
- **goreleaser** — cross-platform binary releases

### Additional Libraries
- **dustin/go-humanize** — relative time ("2 hours ago") and number formatting
- **net/http** (stdlib) — JWKS endpoint fetching

### Design Principles
- Core logic in `internal/` — both CLI and plugin use same engine
- Plugin is a thin wrapper calling the Go binary via Bash tool
- Pipe-friendly: reads stdin, supports `--json` output, respects `NO_COLOR`
- Single binary, zero runtime dependencies
- Stdin detection: if stdin is a TTY and no token arg, print usage help (don't block)
- No config file — all options via flags. Keep it simple.

## CLI Commands

### v0.1 — The Viral Launch
```
jwx decode <token>          # Colorized header/payload/sig, expiry countdown
jwx decode                  # Reads from stdin (pipe-friendly)
jwx decode --json           # Raw JSON output for scripting
jwx version                 # Version info
```

### v0.2 — Sign (create tokens)
`encode` and `sign` are merged into one command — `sign`. Creating a JWT always involves signing.
```
jwx sign --alg HS256 --secret mykey '{"sub":"1234"}'
jwx sign --alg RS256 --key private.pem '{"sub":"1234"}'
jwx sign --alg RS256 --key private.pem --from payload.json
```

### v0.3 — Verify & Inspect
```
jwx verify <token> --secret mykey
jwx verify <token> --jwks https://example.com/.well-known/jwks.json
jwx inspect <token>         # Deep analysis: algo strength, expiry, claims
```

### v0.4 — Security & Keys
```
jwx audit <token>           # alg:none, weak keys, expired, etc.
jwx keygen --alg RS256      # Generate key pairs
jwx keygen --alg HS256      # Generate symmetric secrets
```

## The "Wow" Output

When you run `jwx decode <token>`:

```
+-- Header ------------------------------------------------+
|  alg: RS256                                               |
|  typ: JWT                                                 |
|  kid: abc123                                              |
+-----------------------------------------------------------+
+-- Payload ------------------------------------------------+
|  sub: "1234567890"                                        |
|  name: "John Doe"                                         |
|  iat: 1516239022  (2018-01-18T01:30:22Z)                  |
|  exp: 1716239022  (2024-05-21T01:30:22Z)                  |
|       ! EXPIRED 1 year ago                                |
+-----------------------------------------------------------+
+-- Signature ----------------------------------------------+
|  Status: Not verified (no key provided)                   |
|  Use: jwx verify <token> --key <path>                     |
+-----------------------------------------------------------+
```

- Header box in blue, payload in green, signature in yellow/red
- Timestamps auto-converted to human-readable + relative time
- Expiry gets a warning badge with color
- Signature status: green checkmark if verified, red X if failed, yellow if not checked

### Output Modes
- **Default:** Pretty boxes with colors (for humans)
- **`--json` / `-j`:** Raw JSON (for piping to jq)
- **`--no-color`:** Plain text (for logging, CI)
- Respects `NO_COLOR` environment variable

## Exit Codes
- `0` — success
- `1` — invalid/malformed token
- `2` — expired token (decode still shows output, exits 2 as warning)
- `3` — signature verification failed
- `4` — key error (file not found, wrong format)
- `5` — network error (JWKS fetch failed)

## Claude Code Plugin

The plugin calls the `jwx` binary via Bash. Users must have `jwx` installed (`go install` or `brew`).

### plugin.json
```json
{
  "name": "jwx",
  "description": "Decode, verify, and audit JWT tokens",
  "version": "0.1.0"
}
```

### Commands
- **`/jwt <token>`** — decode a JWT token (runs `jwx decode <token> --no-color`)
- **`/jwt verify <token> --key <path>`** — verify signature
- **`/jwt audit <token>`** — security analysis

### Skills
- **`jwt-decode`** — description triggers on phrases like "decode this JWT", "what's in this token", or when user pastes a JWT (pattern: `eyJ...`). Runs `jwx decode` via Bash and presents results.

### Agents
- **`jwt-auditor`** — security analysis agent. Given a JWT, it runs `jwx audit` and `jwx decode`, then explains findings in context. Has access to Bash, Read, Grep tools.

## Star-Getting Strategy

1. README with animated GIF/screenshot showing colorized decode output
2. Homebrew tap: `brew install <user>/tap/jwx`
3. `go install github.com/<user>/jwx/cmd/jwx@latest`
4. GitHub releases via goreleaser (Linux, macOS, Windows)
5. Post to: r/golang, r/programming, Hacker News, Twitter/X
6. Claude Code plugin listing
7. Pipe-friendly: `echo $TOKEN | jwx decode`, `jwx decode --json | jq .`

## Makefile Targets
- `make build` — `go build ./cmd/jwx`
- `make test` — `go test ./...`
- `make lint` — `golangci-lint run`
- `make install` — `go install ./cmd/jwx`
- `make release` — goreleaser snapshot

## Testing Strategy

- Unit tests for each `internal/` package
- Table-driven tests for JWT decode/encode edge cases
- Integration tests for CLI commands (run binary, check output)
- Golden file tests for display output (catch visual regressions)
- Test fixtures in `testdata/` directory (sample tokens, keys)
- Golden file update: `go test ./... -update` flag

## Security Considerations

- Never log or display secrets/keys in error messages
- Constant-time comparison for signature verification
- Warn users when tokens use weak algorithms
- Audit command checks for known JWT attack vectors (alg:none, HMAC/RSA confusion)
