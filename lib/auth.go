package lib

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"yato/config"
)

func GetNewCodeVerifier() (string, error) {
	bytes := make([]byte, 100)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate code verifier: %w", err)
	}

	encoded := base64.RawURLEncoding.EncodeToString(bytes)

	// Truncate to 128 characters
	if len(encoded) > 128 {
		encoded = encoded[:128]
	}

	return encoded, nil
}

func GetOAuthURL(codeVerifier string) string {
	return fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&code_challenge=%s&code_challenge_method=plain",
		config.MALOAuthBaseURL, config.MALClientID, config.MALRedirectURI, codeVerifier)
}

func OpenBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

func ExchangeToken(code, codeVerifier string) (*config.MyAnimeListConfig, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", config.MALClientID)
	data.Set("client_secret", config.MALClientSecret)
	data.Set("redirect_uri", config.MALRedirectURI)
	data.Set("code", code)

	if codeVerifier != "" {
		data.Set("code_verifier", codeVerifier)
	}

	// POST request to MAL API
	req, err := http.NewRequest("POST", "https://myanimelist.net/v1/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var malConfig config.MyAnimeListConfig
	var tokenResponse struct {
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	malConfig.TokenType = tokenResponse.TokenType
	malConfig.ExpiresIn = tokenResponse.ExpiresIn
	malConfig.AccessToken = tokenResponse.AccessToken
	malConfig.RefreshToken = tokenResponse.RefreshToken

	return &malConfig, nil
}
