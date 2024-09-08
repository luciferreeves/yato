package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"yato/config"
)

type Recommendation struct {
	MALId string `json:"mal_id"`
	URL   string `json:"url"`
	Entry []struct {
		MALId  int    `json:"mal_id"`
		URL    string `json:"url"`
		Images struct {
			JPG struct {
				ImageURL      string `json:"image_url"`
				SmallImageURL string `json:"small_image_url"`
				LargeImageURL string `json:"large_image_url"`
			} `json:"jpg"`
			WebP struct {
				ImageURL      string `json:"image_url"`
				SmallImageURL string `json:"small_image_url"`
				LargeImageURL string `json:"large_image_url"`
			} `json:"webp"`
		} `json:"images"`
		Title string `json:"title"`
	} `json:"entry"`
	Content string `json:"content"`
	Date    string `json:"date"`
	User    struct {
		URL      string `json:"url"`
		Username string `json:"username"`
	} `json:"user"`
}

type recommendationsResponse struct {
	Pagination struct {
		LastVisiblePage int  `json:"last_visible_page"`
		HasNextPage     bool `json:"has_next_page"`
	} `json:"pagination"`
	Data []Recommendation `json:"data"`
}

func getRecentRecommendations(mediaType string) ([]Recommendation, error) {
	url := fmt.Sprintf("%s/recommendations/%s", config.JikanAPIBaseURL, mediaType)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response recommendationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

func GetRecentAnimeRecommendations() ([]Recommendation, error) {
	return getRecentRecommendations("anime")
}

func GetRecentMangaRecommendations() ([]Recommendation, error) {
	return getRecentRecommendations("manga")
}
