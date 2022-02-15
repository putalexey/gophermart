package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/putalexey/gophermart/loyaltyapi/models"
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
	CreateUser(ctx context.Context, user *models.User) (sql.Result, error)
	SaveUser(ctx context.Context, user *models.User) (sql.Result, error)
}

type OrderRepository interface {
	GetOrder(ctx context.Context, ID string) (*models.User, error)
	GetOrders(ctx context.Context, user *models.User) (*models.Order, error)
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

func (r *Repo) GetUser(ctx context.Context, id string) (*models.User, error) {
	var user = &models.User{}
	err := r.DB.GetContext(ctx, user, "SELECT uuid, login, password FROM users WHERE uuid = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *Repo) FindUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user = &models.User{}
	err := r.DB.GetContext(ctx, user, "SELECT uuid, login, password FROM users WHERE login like $1", login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *Repo) CreateUser(ctx context.Context, user *models.User) (sql.Result, error) {
	return r.DB.NamedExecContext(
		ctx,
		"INSERT INTO users (uuid, login, password) VALUES (:uuid, :login, :password)",
		user,
	)
}

func (r *Repo) SaveUser(ctx context.Context, user *models.User) (sql.Result, error) {
	return r.DB.NamedExecContext(
		ctx,
		"UPDATE users SET (login=:login, password=:password) WHERE uuid=:uuid",
		user,
	)
}
