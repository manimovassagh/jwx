---
sidebar_position: 8
title: Privacy & Security
---

# Privacy & Security

Privacy is a core design principle of jwx. This page explains how jwx handles your data and why it's safe for sensitive tokens.

## CLI: everything stays local

When you use the jwx CLI:

- **All processing happens on your machine.** Token decoding, signing, and formatting are performed entirely in your local process.
- **No network requests.** jwx makes zero network connections. It doesn't phone home, check for updates, or transmit telemetry.
- **No disk writes.** Decoded tokens are only written to stdout. jwx does not create log files, cache files, or temporary files.
- **No third-party services.** jwx has no dependencies on external APIs or cloud services.

### What about the secret/key?

When using `jwx sign`, your secret (`--secret`) and private key (`--key`) are:

- Read from the command-line argument or file
- Used in-memory for the signing operation
- Never written to disk, logged, or transmitted

However, be aware of general security practices:

- Secrets passed as command-line arguments may be visible in your shell history. Consider using environment variables:
  ```bash
  jwx sign --alg HS256 --secret "$JWT_SECRET" '{"sub":"1234"}'
  ```
- Private key files should have restricted permissions:
  ```bash
  chmod 600 private.pem
  ```

## Web decoder: client-side only

The [web decoder](https://manimovassagh.github.io/jwx/) runs entirely in your browser:

### What happens when you paste a token

1. JavaScript in the page base64-decodes the header and payload
2. The decoded JSON is rendered in the right panel
3. The token is stored in the URL hash (fragment) for shareability

### What does NOT happen

- No data is sent to any server
- No JavaScript libraries are loaded from CDNs (everything is inline)
- No analytics or tracking scripts are present
- No cookies are set (only `localStorage` for theme preference)

### Content Security Policy

The web decoder enforces a strict Content Security Policy that prevents:

- Loading scripts from external origins
- Making fetch/XHR requests to any server
- Loading fonts from untrusted sources
- Connecting to WebSocket endpoints

This is enforced at the browser level via a `<meta>` CSP tag, providing defense-in-depth even if the page were somehow modified.

### URL hash privacy

When you paste a token, the URL updates to include the token in the hash fragment (the part after `#`). Important to know:

- **The hash fragment is never sent to the server** in HTTP requests -- this is a fundamental property of URLs
- The token is visible in your browser's address bar and history
- If you share the URL, anyone with the link can decode the token

For this reason, avoid pasting production tokens with sensitive claims into the web decoder if your browser history is shared or monitored.

## Best practices

### For development

- Use jwx freely with development and test tokens
- Use `--json` to pipe decoded data into scripts without exposing it in terminal output

### For production tokens

- Prefer the CLI over the web decoder for production tokens
- Use environment variables for secrets instead of command-line arguments
- Clear your shell history if you accidentally paste a sensitive token:
  ```bash
  history -d $(history | tail -1 | awk '{print $1}')  # Bash
  ```

### For shared environments

- Be cautious with `--clipboard` in shared desktop environments
- Don't share web decoder URLs containing production tokens
- Consider `--no-color` when piping to shared logs, as the raw token appears in the output

## Security reporting

If you discover a security vulnerability in jwx, please report it responsibly:

1. **Do not** open a public GitHub issue
2. See the [SECURITY.md](https://github.com/manimovassagh/jwx/blob/main/SECURITY.md) file in the repository for reporting instructions
3. You'll receive acknowledgment within 48 hours

## Dependencies

jwx has a minimal dependency footprint. The CLI is built with:

- [Cobra](https://github.com/spf13/cobra) -- CLI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) -- terminal styling
- [go-humanize](https://github.com/dustin/go-humanize) -- human-readable formatting
- [go-isatty](https://github.com/mattn/go-isatty) -- terminal detection

All cryptographic operations use Go's standard library (`crypto/*` packages), which is well-audited and maintained by the Go team.

The web decoder has **zero external JavaScript dependencies** -- all code is inline in a single HTML file.
