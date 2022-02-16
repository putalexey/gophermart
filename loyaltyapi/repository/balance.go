package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/putalexey/gophermart/loyaltyapi/models"
)

type BalanceRepository interface {
	GetUserBalance(ctx context.Context, userUUID string) (*models.Balance, error)
	BalanceWithdraw(ctx context.Context, withdrawal *models.Withdrawal) (sql.Result, error)
	GetBalanceWithdrawals(ctx context.Context, user *models.User) ([]models.Withdrawal, error)
}

func (r *Repo) CreateBalance(ctx context.Context, balance *models.Balance) (sql.Result, error) {
	return r.db.NamedExecContext(
		ctx,
		"INSERT INTO balances (user_uuid, current, withdrawn) VALUES (:user_uuid, :current, :withdrawn)",
		balance,
	)
}

func (r *Repo) createBalanceTx(tx *sqlx.Tx, ctx context.Context, balance *models.Balance) (sql.Result, error) {
	return tx.NamedExecContext(
		ctx,
		"INSERT INTO balances (user_uuid, current, withdrawn) VALUES (:user_uuid, :current, :withdrawn)",
		balance,
	)
}

func (r *Repo) SaveBalance(ctx context.Context, balance *models.Balance) (sql.Result, error) {
	return r.db.NamedExecContext(
		ctx,
		"UPDATE balances SET current=:current, withdrawn=:withdrawn WHERE user_uuid=:user_uuid",
		balance,
	)
}

func (r *Repo) saveBalanceTx(tx *sqlx.Tx, ctx context.Context, balance *models.Balance) (sql.Result, error) {
	return tx.NamedExecContext(
		ctx,
		"UPDATE balances SET current=:current, withdrawn=:withdrawn WHERE user_uuid=:user_uuid",
		balance,
	)
}

func (r *Repo) GetUserBalance(ctx context.Context, userUUID string) (*models.Balance, error) {
	var balance = &models.Balance{}
	err := r.db.GetContext(ctx, balance, "SELECT user_uuid, current, withdrawn FROM balances WHERE user_uuid = $1", userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return balance, nil
}

func (r *Repo) getUserBalanceTxForUpdate(tx *sqlx.Tx, ctx context.Context, userUUID string) (*models.Balance, error) {
	var balance = &models.Balance{}
	err := tx.GetContext(ctx, balance, "SELECT user_uuid, current, withdrawn FROM balances WHERE user_uuid = $1 FOR UPDATE", userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return balance, nil
}

func (r *Repo) BalanceWithdraw(ctx context.Context, withdrawal *models.Withdrawal) (sql.Result, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	balance, err := r.getUserBalanceTxForUpdate(tx, ctx, withdrawal.UserUUID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if balance.Current < withdrawal.Sum {
		tx.Rollback()
		return nil, ErrNotEnoughBalance
	}

	balance.Current -= withdrawal.Sum
	balance.Withdrawn += withdrawal.Sum
	_, err = r.saveBalanceTx(tx, ctx, balance)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	query := `INSERT INTO withdrawals ("uuid", "user_uuid", "order", "sum", "processed_at") VALUES (:uuid, :user_uuid, :order, :sum, :processed_at)`
	result, err := tx.NamedExecContext(
		ctx,
		query,
		withdrawal,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return result, err
}

func (r *Repo) GetBalanceWithdrawals(ctx context.Context, user *models.User) ([]models.Withdrawal, error) {
	withdrawals := []models.Withdrawal{}
	query := `SELECT "uuid", "user_uuid", "order", "sum", "processed_at" FROM "withdrawals" WHERE "user_uuid" = $1`
	err := r.db.SelectContext(ctx, &withdrawals, query, user.UUID)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}
