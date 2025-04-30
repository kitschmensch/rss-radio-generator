package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	defaultFeedTitle       = "My Radio Stations"
	defaultFeedDescription = "A collection of internet radio stations."
	defaultLanguage        = "en-us"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Itunes  string   `xml:"xmlns:itunes,attr"`
	Content string   `xml:"xmlns:content,attr"`
	Media   string   `xml:"xmlns:media,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Language    string `xml:"language"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string     `xml:"title"`
	Description string     `xml:"description"`
	Enclosure   *Enclosure `xml:"enclosure"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the base64 encoded data
	base64Data := r.URL.Query().Get("data")
	if base64Data == "" {
		http.Error(w, "No data provided", http.StatusBadRequest)
		return
	}

	// Convert URL-safe base64 to standard base64
	base64Data = strings.ReplaceAll(base64Data, "-", "+")
	base64Data = strings.ReplaceAll(base64Data, "_", "/")

	// Add padding if necessary
	padding := len(base64Data) % 4
	if padding > 0 {
		base64Data += strings.Repeat("=", 4-padding)
	}

	// Decode base64 string
	jsonBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid data encoding: %v", err), http.StatusBadRequest)
		return
	}

	// Parse JSON data
	var formData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Language    string `json:"language"`
		Stations    []struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			Description string `json:"description"`
		} `json:"stations"`
	}

	if err := json.Unmarshal(jsonBytes, &formData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Default values if not provided
	if formData.Title == "" {
		formData.Title = defaultFeedTitle
	}
	if formData.Description == "" {
		formData.Description = defaultFeedDescription
	}
	if formData.Language == "" {
		formData.Language = defaultLanguage
	}

	// Build RSS items
	var items []Item
	for _, station := range formData.Stations {
		items = append(items, Item{
			Title:       station.Title,
			Description: station.Description,
			Enclosure: &Enclosure{
				URL:    station.URL,
				Length: "0",
				Type:   "audio/mpeg",
			},
		})
	}

	// Generate RSS feed
	rss := RSS{
		Version: "2.0",
		Itunes:  "http://www.itunes.com/dtds/podcast-1.0.dtd",
		Content: "http://purl.org/rss/1.0/modules/content/",
		Media:   "http://search.yahoo.com/mrss/",
		Channel: Channel{
			Title:       formData.Title,
			Description: formData.Description,
			Language:    formData.Language,
			Items:       items,
		},
	}

	// Output RSS
	w.Header().Set("Content-Type", "application/rss+xml")
	xmlOutput, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		http.Error(w, "Failed to generate RSS", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(xml.Header))
	w.Write(xmlOutput)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	http.Handle("/rss", http.HandlerFunc(rssHandler))       // Serve RSS feed at /rss
	http.Handle("/", http.FileServer(http.Dir("./static"))) // Serve static files
	fmt.Println("Serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
