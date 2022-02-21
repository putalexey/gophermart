package repository

import (
	"context"
	"database/sql"
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
	query := `INSERT INTO "jobs" ("uuid", "order_uuid", "proceed_at", "tries") VALUES (:uuid, :order_uuid, :proceed_at, :tries)`
	_, err := r.db.NamedExecContext(ctx, query, job)
	return err
}

// TakeJob query job to proceed and sets "proceed_at" to now + ttl value
func (r *Repo) TakeJob(ctx context.Context, ttl time.Duration) (*models.Job, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	job := &models.Job{}
	query := `SELECT "uuid", "order_uuid", "proceed_at", "tries" FROM "jobs" WHERE tries < 5 and proceed_at <= $1 order by proceed_at LIMIT 1 FOR UPDATE SKIP LOCKED`
	err = tx.GetContext(ctx, job, query, time.Now())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	job.ProceedAt = time.Now().Add(ttl)
	job.Tries = job.Tries + 1
	_, err = tx.NamedExecContext(ctx, `UPDATE "jobs" SET "proceed_at" = :proceed_at, "tries" = :tries WHERE "uuid" = :uuid`, job)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return job, nil
}

func (r *Repo) UpdateJob(ctx context.Context, job *models.Job) error {
	query := `UPDATE "jobs" SET "proceed_at" = :proceed_at, "tries" = :tries WHERE "uuid" = :uuid`
	_, err := r.db.NamedExecContext(ctx, query, job)
	return err
}

func (r *Repo) DeleteJob(ctx context.Context, job *models.Job) error {
	query := `DELETE FROM "jobs" WHERE "uuid" = :uuid`
	_, err := r.db.NamedExecContext(ctx, query, job)
	return err
}
