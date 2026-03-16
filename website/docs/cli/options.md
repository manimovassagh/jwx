---
sidebar_position: 3
title: "Global Options"
---

# Global Options

These options are available across all jwx commands.

## --no-color

Disable colored terminal output.

```bash
jwx --no-color eyJhbGciOiJIUzI1NiIs...
jwx decode --no-color eyJhbGciOiJIUzI1NiIs...
```

jwx also respects the `NO_COLOR` environment variable, following the [no-color.org](https://no-color.org/) convention:

```bash
export NO_COLOR=1
jwx eyJhbGciOiJIUzI1NiIs...  # output will be uncolored
```

This is useful when:
- Piping output to a file or another program
- Running in a CI/CD environment that doesn't support ANSI escapes
- Using a terminal emulator with limited color support
- Accessibility needs

## --json / -j

Output results as machine-readable JSON instead of the colorized display.

```bash
# Decode as JSON
jwx --json eyJhbGciOiJIUzI1NiIs...

# Sign and output as JSON
jwx sign --alg HS256 --secret key --json '{"sub":"1234"}'
```

For decode, the JSON output contains `header`, `payload`, and `signature` fields. For sign, it outputs `{"token":"..."}`.

This flag is ideal for scripting with `jq`:

```bash
# Extract the subject claim
jwx --json "$TOKEN" | jq -r .payload.sub

# Pretty-print the full payload
jwx --json "$TOKEN" | jq .payload

# Get the signing algorithm
jwx --json "$TOKEN" | jq -r .header.alg
```

## --clipboard / -c

Read the JWT token from the system clipboard instead of requiring it as an argument or via stdin.

```bash
# Copy a token to your clipboard, then:
jwx --clipboard
jwx -c
```

This works across platforms:
- **macOS**: Uses `pbpaste`
- **Linux**: Uses `xclip` or `xsel` (must be installed)
- **Windows**: Uses the Win32 clipboard API

## --help / -h

Display help for any command:

```bash
jwx --help
jwx decode --help
jwx sign --help
```

## --version / -v

Display the jwx version, commit hash, and build date:

```bash
jwx version
```

## Environment variables

| Variable | Description |
|----------|-------------|
| `NO_COLOR` | When set to any value, disables colored output (equivalent to `--no-color`) |

## Flag precedence

When multiple input sources conflict, jwx uses this precedence:

1. **Command-line argument** -- highest priority
2. **`--clipboard` flag** -- reads from clipboard
3. **Stdin** -- reads from piped input (only when stdin is not a terminal)

If no input is provided through any of these methods, jwx shows the help text.
