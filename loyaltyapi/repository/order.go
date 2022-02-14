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
	GetOrders(ctx context.Context, user *models.User) ([]*models.Order, error)
	CreateOrder(ctx context.Context, order *models.Order) (sql.Result, error)
	SaveOrder(ctx context.Context, order *models.Order) (sql.Result, error)
}

func (r *Repo) GetOrder(ctx context.Context, ID string) (*models.Order, error) {
	return nil, errors.New("not implemented")
}
func (r *Repo) GetOrderByNumber(ctx context.Context, number string) (*models.Order, error) {
	return nil, errors.New("not implemented")
}
func (r *Repo) GetOrders(ctx context.Context, user *models.User) ([]*models.Order, error) {
	return nil, errors.New("not implemented")
}
func (r *Repo) CreateOrder(ctx context.Context, order *models.Order) (sql.Result, error) {
	return nil, errors.New("not implemented")
}
func (r *Repo) SaveOrder(ctx context.Context, order *models.Order) (sql.Result, error) {
	return nil, errors.New("not implemented")
}
