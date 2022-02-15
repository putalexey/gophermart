package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/putalexey/gophermart/loyaltyapi/models"
)

type OrderRepository interface {
	GetOrder(ctx context.Context, ID string) (*models.Order, error)
	GetOrderByNumber(ctx context.Context, number string) (*models.Order, error)
	GetUserOrders(ctx context.Context, user *models.User) ([]models.Order, error)
	CreateOrder(ctx context.Context, order *models.Order) (sql.Result, error)
	SaveOrder(ctx context.Context, order *models.Order) (sql.Result, error)
}

func (r *Repo) GetOrder(ctx context.Context, uuid string) (*models.Order, error) {
	order := &models.Order{}
	query := `SELECT uuid, user_uuid, number, status, accrual, uploaded_at FROM orders WHERE uuid=$1 LIMIT 1`
	err := r.db.GetContext(ctx, order, query, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return order, nil
}

func (r *Repo) GetOrderByNumber(ctx context.Context, number string) (*models.Order, error) {
	order := &models.Order{}
	query := `SELECT uuid, user_uuid, number, status, accrual, uploaded_at FROM orders WHERE number=$1 LIMIT 1`
	err := r.db.GetContext(ctx, order, query, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return order, nil
}

func (r *Repo) GetUserOrders(ctx context.Context, user *models.User) ([]models.Order, error) {
	//goland:noinspection GoPreferNilSlice
	orders := []models.Order{}
	query := "SELECT uuid, user_uuid, number, status, accrual, uploaded_at FROM orders WHERE user_uuid=$1 order by uploaded_at"
	err := r.db.SelectContext(ctx, &orders, query, user.UUID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *Repo) CreateOrder(ctx context.Context, order *models.Order) (sql.Result, error) {
	query := `INSERT INTO orders (uuid, user_uuid, number, status, accrual, uploaded_at) VALUES (:uuid, :user_uuid, :number, :status, :accrual, :uploaded_at)`
	return r.db.NamedExecContext(ctx, query, order)
}

// SaveOrder accrual is change by another function
func (r *Repo) SaveOrder(ctx context.Context, order *models.Order) (sql.Result, error) {
	query := `UPDATE orders SET user_uuid=:user_uuid, number=:number, status=:status, accrual=:accrual WHERE uuid=:uuid`
	return r.db.NamedExecContext(ctx, query, order)
}
