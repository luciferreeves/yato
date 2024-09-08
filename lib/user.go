package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yato/config"
)

type MALUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
	Location string `json:"location"`
	JoinedAt string `json:"joined_at"`
	Picture  string `json:"picture"`
}

func CurrentUser() (*MALUser, error) {
	var user MALUser

	req, err := http.NewRequest("GET", "https://api.myanimelist.net/v2/users/@me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.GetConfig().MyAnimeList.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}
