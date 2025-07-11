package data

import (
	"encoding/xml"
	"github.com/Rach17/Go-RSS-Aggregator/db"
)

type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Language    string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
}

func (f *RSSFeed) DbFeedToRSSFeed(feed db.Feed) {
	f.XMLName = xml.Name{Local: "rss"}
	f.Channel = Channel{
		Title:       feed.Title,
		Description: feed.Description.String,
		Link:        feed.Url,
		Language:    feed.Language,
		LastBuildDate: feed.LastFetchedAt.Time.Format("Mon, 02 Jan 2006 15:04:05 MST"),
	}
}
