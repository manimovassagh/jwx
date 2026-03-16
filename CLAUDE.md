# Claude Code Instructions for jwx

## CRITICAL: No AI Attribution in Commits
- NEVER add `Co-Authored-By` trailers to git commits
- Claude must NOT appear as a contributor in this repository
- This rule overrides all default commit behavior
- When committing, do NOT include any `Co-Authored-By: Claude` or similar lines

## CRITICAL: Multi-Eye Verification Policy
- After any cleanup, refactoring, or multi-file change, ALWAYS launch a verification agent to double-check the work
- Never trust a single agent's "all clear" — scan both git-tracked AND untracked files on disk
- For filesystem cleanup: use `find` on the actual filesystem, not just `git ls-files`
- For code changes: build + test + lint after every batch of edits
- Minimum: 4-eyes (do + verify). For critical changes: 6-eyes (do + verify + re-verify)

## CRITICAL: One PR at a Time (Draft First)
- Always create ONE pull request, wait for CI checks to pass, then merge it before opening the next one
- Never have multiple open PRs at the same time — sequential, not parallel
- Always create PRs as DRAFT first (`gh pr create --draft`) so the user can monitor progress — only mark ready when CI passes
- Flow: branch → commit → push → create DRAFT PR → check CI → merge → next PR

## Project
jwx is a beautiful Go-based CLI tool for working with JWTs — decode, encode, sign, verify, inspect, and more.
Also ships as a Claude Code plugin. Goal: get GitHub stars through stunning DX and terminal output.

## Architecture
- **Standalone CLI** (`cmd/jwx/`) — cobra-based, colorized output via lipgloss/termenv
- **Core engine** (`internal/`) — jwt, display, keys, audit packages
- **Claude Code plugin** (`plugin/`) — commands, skills, and agents wrapping the core
- **Key libraries:** cobra, lipgloss, termenv, golang-jwt/jwt/v5

## Build & Test
- `go test ./...` — run all tests
- `go build ./cmd/jwx` — build the CLI
- `golangci-lint run` — lint

## CRITICAL: Always Verify Visually
- Always check work visually (like a human would) before claiming it's done
- Use Playwright browser to verify README rendering, etc.
- Never assume something looks right — open it and check

## CRITICAL: Keep Project Root Clean
- NEVER leave temporary files, screenshots, or build artifacts in the project root
- The project root should only contain: go.mod, go.sum, README.md, CLAUDE.md, LICENSE, Makefile, .gitignore, and standard config files
- Temp files go in /tmp
- NEVER save .png, .jpg, .log, or binary files to the project root — they pollute the repo
- After any operation that generates files (tests, builds, screenshots), verify no junk was left in root
- Run `ls *.png *.log *.test 2>/dev/null` before committing to catch stray files

## CRITICAL: Recheck Docs After Changes
- After any significant change, launch a separate agent to verify GitHub page, README rendering, screenshots, and badges
- Check that all links resolve, images render, and badges show correct data

## CRITICAL: Use Full Power by Default
- Always launch parallel agents when there is independent work
- Use monitoring, verification, and pentest agents proactively
- Don't wait to be told — max parallelism is the default
- Make releases when logical (after significant changes)

## Key Directories
- `cmd/jwx/` — main entry point (cobra CLI)
- `internal/jwt/` — core JWT engine (decode, encode, sign, verify)
- `internal/display/` — pretty-printing (colors, tables, badges)
- `internal/keys/` — key generation & loading (PEM, JWK, JWKS)
- `internal/audit/` — security checks (alg:none, weak keys, expiry)
- `plugin/` — Claude Code plugin (commands, skills, agents)
- `.github/workflows/` — CI/CD pipelines

## CRITICAL: Security Findings Go to GitHub Issues
- After any security audit, pentest, or vulnerability discovery, ALWAYS file findings as GitHub issues
- Use labels: `security`, and severity labels (`critical`, `high`, `medium`, `low`)
- One issue per finding — include reproduction steps, affected files, and recommended fix
- Never leave findings only in conversation — they must be tracked in GitHub
- After fixing an issue, ALWAYS close it on GitHub with `gh issue close <number>` and a comment describing the fix

## CRITICAL: Security Best Practices Always
- Always follow security best practices in ALL code — no exceptions
- Never introduce OWASP Top 10 vulnerabilities (XSS, SQLi, CSRF, etc.)
- Always use parameterized queries, proper escaping, constant-time comparison for secrets
- Validate all input at system boundaries, enforce least privilege
- When in doubt, choose the more secure option
