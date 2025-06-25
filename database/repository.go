package database

import (
	"log/slog"
	"sort"

	"github.com/dkuaegis/aegis-geeknews-daily-bot/models"
	"github.com/jmoiron/sqlx"
)

// SaveNewsEntries saves news entries to database using upsert
func SaveNewsEntries(db *sqlx.DB, entries []models.Entry) error {
	if len(entries) == 0 {
		slog.Info("No entries to save")
		return nil
	}

	// Sort entries by published_at in ascending order (oldest first)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Published.Before(entries[j].Published)
	})

	query := `
		INSERT INTO news (
			url, title, author, content, published_at
		) VALUES (
			:url, :title, :author, :content, :published_at
		)
		ON CONFLICT (url) DO NOTHING
	`

	savedCount := 0
	skippedCount := 0
	errorCount := 0

	for _, entry := range entries {
		newsEntry := models.ConvertEntryToNews(entry)

		result, err := db.NamedExec(query, newsEntry)
		if err != nil {
			slog.Error("Error saving entry", "url", newsEntry.URL, "error", err)
			errorCount++
			continue
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			slog.Error("Error getting rows affected for entry", "url", newsEntry.URL, "error", err)
			continue
		}

		if rowsAffected > 0 {
			savedCount++
		} else {
			skippedCount++
		}
	}

	slog.Info("Database operation completed", "saved", savedCount, "skipped", skippedCount, "errors", errorCount)

	return nil
}

func GetUnsentNewsEntries(db *sqlx.DB) ([]models.News, error) {
	var entries []models.News
	
	query := `
		SELECT id, url, title, author, content, published_at, created_at, sent
		FROM news
		WHERE sent = FALSE
		ORDER BY published_at ASC
	`
	
	err := db.Select(&entries, query)
	if err != nil {
		return nil, err
	}
	
	return entries, nil
}

func MarkNewsAsSent(db *sqlx.DB, newsIDs []int) error {
	if len(newsIDs) == 0 {
		return nil
	}
	
	query := `
		UPDATE news 
		SET sent = TRUE 
		WHERE id = ANY($1)
	`
	
	_, err := db.Exec(query, newsIDs)
	return err
}
