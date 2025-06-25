package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	RSSFeedURL       string
	DiscordWebhookURL string
	CrawlCron        string
	NotificationCron string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if not in production
	env := os.Getenv("ENV")
	if env != "production" {
		if err := godotenv.Load(); err != nil {
			// Don't fail if .env file doesn't exist in non-production
			fmt.Printf("Warning: Could not load .env file: %v\n", err)
		}
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, fmt.Errorf("DB_HOST environment variable is required")
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return nil, fmt.Errorf("DB_PORT environment variable is required")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, fmt.Errorf("DB_USER environment variable is required")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("DB_NAME environment variable is required")
	}

	dbSSLMode := os.Getenv("DB_SSLMODE")
	if dbSSLMode == "" {
		dbSSLMode = "disable" // default value
	}

	rssFeedURL := os.Getenv("RSS_FEED_URL")
	if rssFeedURL == "" {
		rssFeedURL = "https://feeds.feedburner.com/geeknews-feed" // default value
	}

	discordWebhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if discordWebhookURL == "" {
		return nil, fmt.Errorf("DISCORD_WEBHOOK_URL environment variable is required")
	}

	crawlCron := os.Getenv("CRAWL_CRON")
	if crawlCron == "" {
		crawlCron = "59 * * * *" // default: every hour at 59 minutes
	}

	notificationCron := os.Getenv("NOTIFICATION_CRON")
	if notificationCron == "" {
		notificationCron = "0 3 * * *" // default: 12:00 KST (03:00 UTC)
	}

	return &Config{
		DBHost:           dbHost,
		DBPort:           dbPort,
		DBUser:           dbUser,
		DBPassword:       dbPassword,
		DBName:           dbName,
		DBSSLMode:        dbSSLMode,
		RSSFeedURL:       rssFeedURL,
		DiscordWebhookURL: discordWebhookURL,
		CrawlCron:        crawlCron,
		NotificationCron: notificationCron,
	}, nil
}
