package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// DecodedToken holds the parsed parts of a JWT.
type DecodedToken struct {
	Header    map[string]interface{}
	Payload   map[string]interface{}
	Signature string
	Raw       string
	IsExpired bool
	ExpiresAt *time.Time
	IssuedAt  *time.Time
	NotBefore *time.Time
}

// Decode parses a JWT token string without verifying the signature.
func Decode(tokenStr string) (*DecodedToken, error) {
	tokenStr = strings.TrimSpace(tokenStr)

	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token: expected 3 parts, got %d", len(parts))
	}

	header, err := decodeSegment(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid header: %w", err)
	}

	payload, err := decodeSegment(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	dt := &DecodedToken{
		Header:    header,
		Payload:   payload,
		Signature: parts[2],
		Raw:       tokenStr,
	}

	if exp, ok := extractTime(payload, "exp"); ok {
		dt.ExpiresAt = &exp
		dt.IsExpired = time.Now().After(exp)
	}
	if iat, ok := extractTime(payload, "iat"); ok {
		dt.IssuedAt = &iat
	}
	if nbf, ok := extractTime(payload, "nbf"); ok {
		dt.NotBefore = &nbf
	}

	return dt, nil
}

func decodeSegment(seg string) (map[string]interface{}, error) {
	data, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func extractTime(claims map[string]interface{}, key string) (time.Time, bool) {
	val, ok := claims[key]
	if !ok {
		return time.Time{}, false
	}

	switch v := val.(type) {
	case float64:
		return time.Unix(int64(v), 0), true
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return time.Time{}, false
		}
		return time.Unix(n, 0), true
	default:
		return time.Time{}, false
	}
}
