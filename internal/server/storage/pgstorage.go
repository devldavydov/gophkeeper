package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	_databaseInitTimeout time.Duration = 10 * time.Second

	_constraintUniqueViolation pq.ErrorCode = "23505"
	_constraintUsernameCheck   string       = "users_username_key"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

// PgStorage is a Storage implementation for PostgreSQL database.
type PgStorage struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewPgStorage creates PgStorage new object.
func NewPgStorage(pgConnString string, logger *logrus.Logger) (*PgStorage, error) {
	db, err := sql.Open("postgres", pgConnString)
	if err != nil {
		return nil, err
	}

	pgstorage := &PgStorage{db: db, logger: logger}

	if err = pgstorage.init(); err != nil {
		return nil, err
	}

	return pgstorage, nil
}

var _ Storage = (*PgStorage)(nil)

func (pg *PgStorage) CreateUser(ctx context.Context, login, password string) (int64, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var userID int64
	err = tx.QueryRowContext(ctx, _sqlCreateUser, login, password).Scan(&userID)
	if err != nil {
		var pqErr *pq.Error
		if !errors.As(err, &pqErr) {
			return 0, err
		}

		if pqErr.Code == _constraintUniqueViolation && pqErr.Constraint == _constraintUsernameCheck {
			return 0, ErrUserAlreadyExists
		}

		return 0, err
	}

	return userID, tx.Commit()
}

func (pg *PgStorage) FindUser(ctx context.Context, login string) (int64, string, error) {
	var userID int64
	var userPassword string
	err := pg.db.QueryRowContext(ctx, _sqlFindUser, login).Scan(&userID, &userPassword)
	switch {
	case err == sql.ErrNoRows:
		return 0, "", ErrUserNotFound
	case err != nil:
		return 0, "", err
	}

	return userID, userPassword, nil
}

func (pg *PgStorage) Ping(ctx context.Context) bool {
	if err := pg.db.PingContext(ctx); err != nil {
		pg.logger.Errorf("Failed to ping database, err: %v", err)
		return false
	}

	return true
}

func (pg *PgStorage) Close() {
	if pg.db == nil {
		return
	}

	err := pg.db.Close()
	if err != nil {
		pg.logger.Errorf("Database conn close err: %v", err)
	}
}

func (pg *PgStorage) init() error {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseInitTimeout)
	defer cancel()

	for _, createTbl := range []string{_sqlCreateTableUser} {
		_, err := pg.db.ExecContext(ctx, createTbl)
		if err != nil {
			return err
		}
	}

	return nil
}
