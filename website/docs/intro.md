---
slug: /
sidebar_position: 1
title: Welcome to jwx
---

# Welcome to jwx

**jwx** is a fast, beautiful command-line tool for working with JSON Web Tokens. Decode, sign, and inspect JWTs right from your terminal -- or use the [web decoder](https://manimovassagh.github.io/jwx/) in your browser.

## Why jwx?

Working with JWTs during development often means pasting tokens into browser-based tools that may log or transmit your data. **jwx** takes a different approach:

- **Just paste and go** -- no subcommands needed for the common case. Run `jwx <token>` and you're done.
- **Beautiful output** -- colorized rounded boxes with human-readable timestamps, not a raw JSON dump.
- **Privacy first** -- tokens are decoded locally on your machine. Nothing is ever sent to a server.
- **Pipe-friendly** -- reads from stdin, supports `--json` for scripting with `jq`.

## What can jwx do?

| Capability | Description |
|------------|-------------|
| **Decode** | Display JWT header, payload, and signature with colorized, human-readable output |
| **Sign** | Create tokens with HMAC (HS256/384/512), RSA (RS256/384/512), ECDSA (ES256/384/512), and EdDSA |
| **Expiry detection** | Automatically detects expired tokens and shows relative timestamps |
| **JSON mode** | Machine-readable output for piping to `jq` and use in scripts |
| **Clipboard** | Read tokens directly from your system clipboard |
| **Web decoder** | A browser-based decoder that runs entirely client-side |

## Quick example

```bash
# Decode a token -- just paste it
jwx eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

jwx renders a colorized, boxed display showing the header, payload with human-readable timestamps, and signature -- all in your terminal.

## Next steps

- [Install jwx](./installation.md) on your platform
- [Quick Start](./quick-start.md) with common workflows
- [CLI Reference](./cli/decode.md) for all commands and options
- [Web Decoder](./web-decoder.md) to decode tokens in your browser
