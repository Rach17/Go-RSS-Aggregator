package repository

import (
	"context"
	"time"
	"database/sql"
	"github.com/Rach17/Go-RSS-Aggregator/db"
	"github.com/google/uuid"
)

type FeedPostRepository interface {
	Create(ctx context.Context, feedID uuid.UUID, title, description, url string, publishedAt time.Time) error
	GetFeedPosts(ctx context.Context, feedURL string) ([]db.FeedPost, error)

}

type DBFeedPostRepository struct {
	queries *db.Queries
	db      *sql.DB
}

func NewFeedPostRepository(database *sql.DB) *DBFeedPostRepository {
	return &DBFeedPostRepository{
		db:      database,
		queries: db.New(database),
	}
}

func (r *DBFeedPostRepository) Create(ctx context.Context, feedID uuid.UUID, title, description, url string, publishedAt time.Time) error {
	desc := sql.NullString{String: description, Valid: true}
	if description == "" {
		desc = sql.NullString{String: "", Valid: false}
	}
	return r.queries.CreateFeedPost(ctx, db.CreateFeedPostParams{
		FeedID:      feedID,
		Title:       title,
		Description: desc,
		Url:        url,
		PublishedAt: publishedAt,
	})
}

