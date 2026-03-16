<p align="center">
  <h1 align="center">jwx</h1>
  <p align="center"><strong>Decode JWTs instantly. In your terminal or browser — privacy first.</strong></p>
  <p align="center">
    <a href="https://github.com/manimovassagh/jwx/actions/workflows/ci.yml"><img src="https://github.com/manimovassagh/jwx/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
    <a href="https://goreportcard.com/report/github.com/manimovassagh/jwx"><img src="https://goreportcard.com/badge/github.com/manimovassagh/jwx" alt="Go Report Card"></a>
    <a href="https://github.com/manimovassagh/jwx/releases"><img src="https://img.shields.io/github/v/release/manimovassagh/jwx?include_prereleases" alt="Release"></a>
    <a href="LICENSE"><img src="https://img.shields.io/github/license/manimovassagh/jwx" alt="License"></a>
    <a href="https://manimovassagh.github.io/jwx/"><img src="https://img.shields.io/badge/try_in_browser-blue?style=for-the-badge" alt="Try in browser"></a>
  </p>
</p>

<p align="center">
  <img src="docs/assets/jwx-demo.gif" alt="jwx demo" width="800">
</p>

## Why jwx?

- **Just paste and go** -- no subcommands needed for the common case. Run `jwx <token>` and you're done.
- **Beautiful output** -- colorized rounded boxes with human-readable timestamps, not a raw JSON dump.
- **Privacy first** -- tokens are decoded locally on your machine. Nothing is ever sent to a server.
- **Pipe-friendly** -- reads from stdin, supports `--json` for scripting with `jq`.

> **Unlike jwt.io**, your tokens never leave your machine. No tracking, no accounts, 100% local.

## Install

### :apple: macOS

```bash
brew install manimovassagh/tap/jwx
```

### :penguin: Linux

```bash
curl -sL https://github.com/manimovassagh/jwx/releases/latest/download/jwx_linux_amd64.tar.gz | tar xz
sudo mv jwx /usr/local/bin/
```

### Windows

```powershell
Invoke-WebRequest -Uri "https://github.com/manimovassagh/jwx/releases/latest/download/jwx_windows_amd64.zip" -OutFile "$env:TEMP\jwx.zip"
Expand-Archive "$env:TEMP\jwx.zip" -DestinationPath "$env:LOCALAPPDATA\jwx" -Force
$env:PATH += ";$env:LOCALAPPDATA\jwx"
```

Or download `.exe` from [Releases](https://github.com/manimovassagh/jwx/releases/latest) and add to your `PATH`.

### :wrench: From Source

```bash
go install github.com/manimovassagh/jwx/cmd/jwx@latest
```

Or clone and build:

```bash
git clone https://github.com/manimovassagh/jwx.git && cd jwx
make build && make install
```

## Quick Start

```bash
# The simplest way -- just paste a token
jwx eyJhbGciOiJIUzI1NiIs...

# Pipe from clipboard or another command
pbpaste | jwx

# Sign a new token
jwx sign --alg HS256 --secret mykey '{"sub":"1234","role":"admin"}'

# Machine-readable JSON for scripts
jwx --json eyJhbGci... | jq .payload.sub
```

## Features

| Feature | Description |
|---------|-------------|
| **Decode** | Beautiful colorized output with rounded boxes |
| **Sign** | Create tokens with HMAC, RSA, ECDSA, EdDSA |
| **Timestamps** | Auto-converts to human-readable + relative time |
| **Expiry Detection** | `EXPIRED 2 hours ago` / `Expires in 3 days` |
| **JSON Mode** | `--json` flag for piping to `jq` and scripts |
| **Stdin** | Pipe tokens via stdin (`pbpaste \| jwx`) |
| **No Color** | Respects `NO_COLOR` env var and `--no-color` flag |

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

## Shell Completions

```bash
jwx completion bash > /etc/bash_completion.d/jwx   # Bash
jwx completion zsh > "${fpath[1]}/_jwx"             # Zsh
jwx completion fish > ~/.config/fish/completions/jwx.fish  # Fish
jwx completion powershell | Out-String | Invoke-Expression # PowerShell
```

## Roadmap

- [x] `jwx decode` -- beautiful token decoding
- [x] `jwx sign` -- create and sign tokens
- [x] Clipboard support (`jwx --clipboard`)
- [ ] `jwx verify` -- verify token signatures (JWKS support)
- [ ] `jwx inspect` -- deep token analysis
- [ ] `jwx audit` -- security vulnerability checks
- [ ] `jwx keygen` -- generate key pairs

## Contributing

Contributions are welcome! The [Roadmap](#roadmap) above lists planned features — pick one and open a PR. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE)
