package main

import (
	"log"

	"github.com/dkuaegis/aegis-geeknews-daily-bot/config"
	"github.com/dkuaegis/aegis-geeknews-daily-bot/database"
	"github.com/dkuaegis/aegis-geeknews-daily-bot/scheduler"
)

func main() {
	log.Println("GeekNews Daily Bot starting...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded - RSS URL: %s", cfg.RSSFeedURL)

	// Connect to database
	db, err := database.Connect(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Create and start scheduler
	schedulerInstance, err := scheduler.StartScheduler(db, cfg.RSSFeedURL, cfg.DiscordWebhookURL)
	if err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer func() {
		if err := schedulerInstance.Shutdown(); err != nil {
			log.Printf("Error shutting down scheduler: %v", err)
		}
	}()

	// Start the scheduler
	schedulerInstance.Start()

	// Keep the program running
	log.Println("GeekNews Daily Bot is running. Press Ctrl+C to stop.")
	select {} // Block forever
}
