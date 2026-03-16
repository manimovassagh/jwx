package display

import (
	"encoding/json"

	"github.com/mani-sh-reddy/jwx/internal/jwt"
)

// JSONOutput represents the structured JSON output of a decoded token.
type JSONOutput struct {
	Header    map[string]interface{} `json:"header"`
	Payload   map[string]interface{} `json:"payload"`
	Signature string                 `json:"signature"`
	IsExpired bool                   `json:"is_expired"`
}

// RenderJSON returns the token as formatted JSON for piping.
func RenderJSON(token *jwt.DecodedToken) (string, error) {
	out := JSONOutput{
		Header:    token.Header,
		Payload:   token.Payload,
		Signature: token.Signature,
		IsExpired: token.IsExpired,
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
