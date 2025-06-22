package models

import "time"

// News represents a news entry in database
type News struct {
	ID          int       `db:"id"`
	URL         string    `db:"url"`
	Title       string    `db:"title"`
	Author      string    `db:"author"`
	Content     string    `db:"content"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	Sent        bool      `db:"sent"`
}

// ConvertEntryToNews converts RSS Entry to News for database storage
func ConvertEntryToNews(entry Entry) News {
	return News{
		URL:         entry.Link.Href,
		Title:       entry.Title,
		Author:      entry.Author.Name,
		Content:     entry.Content,
		PublishedAt: entry.Published,
	}
}
