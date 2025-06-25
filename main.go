package main

import (
	"log/slog"
	"os"

	"github.com/dkuaegis/aegis-geeknews-daily-bot/config"
	"github.com/dkuaegis/aegis-geeknews-daily-bot/database"
	"github.com/dkuaegis/aegis-geeknews-daily-bot/scheduler"
)

func main() {
	// Setup structured logging with JSON format
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("GeekNews Daily Bot starting...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	slog.Info("Configuration loaded", "rss_url", cfg.RSSFeedURL)

	// Connect to database
	db, err := database.Connect(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Error closing database connection", "error", err)
		}
	}()

	// Create and start scheduler
	schedulerInstance, err := scheduler.StartScheduler(db, cfg.RSSFeedURL, cfg.DiscordWebhookURL, cfg.CrawlCron, cfg.NotificationCron)
	if err != nil {
		slog.Error("Failed to start scheduler", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := schedulerInstance.Shutdown(); err != nil {
			slog.Error("Error shutting down scheduler", "error", err)
		}
	}()

	// Start the scheduler
	schedulerInstance.Start()

	// Keep the program running
	slog.Info("GeekNews Daily Bot is running. Press Ctrl+C to stop.")
	select {} // Block forever
}
