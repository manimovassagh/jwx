---
name: jwt-decode
description: Decode and inspect JWT tokens using the jwx CLI
triggers:
  - "decode JWT"
  - "decode token"
  - "decode this token"
  - "what's in this token"
  - "what is in this token"
  - "inspect token"
  - "inspect JWT"
  - "parse JWT"
  - "parse token"
  - "eyJ"
---

# JWT Decode Skill

You are a JWT decoding assistant. When the user asks you to decode, inspect, or parse a JWT token, or pastes something that looks like a JWT (a string starting with `eyJ` containing two dots), use this skill.

## Steps

1. **Extract the token.** Look for a string matching the JWT format: three base64url-encoded segments separated by dots, typically starting with `eyJ`. The token may appear inline in the user's message or as a standalone string.

2. **Run the decode command.** Use the Bash tool to execute:

   ```
   jwx decode <token>
   ```

   If `jwx` is not in PATH, first build it:

   ```
   go build -o ./jwx ./cmd/jwx && ./jwx decode <token>
   ```

3. **Present the results.** Show the decoded output directly to the user. The `jwx` tool provides formatted output including:
   - Header (algorithm, token type, key ID)
   - Payload (all claims with human-readable timestamps)
   - Signature verification status (if applicable)

4. **Add context.** After showing the decoded output, briefly note anything interesting:
   - Whether the token is expired (compare `exp` claim to the current time)
   - The algorithm used
   - Any notable claims (roles, scopes, custom claims)

## Error Handling

- If the string does not look like a valid JWT, inform the user and explain the expected format (header.payload.signature, each base64url-encoded).
- If `jwx decode` returns an error, relay it to the user with an explanation.
