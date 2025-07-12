package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Rach17/Go-RSS-Aggregator/db"
	"github.com/Rach17/Go-RSS-Aggregator/repository"
	"github.com/google/uuid"
)

type FeedPostService struct {
	FeedRepo  repository.FeedRepository
	PostRepo  repository.FeedPostRepository
	FeedService *FeedService
}

func NewFeedPostService(feedRepository repository.FeedRepository, postRepository repository.FeedPostRepository) *FeedPostService {
	return &FeedPostService{
		FeedRepo:  feedRepository,
		PostRepo:  postRepository,
	}
}

func (s *FeedPostService) CreateFeedPost(ctx context.Context, feedUrl, title, description, url, author string, publishedAt time.Time) error {
	if title == "" || url == "" {
		return fmt.Errorf("title and url cannot be empty")
	}

	feed, err := s.FeedRepo.GetFeedByURL(ctx, feedUrl)
	if err != nil {
		return fmt.Errorf("feed not found: %w", err)
	}

	if feed.ID == uuid.Nil {
		return fmt.Errorf("feed does not exist, please create it first")
	}

	return s.PostRepo.Create(ctx, feed.ID, title, description, url, author, publishedAt)
}


func (s *FeedPostService) GetFeedPosts(ctx context.Context, feedURL string) ([]db.FeedPost, error) {
	posts, err := s.PostRepo.GetFeedPosts(ctx, feedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed posts: %w", err)
	}
	return posts, nil
}
