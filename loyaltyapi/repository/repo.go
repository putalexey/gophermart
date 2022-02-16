package repository

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"time"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrNotEnoughBalance = errors.New("not enough balance")
)

type Repositorier interface {
	UserRepository
	OrderRepository
}

type Repo struct {
	DatabaseDSN string
	db          *sqlx.DB
}

func New(DatabaseDSN string, migrationsDir string) (*Repo, error) {
	repo := &Repo{DatabaseDSN: DatabaseDSN}
	err := repo.connectToDB()
	if err != nil {
		return nil, err
	}

	//migrate
	if migrationsDir != "" {
		err := goose.Up(repo.db.DB, migrationsDir)
		if err != nil {
			return nil, err
		}
	}

	return repo, nil
}

func (r *Repo) connectToDB() error {
	db, err := sqlx.Open("pgx", r.DatabaseDSN)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(30 * time.Second)
	db.SetConnMaxLifetime(2 * time.Minute)

	r.db = db
	err = r.db.Ping()
	if err != nil {
		return err
	}

	return nil
}
