package service

import (
    "context"
    "log"
    "sync"
    "time"

    "github.com/Rach17/Go-RSS-Aggregator/db"
    "github.com/Rach17/Go-RSS-Aggregator/repository"
)

type ScraperService struct {
    FeedService      *FeedService
    FeedRepo         repository.FeedRepository
    ticker           *time.Ticker
    stopChan         chan bool
    wg               sync.WaitGroup
	feedsToFetch     int
}

func NewScraperService(feedService *FeedService, feedRepo repository.FeedRepository, feedsToFetch int) *ScraperService {
    return &ScraperService{
        FeedService:   feedService,
        FeedRepo:      feedRepo,
        stopChan:      make(chan bool),
        feedsToFetch:  feedsToFetch,
    }
}

func (s *ScraperService) Start(interval time.Duration) {
    s.ticker = time.NewTicker(interval)
    s.wg.Add(1)
    
    log.Printf("Scraper started , interval: %v, feeds per cycle: %d", interval, s.feedsToFetch)

    go func() {
        defer s.wg.Done()
        log.Printf("Scraper scheduler started")
        
        for {
            select {
            case <-s.ticker.C:
                s.scrapeFeeds()
            case <-s.stopChan:
                log.Println("Scraper stopped")
                return
            }
        }
    }()
}

func (s *ScraperService) scrapeFeeds() {
    ctx := context.Background()
    
    // Get the least recently fetched feeds
    feeds, err := s.FeedRepo.GetLastFetchedFeeds(ctx, s.feedsToFetch)
    if err != nil {
        log.Printf("Error fetching feeds: %v", err)
        return
    }

    if len(feeds) == 0 {
        log.Println("No feeds found to scrape")
        return
    }

    log.Printf("Starting scrape cycle for %d feeds", len(feeds))
    
    // Create a semaphore to limit concurrent goroutines
    semaphore := make(chan struct{}, s.feedsToFetch)
    var scrapeWg sync.WaitGroup

    // Launch goroutines for each feed
    for i, feed := range feeds {
        scrapeWg.Add(1)
        
        go func(goroutineID int, feedData db.Feed) {
            defer scrapeWg.Done()
            
            // Acquire semaphore
            semaphore <- struct{}{}
            defer func() { <-semaphore }() // Release semaphore
            
            s.scrapeFeed(goroutineID, feedData)
        }(i+1, feed)
    }
    
    // Wait for all goroutines to complete
    scrapeWg.Wait()
    log.Printf("Scrape cycle completed for %d feeds", len(feeds))
}

func (s *ScraperService) scrapeFeed(goroutineID int, feed db.Feed) {
    ctx := context.Background()
    startTime := time.Now()
    
    log.Printf("Goroutine %d: Starting scrape for feed: %s (URL: %s)", 
        goroutineID, feed.Title, feed.Url)
    
    if err := s.FeedService.UpdateFeed(ctx, feed.Url); err != nil {
        log.Printf("Goroutine %d: Error updating feed %s: %v", 
            goroutineID, feed.Title, err)
    } else {
        duration := time.Since(startTime)
        log.Printf("Goroutine %d: Successfully updated feed: %s (took %v)", 
            goroutineID, feed.Title, duration)
    }
}

func (s *ScraperService) Stop() {
    if s.ticker != nil {
        s.ticker.Stop()
    }
    close(s.stopChan)
    s.wg.Wait()
    log.Println("Scraper service stopped gracefully")
}

func (s *ScraperService) ScrapeOnce() error {
    log.Println("Running one-time scrape...")
    s.scrapeFeeds()
    return nil
}
