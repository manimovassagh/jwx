# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 1.x     | Yes       |
| < 1.0   | No        |

## Reporting a Vulnerability

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, please use [GitHub Security Advisories](https://github.com/manimovassagh/jwx/security/advisories/new) to report vulnerabilities privately.

### What to include

- Description of the vulnerability
- Steps to reproduce
- Impact assessment
- Affected versions
- Suggested fix (if any)

### Response timeline

- **Acknowledge**: Within 48 hours
- **Initial assessment**: Within 5 days
- **Fix for critical issues**: Within 7 days
- **Fix for non-critical**: Within 30 days

### Scope

Security issues for jwx include:

- Token parsing vulnerabilities that could lead to incorrect validation
- Key handling issues (accidental exposure, insecure defaults)
- Command injection via CLI arguments
- Dependencies with known CVEs

Out of scope:

- JWTs are not encrypted — payload visibility is by design
- Denial of service via large tokens (known limitation)
