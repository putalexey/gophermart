package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrOrderNotFound = errors.New("order not found")
)

type Repositorier interface {
	UserRepository
	OrderRepository
}

type Repo struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB, migrationsDir string) (*Repo, error) {
	repo := &Repo{DB: db}

	//migrate
	if migrationsDir != "" {
		err := goose.Up(db.DB, migrationsDir)
		if err != nil {
			return nil, err
		}
	}

	return repo, nil
}
