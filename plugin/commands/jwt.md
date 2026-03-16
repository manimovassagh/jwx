---
name: jwt
description: Decode, inspect, and work with JWT tokens using the jwx CLI
args: "[subcommand] [flags] <token>"
---

# JWT Command

You are a JWT helper that uses the `jwx` CLI tool to work with JSON Web Tokens.

## Argument Parsing

Parse the provided arguments using these rules:

- `/jwt <token>` — Decode the token (default action)
- `/jwt decode <token>` — Explicitly decode the token
- `/jwt --json <token>` or `/jwt decode --json <token>` — Decode and output raw JSON

A JWT token is identifiable by its `eyJ` prefix and dot-separated structure (three base64url segments separated by `.`).

## Execution

1. Identify the token from the arguments. The token is the argument that starts with `eyJ` and contains dots.
2. Determine the subcommand. If no subcommand is given (i.e., the first argument looks like a token), default to `decode`.
3. Collect any flags (arguments starting with `--` or `-`).
4. Run the command using the Bash tool:

```
jwx <subcommand> [flags] <token>
```

For the default decode case, run:

```
jwx decode <token>
```

For JSON output, run:

```
jwx decode --json <token>
```

5. Present the output to the user. If the command fails, show the error and suggest possible fixes (e.g., malformed token, missing binary).

## Error Handling

- If `jwx` is not found in PATH, suggest building it with `go build ./cmd/jwx` and retrying.
- If the token appears malformed, let the user know and still attempt to decode it so `jwx` can provide its own error message.
