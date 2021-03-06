package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/putalexey/gophermart/loyaltyapi/models"
)

type UserRepository interface {
	GetUser(ctx context.Context, ID string) (*models.User, error)
	FindUserByLogin(ctx context.Context, Login string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (sql.Result, error)
	SaveUser(ctx context.Context, user *models.User) (sql.Result, error)
}

func (r *Repo) GetUser(ctx context.Context, id string) (*models.User, error) {
	var user = &models.User{}
	err := r.req.GetContext(ctx, user, "SELECT uuid, login, password FROM users WHERE uuid = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *Repo) FindUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user = &models.User{}
	err := r.req.GetContext(ctx, user, "SELECT uuid, login, password FROM users WHERE login like $1", login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *Repo) CreateUser(ctx context.Context, user *models.User) (sql.Result, error) {
	rtx, err := r.Begin()
	if err != nil {
		return nil, err
	}
	defer rtx.Rollback()

	result, err := rtx.req.NamedExecContext(
		ctx,
		"INSERT INTO users (uuid, login, password) VALUES (:uuid, :login, :password)",
		user,
	)
	if err != nil {
		return nil, err
	}

	balance := &models.Balance{
		UserUUID:  user.UUID,
		Current:   0,
		Withdrawn: 0,
	}

	_, err = rtx.CreateBalance(ctx, balance)
	if err != nil {
		return nil, err
	}

	if err = rtx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repo) SaveUser(ctx context.Context, user *models.User) (sql.Result, error) {
	return r.req.NamedExecContext(
		ctx,
		"UPDATE users SET login=:login, password=:password WHERE uuid=:uuid",
		user,
	)
}
