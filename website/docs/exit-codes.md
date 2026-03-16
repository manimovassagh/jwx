---
sidebar_position: 7
title: Exit Codes
---

# Exit Codes

jwx uses meaningful exit codes so you can integrate it into scripts, CI/CD pipelines, and monitoring tools.

## Exit code reference

| Code | Meaning | When it occurs |
|------|---------|----------------|
| `0` | **Success** | Token was decoded/signed successfully; token is not expired |
| `1` | **Error** | Invalid or malformed token; missing required flags; file not found; any general error |
| `2` | **Expired** | Token was decoded successfully but the `exp` claim is in the past |

## Key behavior: expired tokens

When jwx decodes an expired token, it still shows the full decoded output (header, payload, signature) but exits with code **2** instead of 0. This design means:

- You always see the token contents, even when expired
- Scripts can distinguish between "invalid" (code 1) and "expired" (code 2)
- The exit code alone tells you whether the token is still valid

## Using exit codes in scripts

### Bash

```bash
#!/bin/bash

TOKEN="eyJhbGciOiJIUzI1NiIs..."

jwx "$TOKEN" > /dev/null 2>&1
EXIT_CODE=$?

case $EXIT_CODE in
  0)
    echo "Token is valid"
    ;;
  1)
    echo "Token is invalid or malformed"
    ;;
  2)
    echo "Token is expired"
    ;;
  *)
    echo "Unexpected exit code: $EXIT_CODE"
    ;;
esac
```

### Token validation function

```bash
validate_token() {
  local token="$1"
  jwx "$token" > /dev/null 2>&1
  local code=$?

  if [ "$code" -eq 0 ]; then
    return 0
  elif [ "$code" -eq 2 ]; then
    echo "Warning: token is expired" >&2
    return 2
  else
    echo "Error: invalid token" >&2
    return 1
  fi
}

# Usage
if validate_token "$MY_TOKEN"; then
  echo "Proceeding with valid token"
fi
```

### CI/CD pipeline check

```bash
# In a CI script -- fail the build if the token is expired or invalid
jwx "$SERVICE_TOKEN" > /dev/null 2>&1 || {
  echo "Service token is invalid or expired. Exiting."
  exit 1
}
```

### Extracting data regardless of expiry

```bash
# Always extract data, but warn on expiry
PAYLOAD=$(jwx --json "$TOKEN" 2>/dev/null | jq -r .payload)
EXIT_CODE=${PIPESTATUS[0]}

echo "Payload: $PAYLOAD"

if [ "$EXIT_CODE" -eq 2 ]; then
  echo "Note: this token is expired"
fi
```

## Comparison with other tools

| Tool | Valid | Invalid | Expired |
|------|-------|---------|---------|
| **jwx** | 0 | 1 | 2 |
| jwt-cli | 0 | 1 | 0 (no distinction) |
| jq (manual) | 0 | non-zero | 0 (no distinction) |

jwx's distinct exit code for expired tokens is unique among JWT CLI tools, making it more useful for automation.
