# Test Coverage Analysis Report for jwx

**Analysis Date:** 2026-03-16
**Current Coverage:** 81.9%
**Target Coverage:** 90%
**Gap:** 8.1%

## Executive Summary

The jwx project has reasonably good test coverage at 81.9%, but falls short of the 90% target. Analysis of all `*_test.go` files and their corresponding source files identified **15 specific coverage gaps** across 11 source files. Implementing the recommended tests would add approximately **9.1% coverage gain**, reaching an estimated **90.8%+** coverage.

### Coverage by Package

| Package | Coverage | Status |
|---------|----------|--------|
| internal/display | 98.1% | Excellent |
| cmd/jwx/commands | 71.6% | Needs improvement |
| internal/jwt | 80.9% | Good but gaps exist |
| internal/clipboard | 0.0% | Untested |
| cmd/jwx | 0.0% | Main entry point untested |

## Critical Gaps (0% Coverage)

### 1. signEdDSA() - Ed25519 Signing (Severity: HIGH)
- **File:** `internal/jwt/sign.go:208`
- **Functions Untested:** signEdDSA()
- **Missing Test Cases:**
  - Signing token with Ed25519 private key
  - Missing key file error handling
  - Invalid key format error handling
  - Ed25519 algorithm roundtrip (sign + decode)
- **Impact:** EdDSA support completely untested
- **Estimated Gain:** +1.2%

### 2. clipboard.CommandFunc() - OS Platform Detection (Severity: MEDIUM)
- **File:** `internal/clipboard/clipboard.go:12`
- **Functions Untested:** CommandFunc()
- **Missing Test Cases:**
  - macOS (darwin) pbpaste detection
  - Linux xclip preference
  - Linux xsel fallback
  - Unsupported OS error handling
- **Impact:** Clipboard functionality untested across platforms
- **Estimated Gain:** +0.5%

### 3. main() - Entry Point (Severity: LOW)
- **File:** `cmd/jwx/main.go:17`
- **Functions Untested:** main()
- **Missing Test Cases:**
  - Version/commit/date setting
  - Execute() error propagation
- **Impact:** Main entry point not covered
- **Estimated Gain:** +0.3%

## High Priority Gaps (26-90% Coverage)

### 4. Execute() - Error Paths (26.3% Coverage - Severity: HIGH)
- **File:** `cmd/jwx/commands/root.go:44`
- **Problem:** Cobra error interception logic untested
- **Missing Test Cases:**
  - Unknown command containing JWT pattern
  - Flag parsing in error context
  - Non-JWT unknown command handling
  - Error output formatting
- **Estimated Gain:** +1.0%

### 5. runDecode() - Stdin & Expiration (84.6% Coverage - Severity: MEDIUM)
- **File:** `cmd/jwx/commands/decode.go:39`
- **Missing Test Cases:**
  - Reading from stdin (non-TTY pipe)
  - Expired token exit code (os.Exit(2))
  - Empty stdin error handling
  - TTY detection edge cases
- **Estimated Gain:** +0.8%

### 6. signECDSA() - ES384/ES512 Algorithms (78.1% Coverage - Severity: MEDIUM)
- **File:** `internal/jwt/sign.go:154`
- **Missing Test Cases:**
  - ES384 (P-384 curve) signing
  - ES512 (P-521 curve) signing
  - Signature size validation (96 bytes for ES384, 132 for ES512)
  - Hash algorithm selection verification
- **Estimated Gain:** +0.8%

### 7. signRSA() - RS384/RS512 Algorithms (76.2% Coverage - Severity: MEDIUM)
- **File:** `internal/jwt/sign.go:118`
- **Missing Test Cases:**
  - RS384 algorithm path
  - RS512 algorithm path
  - Hash function selection for each
- **Estimated Gain:** +0.6%

### 8. extractTime() - json.Number Type (50% Coverage - Severity: MEDIUM)
- **File:** `internal/jwt/decode.go:77`
- **Missing Test Cases:**
  - json.Number type handling (line 86-91)
  - Invalid json.Number error handling
