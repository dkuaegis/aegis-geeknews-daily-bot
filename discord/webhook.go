package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dkuaegis/aegis-geeknews-daily-bot/models"
)

type WebhookMessage struct {
	Content string `json:"content"`
}

func SendNewsToDiscord(webhookURL string, newsEntries []models.News) error {
	if len(newsEntries) == 0 {
		return nil
	}

	var messageLines []string
	for _, entry := range newsEntries {
		line := fmt.Sprintf("- [%s](<%s>)", entry.Title, entry.URL)
		messageLines = append(messageLines, line)
	}

	message := WebhookMessage{
		Content: strings.Join(messageLines, "\n"),
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook message: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status code: %d", resp.StatusCode)
	}

	return nil
}