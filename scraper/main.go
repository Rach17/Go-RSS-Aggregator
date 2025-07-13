package main

import (
    "database/sql"
    "log"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"

    "github.com/Rach17/Go-RSS-Aggregator/repository"
    "github.com/Rach17/Go-RSS-Aggregator/service"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

func main() {
    // Load environment variables
    err := godotenv.Load(".env")
    if err != nil {
        log.Println("No .env file found, using default environment variables")
    }

    // Get database URL
    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatal("DB_URL environment variable is not set")
    }

    // Get configuration
    config := getScraperConfig()

    // Setup database connection
    connection, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }
    defer connection.Close()

    // Test database connection
    if err := connection.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    // Initialize repositories and services
    feedRepo := repository.NewDBFeedRepository(connection)
    feedPostRepo := repository.NewDBFeedPostRepository(connection)
    feedService := service.NewFeedService(feedRepo, feedPostRepo)
    scraperService := service.NewScraperService(feedService, feedRepo, config.FeedsToFetch)


    // Setup graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Start scraper
    log.Printf("Starting RSS scraper with config: %+v", config)
    scraperService.Start(config.Interval)

    // Run initial scrape if configured
    if config.InitialScrape {
        log.Println("Running initial scrape...")
        scraperService.ScrapeOnce()
    }

    // // Print scraper stats
    // stats := scraperService.GetScraperStats()
    // log.Printf("Scraper stats: %+v", stats)

    // Wait for shutdown signal
    <-sigChan
    log.Println("Shutdown signal received, stopping scraper...")
    
    scraperService.Stop()
    log.Println("Scraper shutdown complete")
}

type ScraperConfig struct {
    Interval       time.Duration
    FeedsToFetch   int
    InitialScrape  bool
}

func getScraperConfig() ScraperConfig {
    config := ScraperConfig{
        Interval:      60 * time.Minute, // Default: 1 hour
        FeedsToFetch:  10,               // Default: 10 feeds per cycle
        InitialScrape: true,             // Default: run initial scrape
    }

    // Get scraper interval (in minutes)
    if intervalStr := os.Getenv("SCRAPER_INTERVAL_MINUTES"); intervalStr != "" {
        if intervalMinutes, err := strconv.Atoi(intervalStr); err == nil {
            config.Interval = time.Duration(intervalMinutes) * time.Minute
        } else {
            log.Printf("Invalid SCRAPER_INTERVAL_MINUTES: %v, using default", err)
        }
    }

    // Get number of feeds to fetch per cycle
    if feedsStr := os.Getenv("SCRAPER_FEEDS_TO_FETCH"); feedsStr != "" {
        if feeds, err := strconv.Atoi(feedsStr); err == nil && feeds > 0 {
            config.FeedsToFetch = int(feeds)
        } else {
            log.Printf("Invalid SCRAPER_FEEDS_TO_FETCH: %v, using default", err)
        }
    }

    // Get initial scrape setting
    if initialStr := os.Getenv("SCRAPER_INITIAL_SCRAPE"); initialStr != "" {
        if initial, err := strconv.ParseBool(initialStr); err == nil {
            config.InitialScrape = initial
        } else {
            log.Printf("Invalid SCRAPER_INITIAL_SCRAPE: %v, using default", err)
        }
    }

    return config
}