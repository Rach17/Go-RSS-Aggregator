package repository

import (
	"context"
	"time"
	"database/sql"
	"github.com/Rach17/Go-RSS-Aggregator/db"
	"github.com/google/uuid"
)

type FeedPostRepository interface {
	Create(ctx context.Context, feedID uuid.UUID, title, description, url, author string, publishedAt time.Time) error
	GetFeedPosts(ctx context.Context, feedURL string) ([]db.FeedPost, error)

}

type DBFeedPostRepository struct {
	queries *db.Queries
	db      *sql.DB
}

func NewDBFeedPostRepository(database *sql.DB) *DBFeedPostRepository {
	return &DBFeedPostRepository{
		db:      database,
		queries: db.New(database),
	}
}

func (r *DBFeedPostRepository) Create(ctx context.Context, feedID uuid.UUID, title, description, url, author string, publishedAt time.Time) error {
	desc := sql.NullString{String: description, Valid: true}
	if description == "" {
		desc = sql.NullString{String: "", Valid: false}
	}
	auth := sql.NullString{String: author, Valid: true}
	if author == "" {
		auth = sql.NullString{String: "", Valid: false}
	}
	return r.queries.CreateFeedPost(ctx, db.CreateFeedPostParams{
		FeedID:      feedID,
		Title:       title,
		Description: desc,
		Url:        url,
		PublishedAt: publishedAt,
		Author:     auth,
	})
}


func (r *DBFeedPostRepository) GetFeedPosts(ctx context.Context, feedURL string) ([]db.FeedPost, error) {
	posts, err := r.queries.GetFeedPosts(ctx, feedURL)
	if err != nil {
		return nil, err
	}
	
	var feedPosts []db.FeedPost
	for _, row := range posts {
		feedPosts = append(feedPosts, row.FeedPost)
	}
	return feedPosts, nil
}

