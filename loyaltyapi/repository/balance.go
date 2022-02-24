package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/putalexey/gophermart/loyaltyapi/models"
)

type BalanceRepository interface {
	GetUserBalance(ctx context.Context, userUUID string) (*models.Balance, error)
	BalanceDeposit(ctx context.Context, deposit *models.Deposit) (sql.Result, error)
	BalanceWithdraw(ctx context.Context, withdrawal *models.Withdrawal) (sql.Result, error)
	GetBalanceWithdrawals(ctx context.Context, user *models.User) ([]models.Withdrawal, error)
}

func (r *Repo) CreateBalance(ctx context.Context, balance *models.Balance) (sql.Result, error) {
	return r.req.NamedExecContext(
		ctx,
		"INSERT INTO balances (user_uuid, current, withdrawn) VALUES (:user_uuid, :current, :withdrawn)",
		balance,
	)
}

func (r *Repo) SaveBalance(ctx context.Context, balance *models.Balance) (sql.Result, error) {
	return r.req.NamedExecContext(
		ctx,
		"UPDATE balances SET current=:current, withdrawn=:withdrawn WHERE user_uuid=:user_uuid",
		balance,
	)
}

func (r *Repo) GetUserBalance(ctx context.Context, userUUID string) (*models.Balance, error) {
	var balance = &models.Balance{}
	err := r.req.GetContext(ctx, balance, "SELECT user_uuid, current, withdrawn FROM balances WHERE user_uuid = $1", userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return balance, nil
}

func (r *Repo) GetUserBalanceForUpdate(ctx context.Context, userUUID string) (*models.Balance, error) {
	var balance = &models.Balance{}
	err := r.req.GetContext(ctx, balance, "SELECT user_uuid, current, withdrawn FROM balances WHERE user_uuid = $1 FOR UPDATE", userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return balance, nil
}

func (r *Repo) BalanceWithdraw(ctx context.Context, withdrawal *models.Withdrawal) (sql.Result, error) {
	rtx, err := r.Begin()
	if err != nil {
		return nil, err
	}
	defer rtx.Rollback()

	balance, err := rtx.GetUserBalanceForUpdate(ctx, withdrawal.UserUUID)
	if err != nil {
		return nil, err
	}
	if balance.Current < withdrawal.Sum {
		return nil, ErrNotEnoughBalance
	}

	balance.Current -= withdrawal.Sum
	balance.Withdrawn += withdrawal.Sum
	_, err = rtx.SaveBalance(ctx, balance)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO withdrawals ("uuid", "user_uuid", "order", "sum", "processed_at") VALUES (:uuid, :user_uuid, :order, :sum, :processed_at)`
	result, err := rtx.req.NamedExecContext(
		ctx,
		query,
		withdrawal,
	)
	if err != nil {
		return nil, err
	}

	if err = rtx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repo) BalanceDeposit(ctx context.Context, deposit *models.Deposit) (sql.Result, error) {
	rtx, err := r.Begin()
	if err != nil {
		return nil, err
	}
	defer rtx.Rollback()

	balance, err := rtx.GetUserBalanceForUpdate(ctx, deposit.UserUUID)
	if err != nil {
		return nil, err
	}

	balance.Current += deposit.Sum
	result, err := rtx.SaveBalance(ctx, balance)
	if err != nil {
		return nil, err
	}

	if err = rtx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func (rtx *RepoTx) BalanceDeposit(ctx context.Context, deposit *models.Deposit) (sql.Result, error) {
	balance, err := rtx.GetUserBalanceForUpdate(ctx, deposit.UserUUID)
	if err != nil {
		return nil, err
	}

	balance.Current += deposit.Sum
	result, err := rtx.SaveBalance(ctx, balance)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repo) GetBalanceWithdrawals(ctx context.Context, user *models.User) ([]models.Withdrawal, error) {
	withdrawals := []models.Withdrawal{}
	query := `SELECT "uuid", "user_uuid", "order", "sum", "processed_at" FROM "withdrawals" WHERE "user_uuid" = $1`
	err := r.req.SelectContext(ctx, &withdrawals, query, user.UUID)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}