- **Note:** Only reachable with custom JSON decoder
- **Estimated Gain:** +0.5%

### 9. Sign() - Edge Cases (86.2% Coverage - Severity: LOW)
- **File:** `internal/jwt/sign.go:39`
- **Missing Test Cases:**
  - Header/claims JSON encoding edge cases
  - Invalid algorithm type handling
  - Error wrapping verification
- **Estimated Gain:** +0.5%

### 10. decodeSegment() - Error Paths (85.7% Coverage - Severity: LOW)
- **File:** `internal/jwt/decode.go:63`
- **Missing Test Cases:**
  - Invalid base64 string
  - Valid base64 with invalid JSON
- **Estimated Gain:** +0.4%

### 11. runSign() - Error Handling (92.3% Coverage - Severity: LOW)
- **File:** `cmd/jwx/commands/sign.go:44`
- **Missing Test Cases:**
  - signFrom file read errors
  - JWT signing errors
  - Empty args with empty signFrom
- **Estimated Gain:** +0.3%

### 12-15. Other Minor Gaps (Estimated: +1.3% combined)
- Decode() - Invalid part count validation (+0.3%)
- formatValue() - Unmapped types (+0.3%)
- RenderJSON() - JSON marshal errors (+0.3%)
- signHMAC() - HS384/HS512 algorithms (+0.2%)
- Decode() - Part count edge cases (+0.2%)

## Test Quality Issues Identified

### 1. Inadequate Error Path Coverage
- Many functions have error paths that are never exercised
- Tests primarily focus on happy paths
- Example: `runDecode()` doesn't test stdin or expiration exit

### 2. Algorithm Variant Testing
- HMAC: HS256 tested, HS384/HS512 missing
- RSA: RS256 tested, RS384/RS512 missing
- ECDSA: ES256 tested, ES384/ES512 missing
- EdDSA: Not tested at all

### 3. Type-Specific Testing Gaps
- extractTime() doesn't test json.Number type
- formatValue() doesn't test array/object types
- No tests for OS-specific code paths

### 4. Integration Point Testing
- stdin/TTY detection untested
- Main entry point untested
- Clipboard platform detection untested

## Recommendations

### Priority 1 (Implement Immediately)
1. Add EdDSA tests - core functionality missing
2. Fix Execute() error handling - 73.7% untested
3. Complete ES384/ES512 tests - major algorithms missing
4. Add stdin testing to runDecode()

### Priority 2 (Implement Next Sprint)
5. Add RS384/RS512 tests
6. Add extractTime() json.Number tests
7. Complete runSign() error paths
8. Add clipboard.CommandFunc() platform tests

### Priority 3 (Polish)
9. Test remaining edge cases
10. Add main() integration test
11. Test formatValue() unmapped types

## Files to Modify

### Tests to Create/Expand
- `/Users/mani/Documents/Projects/jwx/internal/jwt/sign_test.go` - Add EdDSA, ES384, ES512, RS384, RS512 tests
- `/Users/mani/Documents/Projects/jwx/internal/jwt/decode_test.go` - Add json.Number, error path tests
- `/Users/mani/Documents/Projects/jwx/cmd/jwx/commands/commands_test.go` - Add stdin, clipboard.CommandFunc tests
- `/Users/mani/Documents/Projects/jwx/internal/display/render_test.go` - Add formatValue() type tests
- `/Users/mani/Documents/Projects/jwx/internal/display/json_test.go` - Add error handling tests

### New Test File (Optional)
- `/Users/mani/Documents/Projects/jwx/cmd/jwx/main_test.go` - Integration tests for main()

## Estimated Timeline

- Priority 1 tests: ~4-6 hours (gains +3.8%)
- Priority 2 tests: ~3-4 hours (gains +3.0%)
- Priority 3 tests: ~1-2 hours (gains +2.3%)

**Total effort:** ~8-12 hours to reach 90%+ coverage

## Conclusion

The project has good foundational test coverage but needs focused effort on algorithm variants and error paths. Implementing the 15 identified tests would systematically increase coverage from 81.9% to approximately 90.8%, addressing all major gaps while improving code reliability and maintainability.
