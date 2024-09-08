package lib

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"yato/config"
)

// ImageCache handles caching and retrieving images
type ImageCache struct {
	cacheDir string
}

// NewImageCache creates a new ImageCache
func NewImageCache() *ImageCache {
	cacheDir := filepath.Join(config.ConfigDir, config.AppName, "cache")
	return &ImageCache{cacheDir: cacheDir}
}

// GetImage retrieves an image, either from cache or by downloading it
func (c *ImageCache) GetImage(mediaType string, malID int, size string, url string) (image.Image, error) {
	cachePath := c.getCachePath(mediaType, malID, size)

	// Check if the image is already cached
	if img, err := c.loadFromCache(cachePath); err == nil {
		return img, nil
	}

	// If not cached, download and cache the image
	return c.downloadAndCache(url, cachePath)
}

func (c *ImageCache) getCachePath(mediaType string, malID int, size string) string {
	return filepath.Join(c.cacheDir, mediaType, fmt.Sprintf("%d", malID), size+".jpg")
}

func (c *ImageCache) loadFromCache(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (c *ImageCache) downloadAndCache(url, cachePath string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: %s", resp.Status)
	}

	// Ensure the cache directory exists
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return nil, err
	}

	// Create the cache file
	cacheFile, err := os.Create(cachePath)
	if err != nil {
		return nil, err
	}
	defer cacheFile.Close()

	// Download and write to cache file
	_, err = io.Copy(cacheFile, resp.Body)
	if err != nil {
		return nil, err
	}

	// Reset file pointer and decode the image
	_, err = cacheFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(cacheFile)
	if err != nil {
		return nil, err
	}

	return img, nil
}
