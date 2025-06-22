package models

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"time"
)

// Author represents the author information in RSS feed
type Author struct {
	Name  string `xml:"name" json:"name"`
	Email string `xml:"email" json:"email"`
	URI   string `xml:"uri" json:"uri"`
}

// Link represents a link element with attributes
type Link struct {
	Rel  string `xml:"rel,attr" json:"rel,omitempty"`
	Type string `xml:"type,attr" json:"type,omitempty"`
	Href string `xml:"href,attr" json:"href"`
}

// Entry represents each news item in the RSS feed
type Entry struct {
	Title     string    `xml:"title" json:"title"`
	Link      Link      `xml:"link" json:"link"`
	ID        string    `xml:"id" json:"id"`
	Updated   time.Time `xml:"updated" json:"updated"`
	Published time.Time `xml:"published" json:"published"`
	Author    Author    `xml:"author" json:"author"`
	Content   string    `xml:"content" json:"content"`
}

// Category represents the category element
type Category struct {
	Term string `xml:"term,attr" json:"term"`
}

// Feed represents the main RSS feed structure
type Feed struct {
	XMLName  xml.Name  `xml:"feed"`
	Title    string    `xml:"title" json:"title"`
	Subtitle string    `xml:"subtitle" json:"subtitle"`
	Link     []Link    `xml:"link" json:"links"`
	ID       string    `xml:"id" json:"id"`
	Author   Author    `xml:"author" json:"author"`
	Updated  time.Time `xml:"updated" json:"updated"`
	Icon     string    `xml:"icon" json:"icon"`
	Logo     string    `xml:"logo" json:"logo"`
	Category Category  `xml:"category" json:"category"`
	Entries  []Entry   `xml:"entry" json:"entries"`
}

// String implements the Stringer interface to output JSON format
func (f Feed) String() string {
	buffer := &strings.Builder{}
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false) // Prevent HTML escaping in JSON output

	err := encoder.Encode(f)
	if err != nil {
		return "Error marshaling to JSON: " + err.Error()
	}

	// Remove the trailing newline that encoder.Encode adds
	result := buffer.String()
	return strings.TrimSuffix(result, "\n")
}
