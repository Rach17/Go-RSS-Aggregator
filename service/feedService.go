package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Rach17/Go-RSS-Aggregator/data"
	"github.com/Rach17/Go-RSS-Aggregator/db"
	"github.com/Rach17/Go-RSS-Aggregator/repository"
	"github.com/google/uuid"
)

type FeedService struct {
	Repo       repository.FeedRepository
	HTTPClient *http.Client
}

func NewFeedService(repo repository.FeedRepository) *FeedService {
	return &FeedService{
		Repo: repo,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (fs *FeedService) CreateFeed(ctx context.Context, feedURL string) (db.Feed, error) {
	feed, err := fs.ValidateAndFetchNewFeed(ctx, feedURL)
	if err != nil {
		log.Printf("Error validating and fetching feed: %v", err)
		return db.Feed{}, fmt.Errorf("failed to validate and fetch feed: %w", err)
	}
	// Check if feed already exists
	exists, err := fs.FeedExists(ctx, feedURL)
	if err != nil {
		log.Printf("Error checking if feed exists: %v", err)
	}
	if exists {
		log.Printf("Feed already exists: %s", feedURL)
		return db.Feed{}, fmt.Errorf("feed already exists")
	}
	savedFeed, err := fs.Repo.CreateFeed(ctx, feed.Channel.Title, feedURL, feed.Channel.Description, feed.Channel.Language)
	if err != nil {
		log.Printf("Error creating feed: %v", err)
		return db.Feed{}, fmt.Errorf("failed to create feed: %w", err)
	}
	return savedFeed, nil
}

func (fs *FeedService) ValidateAndFetchNewFeed(ctx context.Context, feedURL string) (data.RSSFeed, error) {
	// Validate URL format
	if err := fs.validateURL(feedURL); err != nil {
		log.Printf("Invalid URL: %v", err)
		return data.RSSFeed{}, fmt.Errorf("invalid URL: %w", err)
	}

	// Create request with context
	respBody, err := fs.sendRequest(ctx, feedURL)
	if err != nil {
		log.Printf("Failed to fetch feed: %v", err)
		return data.RSSFeed{}, fmt.Errorf("failed to fetch feed: %w", err)
	}
	defer respBody.Close()

	// Parse RSS feed
	feed, err := fs.parseResponse(respBody)
	if err != nil {
		log.Printf("Failed to parse RSS feed: %v", err)
		return data.RSSFeed{}, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	return feed, nil
}

func (fs *FeedService) validateURL(feedURL string) error {
	parsedURL, err := url.Parse(feedURL)
	if err != nil {
		return err
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a valid host")
	}

	return nil
}

func (fs *FeedService) FeedExists(ctx context.Context, feedURL string) (bool, error) {
	feed, err := fs.Repo.GetFeedByURL(ctx, feedURL)
	if err != nil {
		// Handle "no rows in result set" as normal case - feed doesn't exist
		if strings.Contains(err.Error(), "no rows in result set") ||
			strings.Contains(err.Error(), "sql: no rows") {
			return false, nil
		}
		// Return actual database errors
		return false, err
	}

	// If we successfully fetched the feed, it exists
	if feed.ID != (uuid.UUID{}) {
		return true, nil
	}

	return false, nil
}

func (fs *FeedService) sendRequest(ctx context.Context, feedUrl string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent header
	req.Header.Set("User-Agent", "RSS-Aggregator/1.0")
	req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml")

	// Make HTTP request
	resp, err := fs.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("feed returned status %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (fs *FeedService) parseResponse(body io.Reader) (data.RSSFeed, error) {
	var feed data.RSSFeed
	if err := xml.NewDecoder(body).Decode(&feed); err != nil {
		return data.RSSFeed{}, fmt.Errorf("failed to parse response: %w", err)
	}
	return feed, nil
}

func (fs *FeedService) GetFeedByURL(ctx context.Context, feedURL string) (db.Feed, error) {
	feed, err := fs.Repo.GetFeedByURL(ctx, feedURL)
	if err != nil {
		return db.Feed{}, fmt.Errorf("failed to get feed by URL: %w", err)
	}
	return feed, nil
}


func (fs *FeedService) GetAllFeeds(ctx context.Context) ([]db.Feed, error) {
	feeds, err := fs.Repo.GetAllFeeds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all feeds: %w", err)
	}
	return feeds, nil
}


func (fs *FeedService) UpdateFeed(ctx context.Context, url string) error {
	feed, err := fs.Repo.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to get feed by URL: %w", err)
	}

	respBody, err := fs.sendRequest(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}
	defer respBody.Close()

	// Parse RSS feed
	fetchedFeed, err := fs.parseResponse(respBody)
	if err != nil {
		return fmt.Errorf("failed to parse RSS feed: %w", err)
	}
	
	f := data.RSSFeed{}
	f.DbFeedToRSSFeed(feed)

	if (f.Channel.LastBuildDate >= fetchedFeed.Channel.LastBuildDate) {
		log.Printf("No new updates for feed: %s", fetchedFeed.Channel.Title)
		return nil
	}

	// Update feed last fetched time
	if err := fs.Repo.UpdateFeedLastFetchedAt(ctx, url); err != nil {
		return fmt.Errorf("failed to update feed last fetched time: %w", err)
	}

	log.Printf("Successfully fetched and updated feed: %s", fetchedFeed.Channel.Title)
	return nil
}

func (fs *FeedService) FollowFeed(ctx context.Context, feedURL string, userID uuid.UUID)  error {
	feed, err := fs.Repo.GetFeedByURL(ctx, feedURL)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") || strings.Contains(err.Error(), "sql: no rows") {
			return fmt.Errorf("feed does not exist: %w", err)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to follow feed: %w", err)
	}

	// Associate the feed with the user
	if err := fs.Repo.FollowFeed(ctx, userID, feed.ID); err != nil {
		return fmt.Errorf("failed to follow feed: %w", err)
	}

	return nil
}

func (fs *FeedService) CreateAndFollowFeed(ctx context.Context, feedURL string, userID uuid.UUID) (db.Feed, error) {
	feed, err := fs.CreateFeed(ctx, feedURL)
	if err != nil {
		return db.Feed{}, fmt.Errorf("failed to create feed: %w", err)
	}

	// Follow the newly created feed
	if err := fs.Repo.FollowFeed(ctx, userID, feed.ID); err != nil {
		return db.Feed{}, fmt.Errorf("failed to follow feed: %w", err)
	}

	return feed, nil
}