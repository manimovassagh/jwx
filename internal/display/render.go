package display

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/mani-sh-reddy/jwx/internal/jwt"
)

var (
	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")). // blue
			Padding(0, 1).
			MarginBottom(0)

	payloadStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("42")). // green
			Padding(0, 1).
			MarginBottom(0)

	signatureStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("214")). // yellow/orange
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Bold(true)

	headerTitle = titleStyle.Foreground(lipgloss.Color("63"))
	payloadTitle = titleStyle.Foreground(lipgloss.Color("42"))
	sigTitle     = titleStyle.Foreground(lipgloss.Color("214"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	warnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")). // red
			Bold(true)

	okStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42")). // green
		Bold(true)

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("251")).
			Bold(true)

	stringStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("114")) // light green

	numberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("178")) // gold
)

// Render produces the beautiful colorized output for a decoded JWT.
func Render(token *jwt.DecodedToken) string {
	var sections []string

	// Header
	headerContent := headerTitle.Render("Header") + "\n" + renderClaims(token.Header)
	sections = append(sections, headerStyle.Render(headerContent))

	// Payload
	payloadContent := payloadTitle.Render("Payload") + "\n" + renderPayload(token)
	sections = append(sections, payloadStyle.Render(payloadContent))

	// Signature
	sigContent := sigTitle.Render("Signature") + "\n" + renderSignature(token)
	sections = append(sections, signatureStyle.Render(sigContent))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func renderClaims(claims map[string]interface{}) string {
	keys := sortedKeys(claims)
	var lines []string
	for _, k := range keys {
		v := claims[k]
		lines = append(lines, fmt.Sprintf("  %s %s", keyStyle.Render(k+":"), formatValue(v)))
	}
	return strings.Join(lines, "\n")
}

func renderPayload(token *jwt.DecodedToken) string {
	// Sort keys but put standard claims first
	keys := sortedKeys(token.Payload)
	standardOrder := []string{"iss", "sub", "aud", "exp", "nbf", "iat", "jti"}
	ordered := orderKeys(keys, standardOrder)

	var lines []string
	for _, k := range ordered {
		v := token.Payload[k]
		line := fmt.Sprintf("  %s %s", keyStyle.Render(k+":"), formatValue(v))

		// Add timestamp annotations
		switch k {
		case "exp":
			if token.ExpiresAt != nil {
				ts := dimStyle.Render(fmt.Sprintf("(%s)", token.ExpiresAt.UTC().Format(time.RFC3339)))
				line = fmt.Sprintf("  %s %s  %s", keyStyle.Render(k+":"), formatValue(v), ts)
				if token.IsExpired {
					expired := humanize.Time(*token.ExpiresAt)
					line += "\n       " + warnStyle.Render("⚠ EXPIRED "+expired)
				} else {
					expiresIn := humanize.Time(*token.ExpiresAt)
					line += "\n       " + okStyle.Render("✓ Expires "+expiresIn)
				}
			}
		case "iat":
			if token.IssuedAt != nil {
				ts := dimStyle.Render(fmt.Sprintf("(%s)", token.IssuedAt.UTC().Format(time.RFC3339)))
				line = fmt.Sprintf("  %s %s  %s", keyStyle.Render(k+":"), formatValue(v), ts)
			}
		case "nbf":
			if token.NotBefore != nil {
				ts := dimStyle.Render(fmt.Sprintf("(%s)", token.NotBefore.UTC().Format(time.RFC3339)))
				line = fmt.Sprintf("  %s %s  %s", keyStyle.Render(k+":"), formatValue(v), ts)
			}
		}

		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func renderSignature(token *jwt.DecodedToken) string {
	alg := "unknown"
	if a, ok := token.Header["alg"].(string); ok {
		alg = a
	}

	sigPreview := token.Signature
	if len(sigPreview) > 32 {
		sigPreview = sigPreview[:32] + "..."
	}

	lines := []string{
		fmt.Sprintf("  %s %s", keyStyle.Render("Algorithm:"), formatValue(alg)),
		fmt.Sprintf("  %s %s", keyStyle.Render("Data:"), dimStyle.Render(sigPreview)),
		fmt.Sprintf("  %s %s", keyStyle.Render("Status:"), dimStyle.Render("Not verified (no key provided)")),
		"  " + dimStyle.Render("Use: jwx verify <token> --key <path>"),
	}
	return strings.Join(lines, "\n")
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return stringStyle.Render(fmt.Sprintf("%q", val))
	case float64:
		if val == float64(int64(val)) {
			return numberStyle.Render(fmt.Sprintf("%d", int64(val)))
		}
		return numberStyle.Render(fmt.Sprintf("%g", val))
	case bool:
		return numberStyle.Render(fmt.Sprintf("%t", val))
	case nil:
		return dimStyle.Render("null")
	default:
		return fmt.Sprintf("%v", val)
	}
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func orderKeys(keys []string, priority []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, p := range priority {
		for _, k := range keys {
			if k == p {
				result = append(result, k)
				seen[k] = true
				break
			}
		}
	}

	for _, k := range keys {
		if !seen[k] {
			result = append(result, k)
		}
	}

	return result
}
