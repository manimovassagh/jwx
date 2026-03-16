# Why Your JWT Decoder Might Be Leaking Your Tokens

*I built jwx — a CLI and browser tool that keeps your tokens local*

---

Every developer has done it — copy a JWT from a log or API response, paste it into jwt.io, and decode it. Quick and easy.

But have you thought about what happens to that token?

Browser-based JWT decoders can log your tokens server-side. That token might contain user IDs, emails, roles, permissions — sensitive data you just handed to a third party.

---

## The Problem

Most JWT tools send your token to a server for decoding or signature verification. Even if they claim "client-side only," you're trusting their code and infrastructure.

If you're working with production tokens, staging tokens with real user data, or tokens from security audits — that's a risk.

---

## The Solution: Decode Locally

I built [jwx](https://github.com/manimovassagh/jwx) — a JWT tool that runs entirely on your machine.

---

## In Your Terminal

Just paste and go:

    jwx eyJhbGciOiJIUzI1NiIs...

You get colorized output with human-readable timestamps, expiry detection, and algorithm info. No subcommands, no configuration.

Need it in a script? Pipe it:

    pbpaste | jwx --json | jq .payload.sub

---

## In Your Browser

Don't want to install anything? The [web decoder](https://manimovassagh.github.io/jwx/) runs entirely client-side. Zero network requests. Open DevTools and check — nothing leaves your browser.

It includes:

**Color-coded token parts** — header in blue, payload in green, signature in gold. You can immediately see the three parts of your JWT, just like jwt.io.

**Hover tooltips** — hover over any claim name like "iss", "sub", or "exp" and get a plain-English explanation of what it means.

**Security linter** — automatically flags issues like missing expiry, weak algorithms, or overly long token lifetimes. No button to click — it just runs.

**Signature verification** — verify HMAC, RSA, and ECDSA signatures right in your browser using the Web Crypto API. Your key never leaves the page.

---

## What Can jwx Do?

**Decode** — Beautiful colorized terminal output with rounded boxes and human-readable timestamps.

**Sign** — Create tokens with HMAC (HS256/384/512), RSA (RS256/384/512), ECDSA (ES256/384/512), and EdDSA.

**Security audit** — Automatic checks for common JWT mistakes.

**Verify signatures** — Locally in the browser using Web Crypto API.

**JSON mode** — Machine-readable output for piping to jq and scripts.

**Clipboard** — Read tokens directly from your system clipboard.

---

## Install

**macOS:**

    brew install manimovassagh/tap/jwx

**Windows:**

    choco install jwx

**Linux (Debian/Ubuntu):**

    echo "deb [trusted=yes] https://raw.githubusercontent.com/manimovassagh/jwx/apt-repo/ /" | sudo tee /etc/apt/sources.list.d/jwx.list
    sudo apt update && sudo apt install jwx

**Go:**

    go install github.com/manimovassagh/jwx/cmd/jwx@latest

Or just [try the browser version](https://manimovassagh.github.io/jwx/) — no install needed.

---

## The Security Linter

One feature I'm particularly proud of is the built-in security linter. When you decode a token, jwx automatically checks for:

**Algorithm "none"** — tokens with no signature that anyone can forge.

**Missing expiry** — tokens that never expire give attackers unlimited access.

**Expired tokens** — catches tokens that should have been rejected.

**Missing issuer/audience** — tokens that can't be properly validated.

**Missing JWT ID** — no replay protection.

**Overly long expiry** — tokens valid for more than a year.

These checks run automatically — no configuration needed. In the web decoder, findings appear in a collapsible panel with severity levels from critical to info.

---

## Why Go?

Go was a natural choice for a CLI tool:

**Single binary** — no runtime dependencies, just download and run.

**Cross-platform** — builds for macOS, Linux, and Windows from one codebase.

**Fast startup** — no JVM warmup, no interpreter overhead.

**Standard library crypto** — built-in support for all the algorithms JWTs use.

---

## Open Source

jwx is MIT licensed and fully open source.

**GitHub:** [github.com/manimovassagh/jwx](https://github.com/manimovassagh/jwx)

**Web decoder:** [manimovassagh.github.io/jwx](https://manimovassagh.github.io/jwx/)

**Documentation:** [manimovassagh.github.io/jwx/docs](https://manimovassagh.github.io/jwx/docs/)

---

If you work with JWTs regularly, give it a try. Your tokens will thank you.
