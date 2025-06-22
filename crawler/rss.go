package crawler

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dkuaegis/aegis-geeknews-daily-bot/models"
)

// CrawlRSSFeed fetches and parses the RSS feed from the given URL
func CrawlRSSFeed(url string) (*models.Feed, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent header
	req.Header.Set("User-Agent", "GeekNews-Daily-Bot/1.0")

	// Make HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse XML
	var feed models.Feed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &feed, nil
}
