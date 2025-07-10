package repository

import (
	"context"
	"database/sql"
	"github.com/Rach17/Go-RSS-Aggregator/db"
)


type UserRepository interface {
	CreateUser(ctx context.Context, username, passwordhash string) (db.User, error)
	GetUserByAPIKey(ctx context.Context, apiKey string) (db.User, error)
}

type DBUserRepository struct {
	queries *db.Queries
	db      *sql.DB
}

func NewDBUserRepository(database *sql.DB) *DBUserRepository {
	return &DBUserRepository{
		queries: db.New(database),
		db:      database,
	}
}

func (r *DBUserRepository) CreateUser(ctx context.Context, username, passwordhash string) (db.User, error) {
	return r.queries.CreateUser(ctx, db.CreateUserParams{Username: username, PasswordHash: passwordhash})
}

func (r *DBUserRepository) GetUserByAPIKey(ctx context.Context, apiKey string) (db.User, error) {
	user, err := r.queries.GetUserByAPIKey(ctx, apiKey)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}