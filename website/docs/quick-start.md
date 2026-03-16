---
sidebar_position: 3
title: Quick Start
---

# Quick Start

This guide walks through common jwx workflows to get you productive in minutes.

## Decoding tokens

The most common use case is decoding a JWT to see what's inside. jwx makes this as simple as possible.

### Paste a token directly

Just pass the token as an argument -- no subcommand needed:

```bash
jwx eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

jwx displays the token in colorized, rounded boxes showing:
- **Header** -- the algorithm and token type
- **Payload** -- all claims with timestamps converted to human-readable format
- **Signature** -- the raw signature bytes

### Pipe from other commands

jwx reads from stdin, so you can pipe tokens from any source:

```bash
# From clipboard (macOS)
pbpaste | jwx

# From clipboard (Linux with xclip)
xclip -selection clipboard -o | jwx

# From a file
cat token.txt | jwx

# From an API response
curl -s https://api.example.com/auth | jq -r .token | jwx
```

### Read from clipboard

Use the `--clipboard` flag to read directly from your system clipboard:

```bash
jwx --clipboard
```

This is equivalent to `pbpaste | jwx` on macOS but works cross-platform.

### Use the decode subcommand

The `decode` subcommand is explicit but behaves identically:

```bash
jwx decode eyJhbGciOiJIUzI1NiIs...
jwx decode --clipboard
echo "eyJhbGci..." | jwx decode
```

## Signing tokens

Create signed JWT tokens for testing and development.

### HMAC (shared secret)

The simplest signing method -- uses a shared secret key:

```bash
# Minimal example
jwx sign --alg HS256 --secret mykey '{"sub":"1234"}'

# With more claims
jwx sign --alg HS256 --secret mykey '{"sub":"user_42","name":"Mani","role":"admin","iat":1710600000}'
```

### RSA (key pair)

Use a PEM-encoded private key for RSA signing:

```bash
# Generate a key pair (if you don't have one)
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -o public.pem

# Sign with the private key
jwx sign --alg RS256 --key private.pem '{"sub":"1234","role":"admin"}'
```

### ECDSA

```bash
# Generate an EC key pair
openssl ecparam -genkey -name prime256v1 -noout -out ec-private.pem

# Sign
jwx sign --alg ES256 --key ec-private.pem '{"sub":"1234"}'
```

### Claims from a file

For complex payloads, put your claims in a JSON file:

```bash
# Create a claims file
echo '{"sub":"user_42","name":"Mani","roles":["admin","editor"],"iat":1710600000,"exp":1900000000}' > claims.json

# Sign from the file
jwx sign --alg HS256 --secret mykey --from claims.json
```

## JSON output for scripting

Use `--json` to get machine-readable output that works with `jq` and other tools:

```bash
# Decode as JSON
jwx --json eyJhbGciOiJIUzI1NiIs... | jq .

# Extract a specific claim
jwx --json eyJhbGciOiJIUzI1NiIs... | jq -r .payload.sub

# Check if a token is expired
jwx --json eyJhbGciOiJIUzI1NiIs... | jq .payload.exp

# Sign and output as JSON
jwx sign --alg HS256 --secret mykey --json '{"sub":"1234"}' | jq -r .token
```

## Decode-then-sign workflow

A common development workflow is to decode a token, modify claims, and re-sign:

```bash
# 1. Decode to see current claims
jwx --json eyJhbGciOiJIUzI1NiIs... | jq .payload

# 2. Sign a modified version
jwx sign --alg HS256 --secret dev-secret '{"sub":"1234","role":"admin","exp":1900000000}'

# 3. Verify the new token
jwx $(jwx sign --alg HS256 --secret dev-secret '{"sub":"1234","role":"admin"}')
```

## Handling expired tokens

jwx detects expired tokens automatically:

```bash
# An expired token still decodes, but shows a warning
jwx eyJhbGciOiJIUzI1NiIs...
# Output includes: EXPIRED 2 hours ago

# The exit code is 2 for expired tokens
jwx eyJhbGciOiJIUzI1NiIs...
echo $?  # prints 2
```

This makes it easy to check token validity in scripts:

```bash
if jwx "$TOKEN" > /dev/null 2>&1; then
  echo "Token is valid"
else
  code=$?
  if [ "$code" -eq 2 ]; then
    echo "Token is expired"
  else
    echo "Token is invalid"
  fi
fi
```

## Disabling color output

For environments that don't support ANSI colors:

```bash
# Using the flag
jwx --no-color eyJhbGciOiJIUzI1NiIs...

# Using the environment variable
NO_COLOR=1 jwx eyJhbGciOiJIUzI1NiIs...
```

## Next steps

- [Decode command reference](./cli/decode.md) -- all decode options in detail
- [Sign command reference](./cli/sign.md) -- all signing algorithms and options
- [Supported algorithms](./algorithms.md) -- when to use which algorithm
- [Exit codes](./exit-codes.md) -- for scripting and automation
