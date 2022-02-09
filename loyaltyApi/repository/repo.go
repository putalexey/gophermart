package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/pressly/goose/v3"
	"github.com/putalexey/gophermart/loyaltyApi/models"
)

type Repositorier interface {
	UserRepository
	OrderRepository
}

type UserRepository interface {
	GetUser(ctx context.Context, ID string) (*models.User, error)
	FindUserByLogin(ctx context.Context, Login string) (*models.User, error)
}

type OrderRepository interface {
	GetOrder(ctx context.Context, ID string) (*models.User, error)
	GetOrders(ctx context.Context, user *models.User) (*models.Order, error)
}

type Repo struct {
	DB *sql.DB
}

func New(db *sql.DB, migrationsDir string) (*Repo, error) {
	repo := &Repo{DB: db}

	//migrate
	if migrationsDir != "" {
		err := goose.Up(db, migrationsDir)
		if err != nil {
			return nil, err
		}
	}

	return repo, nil
}

func (r *Repo) GetUser(ctx context.Context, ID string) (*models.User, error) {
	return nil, errors.New("User not found")
}

func (r *Repo) FindUserByLogin(ctx context.Context, Login string) (*models.User, error) {
	return nil, nil
}
