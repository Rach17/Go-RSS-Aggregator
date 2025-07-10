package service

import (
	"context"
	"github.com/Rach17/Go-RSS-Aggregator/repository"
	"github.com/Rach17/Go-RSS-Aggregator/db"
)

type AuthService struct {
	Repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{
		Repo: repo,
	}
}

func (s *AuthService) IsAuth(ctx context.Context, apiKey string) (db.User, error) {
	user, err := s.Repo.GetUserByAPIKey(ctx, apiKey)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}
