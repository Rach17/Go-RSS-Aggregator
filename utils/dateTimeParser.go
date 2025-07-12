package utils

import (
	"fmt"
	"strings"
	"time"
)

// ParseRSSDate parses common RSS date formats
func ParseRSSDate(dateStr string) (time.Time, error) {
    if dateStr == "" {
        return time.Now(), nil
    }

    // Common RSS date formats
    formats := []string{
        time.RFC1123,     // "Mon, 02 Jan 2006 15:04:05 MST"
        time.RFC1123Z,    // "Mon, 02 Jan 2006 15:04:05 -0700"
        time.RFC822,      // "02 Jan 06 15:04 MST"
        time.RFC822Z,     // "02 Jan 06 15:04 -0700"
        time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
        "2006-01-02 15:04:05",
        "2006-01-02T15:04:05",
        "Mon, 2 Jan 2006 15:04:05 MST",
        "Mon, 2 Jan 2006 15:04:05 -0700",
    }

    // Try each format
    for _, format := range formats {
        if t, err := time.Parse(format, strings.TrimSpace(dateStr)); err == nil {
            return t, nil
        }
    }

    return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}