package repository

import (
	"context"
	"errors"
	"github.com/putalexey/gophermart/loyaltyapi/models"
	"time"
)

type JobRepository interface {
	CreateJob(ctx context.Context, job *models.Job) error
	TakeJob(ctx context.Context, ttl time.Duration) (*models.Job, error)
	UpdateJob(ctx context.Context, job *models.Job) error
	DeleteJob(ctx context.Context, job *models.Job) error
}

func (r *Repo) CreateJob(ctx context.Context, job *models.Job) error {
	return errors.New("not implemented")
}

func (r *Repo) TakeJob(ctx context.Context, ttl time.Duration) (*models.Job, error) {
	return nil, errors.New("not implemented")
}

func (r *Repo) UpdateJob(ctx context.Context, job *models.Job) error {
	return errors.New("not implemented")
}

func (r *Repo) DeleteJob(ctx context.Context, job *models.Job) error {
	return errors.New("not implemented")
}
