---
sidebar_position: 6
title: Supported Algorithms
---

# Supported Algorithms

jwx supports all major JWT signing algorithms. This page explains each one and when to use it.

## Algorithm overview

| Family | Algorithms | Key Type | Use Case |
|--------|-----------|----------|----------|
| **HMAC** | HS256, HS384, HS512 | Shared secret | Simple services, internal APIs |
| **RSA** | RS256, RS384, RS512 | RSA key pair | Widely supported, most identity providers |
| **ECDSA** | ES256, ES384, ES512 | EC key pair | Smaller keys, mobile and IoT |
| **EdDSA** | Ed25519 | Ed25519 key pair | Modern, fast, small signatures |

## HMAC (HS256, HS384, HS512)

**Symmetric** algorithms -- the same secret is used for both signing and verification.

```bash
jwx sign --alg HS256 --secret "my-secret-key" '{"sub":"1234"}'
```

### How it works

HMAC computes a keyed hash (SHA-256, SHA-384, or SHA-512) over the header and payload. The verifier re-computes the hash with the same secret and compares.

### When to use HMAC

- **Internal services** where both the issuer and verifier share a secret
- **Simple setups** where managing key pairs adds unnecessary complexity
- **Development and testing** where security requirements are lower

### When to avoid HMAC

- **Third-party verification** -- you'd have to share the secret with every verifier, which is a security risk
- **Public APIs** -- anyone with the secret can both create and verify tokens
- **Key rotation** -- rotating a shared secret requires coordinating all parties simultaneously

### Key size recommendations

| Algorithm | Minimum secret length | Recommended |
|-----------|----------------------|-------------|
| HS256 | 32 bytes (256 bits) | 32+ bytes |
| HS384 | 48 bytes (384 bits) | 48+ bytes |
| HS512 | 64 bytes (512 bits) | 64+ bytes |

## RSA (RS256, RS384, RS512)

**Asymmetric** algorithms -- sign with a private key, verify with the corresponding public key.

```bash
# Generate keys
openssl genrsa -out private.pem 2048

# Sign
jwx sign --alg RS256 --key private.pem '{"sub":"1234"}'
```

### How it works

RSA uses the RSASSA-PKCS1-v1_5 signature scheme with SHA-256, SHA-384, or SHA-512. The private key signs the data, and anyone with the public key can verify without being able to forge tokens.

### When to use RSA

- **Identity providers (IdPs)** like Auth0, Okta, and Keycloak -- this is the industry default
- **Public verification** -- distribute the public key via JWKS endpoint; verifiers can't forge tokens
- **Interoperability** -- RS256 is the most widely supported JWT algorithm across libraries and platforms

### Key size recommendations

| Key size | Security level | Status |
|----------|---------------|--------|
| 2048-bit | 112-bit | Minimum recommended |
| 3072-bit | 128-bit | Good for new deployments |
| 4096-bit | ~140-bit | Maximum practical security |

### Trade-offs

- Larger signatures (~256 bytes for 2048-bit keys) compared to ECDSA
- Slower signing and verification than ECDSA or EdDSA
- Universally supported across all JWT libraries

## ECDSA (ES256, ES384, ES512)

**Asymmetric** algorithms using elliptic curve cryptography. Smaller keys and signatures than RSA with equivalent security.

```bash
# Generate keys
openssl ecparam -genkey -name prime256v1 -noout -out ec-private.pem

# Sign
jwx sign --alg ES256 --key ec-private.pem '{"sub":"1234"}'
```

### Algorithm-to-curve mapping

| Algorithm | Curve | Key size | Signature size |
|-----------|-------|----------|---------------|
| ES256 | P-256 (prime256v1) | 256-bit | ~64 bytes |
| ES384 | P-384 (secp384r1) | 384-bit | ~96 bytes |
| ES512 | P-521 (secp521r1) | 521-bit | ~132 bytes |

### When to use ECDSA

- **Mobile and IoT** -- smaller tokens due to compact signatures
- **Performance-sensitive** applications -- faster than RSA for most operations
- **Modern infrastructure** -- well-supported in Go, Node.js, Python, Java 7+

### When to avoid ECDSA

- **Legacy systems** that only support RSA
- **Deterministic signatures needed** -- standard ECDSA uses a random nonce; improper RNG can leak the private key (though most modern libraries handle this safely)

## EdDSA (Ed25519)

**Asymmetric** algorithm using the Edwards-curve Digital Signature Algorithm with Curve25519.

```bash
# Generate keys
openssl genpkey -algorithm Ed25519 -out ed25519-private.pem

# Sign
jwx sign --alg EdDSA --key ed25519-private.pem '{"sub":"1234"}'
```

### When to use EdDSA

- **New projects** with no legacy constraints -- Ed25519 is considered state-of-the-art
- **Performance** -- faster than both RSA and ECDSA for signing and verification
- **Security** -- deterministic signatures (no RNG dependency), resistant to timing attacks
- **Compact** -- 32-byte public keys, 64-byte signatures

### When to avoid EdDSA

- **Broad interoperability** -- not all JWT libraries and platforms support EdDSA yet (though support is growing rapidly)
- **Compliance requirements** -- some standards and regulations specify RSA or ECDSA

### Comparison with other algorithms

| Property | Ed25519 | ES256 | RS256 (2048) |
|----------|---------|-------|-------------|
| Public key size | 32 bytes | 64 bytes | 294 bytes |
| Signature size | 64 bytes | 64 bytes | 256 bytes |
| Sign speed | Fastest | Fast | Slow |
| Verify speed | Fast | Fast | Moderate |
| Deterministic | Yes | No* | N/A |

*Some libraries implement deterministic ECDSA (RFC 6979), but it's not guaranteed.

## Choosing an algorithm

For most new projects, here's a simple decision tree:

1. **Do both parties share a secret?** Use **HS256**.
2. **Do you need broad compatibility?** Use **RS256**.
3. **Do you need small tokens?** Use **ES256**.
4. **Building something new with no legacy constraints?** Use **EdDSA** (Ed25519).

When in doubt, **RS256** is the safest default -- it works everywhere and is well-understood.
