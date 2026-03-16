<p align="center">
  <h1 align="center">jwx</h1>
  <p align="center">A beautiful CLI for working with JSON Web Tokens</p>
  <p align="center">
    <a href="https://github.com/manimovassagh/jwx/actions/workflows/ci.yml"><img src="https://github.com/manimovassagh/jwx/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
    <a href="https://goreportcard.com/report/github.com/manimovassagh/jwx"><img src="https://goreportcard.com/badge/github.com/manimovassagh/jwx" alt="Go Report Card"></a>
    <a href="https://github.com/manimovassagh/jwx/releases"><img src="https://img.shields.io/github/v/release/manimovassagh/jwx?include_prereleases" alt="Release"></a>
    <a href="LICENSE"><img src="https://img.shields.io/github/license/manimovassagh/jwx" alt="License"></a>
  </p>
</p>

---

**jwx** decodes, signs, verifies, and audits JWTs — with beautiful, colorized terminal output.

<p align="center">
  <img src="docs/assets/demo.gif" alt="jwx demo" width="800">
</p>

> **Try it in your browser**: [manimovassagh.github.io/jwx](https://manimovassagh.github.io/jwx/)

## Install

```bash
# Go (recommended)
go install github.com/manimovassagh/jwx/cmd/jwx@latest

# Homebrew (macOS/Linux)
brew install manimovassagh/tap/jwx

# Download binary (Linux, macOS, Windows)
# → https://github.com/manimovassagh/jwx/releases/latest
```

### Other platforms

| Platform | Command |
|----------|---------|
| **Go** | `go install github.com/manimovassagh/jwx/cmd/jwx@latest` |
| **Homebrew** | `brew install manimovassagh/tap/jwx` |
| **Binary** | Download from [Releases](https://github.com/manimovassagh/jwx/releases/latest) |
| **Docker** | `docker run --rm ghcr.io/manimovassagh/jwx decode <token>` |
| **From source** | `git clone && make install` |

## Shell Completions

```bash
# Bash
jwx completion bash > /etc/bash_completion.d/jwx

# Zsh
jwx completion zsh > "${fpath[1]}/_jwx"

# Fish
jwx completion fish > ~/.config/fish/completions/jwx.fish

# PowerShell
jwx completion powershell | Out-String | Invoke-Expression
```

## Quick Start

```bash
# Decode a token
jwx decode eyJhbGciOiJIUzI1NiIs...

# Pipe from clipboard (macOS)
pbpaste | jwx decode

# Pipe from environment variable
echo $JWT_TOKEN | jwx decode

# JSON output for scripting
jwx decode --json eyJhbGci... | jq .payload.sub

# Sign a new token
jwx sign --alg HS256 --secret mykey '{"sub":"1234","role":"admin"}'

# Sign with RSA key
jwx sign --alg RS256 --key private.pem '{"sub":"1234"}'

# Sign from a file
jwx sign --alg HS256 --secret mykey --from claims.json
```

## Features

| Feature | Description |
|---------|-------------|
| **Decode** | Beautiful colorized output with rounded boxes |
| **Sign** | Create tokens with HMAC, RSA, ECDSA, EdDSA |
| **Timestamps** | Auto-converts to human-readable + relative time |
| **Expiry Detection** | `⚠ EXPIRED 2 hours ago` / `✓ Expires in 3 days` |
| **JSON Mode** | `--json` flag for piping to `jq` and scripts |
| **Stdin** | Pipe tokens via stdin (`pbpaste \| jwx decode`) |
| **No Color** | Respects `NO_COLOR` env var and `--no-color` flag |
| **Exit Codes** | Scriptable: 0=ok, 1=invalid, 2=expired |

## Supported Algorithms

| Type | Algorithms |
|------|-----------|
| HMAC | HS256, HS384, HS512 |
| RSA | RS256, RS384, RS512 |
| ECDSA | ES256, ES384, ES512 |
| EdDSA | Ed25519 |

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Invalid or malformed token |
| `2` | Token is expired (output still shown) |
| `3` | Signature verification failed |
| `4` | Key error (file not found, wrong format) |

## Roadmap

- [x] `jwx decode` — beautiful token decoding
- [x] `jwx sign` — create and sign tokens
- [ ] `jwx verify` — verify token signatures (JWKS support)
- [ ] `jwx inspect` — deep token analysis
- [ ] `jwx audit` — security vulnerability checks
- [ ] `jwx keygen` — generate key pairs
- [ ] Claude Code plugin

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

```bash
# Clone
git clone https://github.com/manimovassagh/jwx.git
cd jwx

# Build
make build

# Test
make test

# Install locally
make install
```

## License

[MIT](LICENSE)
