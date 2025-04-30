package main

import (
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
	// Extract query parameters for feed title and description
	query := r.URL.Query()
	feedTitle := query.Get("title")
	if feedTitle == "" {
		feedTitle = defaultFeedTitle
	}
	feedDescription := query.Get("description")
	if feedDescription == "" {
		feedDescription = defaultFeedDescription
	}
	language := query.Get("language")
	if language == "" {
		language = defaultLanguage
	}

	stationsParam := query.Get("stations")
	if stationsParam == "" {
		http.Error(w, "No stations provided", http.StatusBadRequest)
		return
	}

	// Parse the stations parameter (format: "Title1|URL1|Description1,Title2|URL2|Description2,...")
	stations := strings.Split(stationsParam, ",")
	var items []Item
	for _, station := range stations {
		parts := strings.SplitN(station, "|", 3)
		if len(parts) != 3 {
			http.Error(w, "Invalid station format. Use Title|URL|Description.", http.StatusBadRequest)
			return
		}
		title := parts[0]
		url := parts[1]
		description := parts[2]
		items = append(items, Item{
			Title:       title,
			Description: description,
			Enclosure: &Enclosure{
				URL:    url,
				Length: "0",
				Type:   "audio/mpeg",
			},
		})
	}

	rss := RSS{
		Version: "2.0",
		Itunes:  "http://www.itunes.com/dtds/podcast-1.0.dtd",
		Content: "http://purl.org/rss/1.0/modules/content/",
		Media:   "http://search.yahoo.com/mrss/",
		Channel: Channel{
			Title:       feedTitle,
			Description: feedDescription,
			Language:    language,
			Items:       items,
		},
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	xmlOutput, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		http.Error(w, "Failed to generate RSS", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(xml.Header))
	w.Write(xmlOutput)
}

func main() {
	http.HandleFunc("/", rssHandler)
	fmt.Println("Serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
