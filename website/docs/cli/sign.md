---
sidebar_position: 2
title: "Sign"
---

# jwx sign

Create and sign a JWT token from JSON claims.

## Synopsis

```bash
jwx sign --alg <algorithm> --secret <key> '<claims-json>'
jwx sign --alg <algorithm> --key <keyfile> '<claims-json>'
jwx sign --alg <algorithm> --secret <key> --from <file>
```

## Description

The `sign` command creates a new JWT by signing JSON claims with the specified algorithm and key. It supports HMAC (symmetric), RSA, ECDSA, and EdDSA (asymmetric) algorithms.

The resulting token is printed to stdout, making it easy to pipe into other commands or capture in a variable.

## Flags

| Flag | Description | Required |
|------|-------------|----------|
| `--alg` | Signing algorithm (e.g., `HS256`, `RS256`, `ES256`, `EdDSA`) | Yes |
| `--secret` | Secret key string for HMAC algorithms | For HMAC |
| `--key` | Path to a PEM-encoded private key file for RSA, ECDSA, or EdDSA | For asymmetric |
| `--from` | Read claims from a JSON file instead of the command-line argument | No |
| `--json`, `-j` | Output as JSON object (`{"token":"..."}`) instead of raw token | No |

## Claims input

Claims can be provided in two ways:

### Inline JSON argument

Pass the claims as a JSON string:

```bash
jwx sign --alg HS256 --secret mykey '{"sub":"1234","name":"John Doe","role":"admin"}'
```

:::tip
Use single quotes around the JSON to prevent your shell from interpreting special characters. On Windows PowerShell, you may need to escape differently or use the `--from` flag.
:::

### From a file

For complex claims or to avoid shell escaping issues, write your claims to a file:

```bash
# Create a claims file
cat > claims.json << 'EOF'
{
  "sub": "user_42",
  "name": "Mani Movassagh",
  "roles": ["admin", "editor"],
  "permissions": {
    "read": true,
    "write": true,
    "delete": false
  },
  "iat": 1710600000,
  "exp": 1900000000
}
EOF

# Sign from the file
jwx sign --alg HS256 --secret mykey --from claims.json
```

## Algorithm-specific usage

### HMAC (HS256, HS384, HS512)

HMAC algorithms use a shared secret. Both the signer and verifier must know the same secret.

```bash
# HS256 (most common)
jwx sign --alg HS256 --secret "my-secret-key" '{"sub":"1234"}'

# HS384
jwx sign --alg HS384 --secret "my-secret-key" '{"sub":"1234"}'

# HS512
jwx sign --alg HS512 --secret "my-secret-key" '{"sub":"1234"}'
```

:::caution
Keep your secret keys secure. Never commit them to version control or pass them in URLs. For production use, consider using environment variables:

```bash
jwx sign --alg HS256 --secret "$JWT_SECRET" '{"sub":"1234"}'
```
:::

### RSA (RS256, RS384, RS512)

RSA algorithms use a private key for signing and a public key for verification.

```bash
# Generate a key pair
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -out public.pem

# Sign with the private key
jwx sign --alg RS256 --key private.pem '{"sub":"1234","iss":"https://auth.example.com"}'
```

The `--key` flag accepts any PEM-encoded RSA private key file.

### ECDSA (ES256, ES384, ES512)

ECDSA algorithms use elliptic curve keys, which are smaller and faster than RSA.

```bash
# ES256 (P-256 curve)
openssl ecparam -genkey -name prime256v1 -noout -out ec256-private.pem
jwx sign --alg ES256 --key ec256-private.pem '{"sub":"1234"}'

# ES384 (P-384 curve)
openssl ecparam -genkey -name secp384r1 -noout -out ec384-private.pem
jwx sign --alg ES384 --key ec384-private.pem '{"sub":"1234"}'

# ES512 (P-521 curve)
openssl ecparam -genkey -name secp521r1 -noout -out ec512-private.pem
jwx sign --alg ES512 --key ec512-private.pem '{"sub":"1234"}'
```

### EdDSA (Ed25519)

EdDSA with Ed25519 offers high security with small key sizes and fast operations.

```bash
# Generate an Ed25519 key
openssl genpkey -algorithm Ed25519 -out ed25519-private.pem

# Sign
jwx sign --alg EdDSA --key ed25519-private.pem '{"sub":"1234"}'
```

## Output formats

### Default (raw token)

By default, the signed token is printed as a raw string:

```bash
$ jwx sign --alg HS256 --secret mykey '{"sub":"1234"}'
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0In0.xxxxx
```

This is ideal for capturing in a variable:

```bash
TOKEN=$(jwx sign --alg HS256 --secret mykey '{"sub":"1234"}')
echo "Bearer $TOKEN"
```

### JSON mode

With `--json`, the output is wrapped in a JSON object:

```bash
$ jwx sign --alg HS256 --secret mykey --json '{"sub":"1234"}'
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0In0.xxxxx"}
```

Useful for APIs and scripts:

```bash
jwx sign --alg HS256 --secret mykey --json '{"sub":"1234"}' | jq -r .token
```

## Common workflows

### Create a test token with expiration

```bash
# Token that expires in 1 hour (3600 seconds from now)
EXP=$(( $(date +%s) + 3600 ))
jwx sign --alg HS256 --secret dev-secret "{\"sub\":\"test\",\"exp\":$EXP}"
```

### Round-trip: sign then decode

```bash
# Sign a token and immediately decode it to verify
TOKEN=$(jwx sign --alg HS256 --secret mykey '{"sub":"1234","role":"admin"}')
jwx "$TOKEN"
```

### Use in curl requests

```bash
TOKEN=$(jwx sign --alg HS256 --secret mykey '{"sub":"1234"}')
curl -H "Authorization: Bearer $TOKEN" https://api.example.com/protected
```

## Error handling

The sign command exits with code 1 if:

- No algorithm is specified (`--alg` is required)
- No claims are provided (neither argument nor `--from`)
- The claims string is not valid JSON
- The key file cannot be read or is not a valid PEM key
- The algorithm is not supported
- The key type doesn't match the algorithm (e.g., RSA key with HS256)
