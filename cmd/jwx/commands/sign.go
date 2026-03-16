package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/manimovassagh/jwx/internal/jwt"
)

var (
	signAlg    string
	signSecret string
	signKey    string
	signFrom   string
	signJSON   bool
)

var signCmd = &cobra.Command{
	Use:   "sign [claims-json]",
	Short: "Create and sign a JWT token",
	Long: `Create a signed JWT token from JSON claims.

HMAC (symmetric):
  jwx sign --alg HS256 --secret mykey '{"sub":"1234","name":"John"}'

RSA (asymmetric):
  jwx sign --alg RS256 --key private.pem '{"sub":"1234"}'

From file:
  jwx sign --alg HS256 --secret mykey --from claims.json`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSign,
}

func init() {
	signCmd.Flags().StringVar(&signAlg, "alg", "", "Signing algorithm (HS256, RS256, ES256, EdDSA, etc.)")
	signCmd.Flags().StringVar(&signSecret, "secret", "", "Secret key for HMAC algorithms")
	signCmd.Flags().StringVar(&signKey, "key", "", "Path to private key file (PEM) for RSA/EC/EdDSA")
	signCmd.Flags().StringVar(&signFrom, "from", "", "Read claims from a JSON file instead of argument")
	signCmd.Flags().BoolVarP(&signJSON, "json", "j", false, "Output as JSON ({\"token\":\"...\"})")
	_ = signCmd.MarkFlagRequired("alg")
}

func runSign(cmd *cobra.Command, args []string) error {
	var claims string

	switch {
	case signFrom != "":
		data, err := os.ReadFile(signFrom)
		if err != nil {
			return fmt.Errorf("failed to read claims file: %w", err)
		}
		claims = string(data)
	case len(args) > 0:
		claims = args[0]
	default:
		return fmt.Errorf("no claims provided\n\nUsage: jwx sign --alg HS256 --secret mykey '{\"sub\":\"1234\"}'")
	}

	token, err := jwt.Sign(jwt.SignOptions{
		Algorithm: signAlg,
		Secret:    signSecret,
		KeyFile:   signKey,
		Claims:    claims,
	})
	if err != nil {
		return err
	}

	if signJSON {
		out, err := json.Marshal(map[string]string{"token": token})
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(out))
	} else {
		fmt.Println(token)
	}
	return nil
}
