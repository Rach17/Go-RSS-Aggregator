package repository

import (
	"context"
	"github.com/google/uuid"
	"database/sql"
	"github.com/Rach17/Go-RSS-Aggregator/db"
)

type FeedRepository interface {
	CreateFeed(ctx context.Context, title, url, description, language string) (db.Feed, error)
	GetFeedByID(ctx context.Context, id uuid.UUID) (db.Feed, error)
	GetFeedByURL(ctx context.Context, url string) (db.Feed, error)
	UpdateFeedLastFetchedAt(ctx context.Context, id uuid.UUID) error
}

type DBFeedRepository struct {
	queries *db.Queries
	db      *sql.DB
}

func NewDBFeedRepository(database *sql.DB) *DBFeedRepository {
	return &DBFeedRepository{
		queries: db.New(database),
		db:      database,
	}
}

func (r *DBFeedRepository) CreateFeed(ctx context.Context, title, url, description, language string) (db.Feed, error) {
	return r.queries.CreateFeed(ctx, db.CreateFeedParams{
		Title:       title,
		Url:         url,
		Description: sql.NullString{String: description, Valid: description != ""},
		Language:    language,
	})
}

func (r *DBFeedRepository) GetFeedByID(ctx context.Context, id uuid.UUID) (db.Feed, error) {
	feed, err := r.queries.GetFeedByID(ctx, id)
	if err != nil {
		return db.Feed{}, err
	}
	return feed, nil
}

func (r *DBFeedRepository) GetFeedByURL(ctx context.Context, url string) (db.Feed, error) {
	feed, err := r.queries.GetFeedByURL(ctx, url)
	if err != nil {
		return db.Feed{}, err
	}
	return feed, nil
}


func (r *DBFeedRepository) UpdateFeedLastFetchedAt(ctx context.Context, id uuid.UUID) error {
	return r.queries.UpdateFeedLastFetchedAt(ctx, id)
}

