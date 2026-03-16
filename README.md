# jwx

A beautiful CLI for working with JSON Web Tokens.

```
jwx decode eyJhbGciOiJIUzI1NiIs...
```

```
╭────────────────╮
│ Header         │
│   alg: "HS256" │
│   typ: "JWT"   │
╰────────────────╯
╭───────────────────────────────────────────╮
│ Payload                                   │
│   sub: "1234567890"                       │
│   iat: 1516239022  (2018-01-18T01:30:22Z) │
│   name: "John Doe"                        │
╰───────────────────────────────────────────╯
╭─────────────────────────────────────────────╮
│ Signature                                   │
│   Algorithm: "HS256"                        │
│   Status: Not verified (no key provided)    │
╰─────────────────────────────────────────────╯
```

## Install

```bash
go install github.com/mani-sh-reddy/jwx/cmd/jwx@latest
```

## Usage

```bash
# Decode a JWT
jwx decode <token>

# Pipe from stdin
echo $TOKEN | jwx decode
pbpaste | jwx decode

# JSON output (for piping to jq)
jwx decode --json <token>
jwx decode -j <token> | jq .payload

# Version
jwx version
```

## Features

- Beautiful colorized output with rounded boxes
- Automatic timestamp conversion to human-readable format
- Expiry detection with warnings (`⚠ EXPIRED 2 hours ago`)
- JSON output mode for scripting (`--json`)
- Stdin support for piping
- Respects `NO_COLOR` environment variable
- Single binary, zero runtime dependencies

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Invalid/malformed token |
| 2 | Token is expired (output still shown) |

## Roadmap

- [ ] `jwx sign` — create and sign tokens
- [ ] `jwx verify` — verify token signatures
- [ ] `jwx inspect` — deep token analysis
- [ ] `jwx audit` — security vulnerability checks
- [ ] `jwx keygen` — generate keys
- [ ] Claude Code plugin

## License

MIT
