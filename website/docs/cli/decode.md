---
sidebar_position: 1
title: "Decode"
---

# jwx decode

Decode and display a JWT token with colorized, human-readable output.

## Synopsis

```bash
jwx [token]
jwx decode [token]
jwx decode --clipboard
echo "token" | jwx decode
```

## Description

The `decode` command parses a JWT token and displays its three parts -- header, payload, and signature -- in a formatted, colorized terminal display.

You can invoke it in two ways:
- **Implicitly**: `jwx <token>` -- if the argument looks like a JWT (starts with `eyJ` and has three dot-separated parts), jwx decodes it automatically without needing the `decode` subcommand.
- **Explicitly**: `jwx decode <token>` -- uses the decode subcommand directly.

## Input methods

### Command-line argument

Pass the token directly as an argument:

```bash
jwx eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

### Standard input (stdin)

Pipe a token from another command or file:

```bash
# From clipboard (macOS)
pbpaste | jwx

# From clipboard (Linux)
xclip -selection clipboard -o | jwx

# From a file
cat token.txt | jwx decode

# From an API response
curl -s https://api.example.com/auth | jq -r .access_token | jwx
```

When stdin is a pipe (not a terminal), jwx reads the first line and treats it as the token.

### Clipboard

Read directly from the system clipboard:

```bash
jwx --clipboard
jwx -c
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | `-j` | Output as JSON instead of the colorized display |
| `--clipboard` | `-c` | Read the token from the system clipboard |
| `--no-color` | | Disable colored output (also respects the `NO_COLOR` environment variable) |

## Output format

### Default (colorized)

By default, jwx renders the token in colorized, rounded boxes:

- **Header section** -- shows the algorithm (`alg`), token type (`typ`), and any additional header claims like `kid`
- **Payload section** -- shows all claims with special formatting for:
  - `exp` (expiration): displayed as a human-readable date with relative time ("Expires in 3 days" or "EXPIRED 2 hours ago")
  - `iat` (issued at): displayed as a human-readable date with relative time
  - `nbf` (not before): displayed with validity status
- **Signature section** -- shows the raw base64url-encoded signature

### JSON mode

With `--json`, jwx outputs a JSON object containing the decoded token:

```bash
jwx --json eyJhbGciOiJIUzI1NiIs... | jq .
```

```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "1234567890",
    "name": "John Doe",
    "iat": 1516239022
  },
  "signature": "SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}
```

This is ideal for scripting:

```bash
# Extract a single claim
jwx --json "$TOKEN" | jq -r .payload.sub

# Check the algorithm
jwx --json "$TOKEN" | jq -r .header.alg

# Get all roles
jwx --json "$TOKEN" | jq '.payload.roles[]'
```

## Token validation

jwx performs basic structural validation:

- The token must have exactly three dot-separated parts
- The first part (header) must start with `eyJ` (base64url-encoded `{`)
- Both header and payload must be valid JSON objects after base64url decoding

jwx does **not** verify signatures during decode. Signature verification is a separate concern (see the roadmap for the upcoming `jwx verify` command).

## Expiry detection

If the payload contains an `exp` claim, jwx automatically checks it against the current time:

- **Valid token**: shows "Expires in X days/hours/minutes"
- **Expired token**: shows "EXPIRED X days/hours/minutes ago" and exits with code 2

The token is still fully decoded and displayed even when expired.

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | Token decoded successfully and is not expired |
| `1` | Invalid or malformed token (could not be decoded) |
| `2` | Token is expired (output is still shown) |

See [Exit Codes](../exit-codes.md) for more details.

## Examples

```bash
# Basic decode
jwx eyJhbGciOiJIUzI1NiIs...

# Decode with JSON output
jwx decode --json eyJhbGciOiJIUzI1NiIs...

# Decode from clipboard
jwx decode --clipboard

# Decode without color
jwx --no-color eyJhbGciOiJIUzI1NiIs...

# Pipe and extract a claim
pbpaste | jwx --json | jq -r .payload.email

# Check expiry in a script
if jwx "$TOKEN" > /dev/null 2>&1; then
  echo "Token is valid"
fi
```
