package scheduler

import (
	"log"
	"time"

	"github.com/dkuaegis/aegis-geeknews-daily-bot/crawler"
	"github.com/dkuaegis/aegis-geeknews-daily-bot/database"
	"github.com/dkuaegis/aegis-geeknews-daily-bot/discord"
	"github.com/go-co-op/gocron/v2"
	"github.com/jmoiron/sqlx"
)

// CrawlAndSave performs the RSS crawling and saves to database
func CrawlAndSave(db *sqlx.DB, rssFeedURL string) error {
	log.Println("Starting RSS crawling...")

	// Crawl RSS feed
	feed, err := crawler.CrawlRSSFeed(rssFeedURL)
	if err != nil {
		log.Printf("RSS crawling failed: %v", err)
		return err
	}

	log.Printf("Crawled %d entries from RSS feed", len(feed.Entries))

	// Save to database
	err = database.SaveNewsEntries(db, feed.Entries)
	if err != nil {
		log.Printf("Failed to save entries to database: %v", err)
		return err
	}

	log.Printf("Successfully completed RSS crawling and saved %d entries to database", len(feed.Entries))
	return nil
}

func SendDiscordNotification(db *sqlx.DB, webhookURL string) error {
	log.Println("Starting Discord notification...")

	entries, err := database.GetUnsentNewsEntries(db)
	if err != nil {
		log.Printf("Failed to get unsent news entries: %v", err)
		return err
	}

	if len(entries) == 0 {
		log.Println("No unsent news entries found")
		return nil
	}

	log.Printf("Found %d unsent news entries", len(entries))

	err = discord.SendNewsToDiscord(webhookURL, entries)
	if err != nil {
		log.Printf("Failed to send Discord notification: %v", err)
		return err
	}

	var newsIDs []int
	for _, entry := range entries {
		newsIDs = append(newsIDs, entry.ID)
	}

	err = database.MarkNewsAsSent(db, newsIDs)
	if err != nil {
		log.Printf("Failed to mark news as sent: %v", err)
		return err
	}

	log.Printf("Successfully sent %d news entries to Discord", len(entries))
	return nil
}

// StartScheduler starts the scheduler with configurable RSS crawling and Discord notifications
func StartScheduler(db *sqlx.DB, rssFeedURL, webhookURL, crawlCron, notificationCron string) (gocron.Scheduler, error) {
	// Set timezone to UTC
	loc, _ := time.LoadLocation("UTC")

	// Create scheduler
	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(loc),
	)
	if err != nil {
		return nil, err
	}

	// Add job to run crawling based on configured cron expression
	crawlJob, err := scheduler.NewJob(
		gocron.CronJob(crawlCron, false), // Configurable cron expression
		gocron.NewTask(func() {
			if err := CrawlAndSave(db, rssFeedURL); err != nil {
				log.Printf("Scheduled crawling failed: %v", err)
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	// Add job to run Discord notifications based on configured cron expression
	discordJob, err := scheduler.NewJob(
		gocron.CronJob(notificationCron, false), // Configurable cron expression
		gocron.NewTask(func() {
			if err := SendDiscordNotification(db, webhookURL); err != nil {
				log.Printf("Discord notification failed: %v", err)
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	log.Printf("RSS crawling job created with ID: %s (cron: %s)", crawlJob.ID().String(), crawlCron)
	log.Printf("Discord notification job created with ID: %s (cron: %s)", discordJob.ID().String(), notificationCron)
	log.Printf("Scheduler configured. RSS crawling: %s, Discord notifications: %s", crawlCron, notificationCron)

	return scheduler, nil
}
