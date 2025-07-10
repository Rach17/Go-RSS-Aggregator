package service

import (
	"context"
	"github.com/Rach17/Go-RSS-Aggregator/repository"
	"github.com/Rach17/Go-RSS-Aggregator/db"
	"github.com/Rach17/Go-RSS-Aggregator/utils"
)


type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		Repo: repo,
	}
}

func (s *UserService) CreateUser(context context.Context, username, password string) (db.User, error) {
	hashedPassword, err := utils.Hash(password)
	if err != nil {
		return db.User{}, err
	}
	return s.Repo.CreateUser(context, username, hashedPassword)
}

