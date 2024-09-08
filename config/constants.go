package config

import (
	"encoding/base64"
	"fmt"
	"os"
)

var (
	AppName         = "yato"
	Version         = "0.1.0"
	MALOAuthBaseURL = "https://myanimelist.net/v1/oauth2/authorize"
	MALClientID     string
	MALClientSecret string
	MALRedirectURI  = "http://localhost:42069/authenticate"
)

// These variables will be set by the linker during build
var (
	encodedClientID     string
	encodedClientSecret string
)

func init() {
	var err error

	// Try to decode from build flags first
	MALClientID, err = decodeSecret(encodedClientID)
	if err != nil || MALClientID == "" {
		// Fallback to environment variable
		MALClientID = os.Getenv("MAL_CLIENT_ID")
	}

	MALClientSecret, err = decodeSecret(encodedClientSecret)
	if err != nil || MALClientSecret == "" {
		// Fallback to environment variable
		MALClientSecret = os.Getenv("MAL_CLIENT_SECRET")
	}

	if MALClientID == "" || MALClientSecret == "" {
		fmt.Println("Warning: Client ID or Secret not set. Please set MAL_CLIENT_ID and MAL_CLIENT_SECRET environment variables or use the build script.")
	}
}

func decodeSecret(encoded string) (string, error) {
	if encoded == "" {
		return "", fmt.Errorf("encoded string is empty")
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// GetMALClientID returns the MAL Client ID
func GetMALClientID() string {
	return MALClientID
}

// GetMALClientSecret returns the MAL Client Secret
func GetMALClientSecret() string {
	return MALClientSecret
}
