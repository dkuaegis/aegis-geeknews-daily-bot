package scheduler

import (
	"log/slog"
	"time"

	"github.com/dkuaegis/aegis-geeknews-daily-bot/crawler"
	"github.com/dkuaegis/aegis-geeknews-daily-bot/database"
	"github.com/dkuaegis/aegis-geeknews-daily-bot/discord"
	"github.com/go-co-op/gocron/v2"
	"github.com/jmoiron/sqlx"
)

// CrawlAndSave performs the RSS crawling and saves to database
func CrawlAndSave(db *sqlx.DB, rssFeedURL string) error {
	slog.Info("Starting RSS crawling...")

	// Crawl RSS feed
	feed, err := crawler.CrawlRSSFeed(rssFeedURL)
	if err != nil {
		slog.Error("RSS crawling failed", "error", err)
		return err
	}

	slog.Info("Crawled entries from RSS feed", "count", len(feed.Entries))

	// Save to database
	err = database.SaveNewsEntries(db, feed.Entries)
	if err != nil {
		slog.Error("Failed to save entries to database", "error", err)
		return err
	}

	slog.Info("Successfully completed RSS crawling and saved entries to database", "count", len(feed.Entries))
	return nil
}

func SendDiscordNotification(db *sqlx.DB, webhookURL string) error {
	slog.Info("Starting Discord notification...")

	entries, err := database.GetUnsentNewsEntries(db)
	if err != nil {
		slog.Error("Failed to get unsent news entries", "error", err)
		return err
	}

	if len(entries) == 0 {
		slog.Info("No unsent news entries found")
		return nil
	}

	slog.Info("Found unsent news entries", "count", len(entries))

	err = discord.SendNewsToDiscord(webhookURL, entries)
	if err != nil {
		slog.Error("Failed to send Discord notification", "error", err)
		return err
	}

	var newsIDs []int
	for _, entry := range entries {
		newsIDs = append(newsIDs, entry.ID)
	}

	err = database.MarkNewsAsSent(db, newsIDs)
	if err != nil {
		slog.Error("Failed to mark news as sent", "error", err)
		return err
	}

	slog.Info("Successfully sent news entries to Discord", "count", len(entries))
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
				slog.Error("Scheduled crawling failed", "error", err)
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
				slog.Error("Discord notification failed", "error", err)
			}
		}),
	)
	if err != nil {
		return nil, err
	}

	slog.Info("RSS crawling job created", "job_id", crawlJob.ID().String(), "cron", crawlCron)
	slog.Info("Discord notification job created", "job_id", discordJob.ID().String(), "cron", notificationCron)
	slog.Info("Scheduler configured", "crawl_cron", crawlCron, "notification_cron", notificationCron)

	return scheduler, nil
}
