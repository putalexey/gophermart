package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/pressly/goose/v3"
	"github.com/putalexey/gophermart/loyaltyApi/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
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

func (r *Repo) GetUser(ctx context.Context, id string) (*models.User, error) {
	return nil, ErrUserNotFound
}

func (r *Repo) FindUserByLogin(ctx context.Context, login string) (*models.User, error) {
	if login == "test" {
		return &models.User{
			UUID:     "asdasdasd-as-da-sdasd",
			Login:    "test",
			Password: "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92", // 123456
		}, nil
	}
	return nil, ErrUserNotFound
}
