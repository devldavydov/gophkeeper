package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	_databaseInitTimeout = 10 * time.Second
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

	for _, createTbl := range []string{} {
		_, err := pg.db.ExecContext(ctx, createTbl)
		if err != nil {
			return err
		}
	}

	return nil
}
