package repository

import (
	"context"
	"database/sql"
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
	RepoTxStarter
	UserRepository
	OrderRepository
	BalanceRepository
	JobRepository
}

type RepoTxStarter interface {
	Begin() (*RepoTx, error)
}

type dbRequester interface {
	sqlx.QueryerContext
	sqlx.ExecerContext
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Repo struct {
	DatabaseDSN string
	db          *sqlx.DB
	req         dbRequester
}

type RepoTx struct {
	*sqlx.Tx
	Repo
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
	r.req = db
	err = r.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Begin() (*RepoTx, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	r2 := *r
	r2.req = tx

	rtx := &RepoTx{
		Tx:   tx,
		Repo: r2,
	}
	return rtx, nil
}
func (rtx *RepoTx) Commit() error {
	return rtx.Tx.Commit()
}
func (rtx *RepoTx) Rollback() error {
	return rtx.Tx.Rollback()
}

func (rtx *RepoTx) Begin() (*RepoTx, error) {
	return nil, errors.New("transaction already started")
}
