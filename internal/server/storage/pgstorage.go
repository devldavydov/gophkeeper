package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/devldavydov/gophkeeper/internal/common/model"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	_databaseInitTimeout time.Duration = 10 * time.Second

	_constraintUniqueViolation pq.ErrorCode = "23505"
	_constraintUsernameCheck   string       = "users_username_key"
	_constraintSecretCheck     string       = "secrets_pkey"
)

var (
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrSecretAlreadyExists = errors.New("secret already exists")
	ErrSecretNotFound      = errors.New("secret not found")
	ErrNoSecrets           = errors.New("no secrets")
	ErrSecretOutdated      = errors.New("secret outdated")
	ErrSecretWrongVersion  = errors.New("secret wrong version")
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

// CreateUser - creates new user in storage. Accepts login and password.
//
// Returns new user id or error:
//
// - ErrUserAlreadyExists - if user already exists.
//
// - internal PG error.
func (pg *PgStorage) CreateUser(ctx context.Context, login, password string) (int64, error) {
	var userID int64
	err := pg.db.QueryRowContext(ctx, _sqlCreateUser, login, password).Scan(&userID)
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

	return userID, nil
}

// FindUser - finds user in storage. Accepts login.
//
// Returns user id or error:
//
// - ErrUserNotFound - if user not exists.
//
// - internal PG error.
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

// CreateSecret - create user's secret in storage. Accepts userID and secret.
//
// Returns nil or error:
//
// - ErrSecretAlreadyExists - if secret already exists.
//
// - internal PG error.
func (pg *PgStorage) CreateSecret(ctx context.Context, userID int64, secret *model.Secret) error {
	_, err := pg.db.ExecContext(
		ctx,
		_sqlCreateSecret,
		userID,
		secret.Type,
		secret.Name,
		secret.Meta,
		secret.Version,
		secret.PayloadRaw,
	)
	if err != nil {
		var pqErr *pq.Error
		if !errors.As(err, &pqErr) {
			return err
		}

		if pqErr.Code == _constraintUniqueViolation && pqErr.Constraint == _constraintSecretCheck {
			return ErrSecretAlreadyExists
		}

		return err
	}

	return nil
}

// GetSecret - gets user's secret from storage. Accepts userID and secret name.
//
// Returns secret or error:
//
// - ErrSecretNotFound - if secret not exists.
//
// - internal PG error.
func (pg *PgStorage) GetSecret(ctx context.Context, userID int64, name string) (*model.Secret, error) {
	secret := &model.Secret{}
	err := pg.db.
		QueryRowContext(ctx, _sqlGetSecret, userID, name).
		Scan(&secret.Type, &secret.Name, &secret.Meta, &secret.Version, &secret.PayloadRaw)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrSecretNotFound
	case err != nil:
		return nil, err
	}

	return secret, nil
}

// GetAllSecrets - gets all user's secrets from storage. Accepts userID.
//
// Returns secrets or error:
//
// - ErrNoSecrets - if no secrets exist.
//
// - internal PG error.
func (pg *PgStorage) GetAllSecrets(ctx context.Context, userID int64) ([]model.SecretInfo, error) {
	rows, err := pg.db.QueryContext(ctx, _sqlGetAllSecrets, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.SecretInfo

	for rows.Next() {
		secretItem := model.SecretInfo{}
		err = rows.Scan(&secretItem.Type, &secretItem.Name, &secretItem.Version)
		if err != nil {
			return nil, err
		}
		items = append(items, secretItem)
	}

	if len(items) == 0 {
		return nil, ErrNoSecrets
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// DeleteSecret - deletes user's secret from storage. Accepts userID and secret name.
//
// Returns nil or error:
//
// - internal PG error.
func (pg *PgStorage) DeleteSecret(ctx context.Context, userID int64, name string) error {
	_, err := pg.db.ExecContext(ctx, _sqlDeleteSecret, userID, name)
	return err
}

// DeleteSecret - deletes all user's secrets from storage. Accepts userID.
//
// Returns nil or error:
//
// - internal PG error.
func (pg *PgStorage) DeleteAllSecrets(ctx context.Context, userID int64) error {
	_, err := pg.db.ExecContext(ctx, _sqlDeleteAllSecrets, userID)
	return err
}

// UpdateSecret - updates user's secret in storage. Accepts userIDm secret name and secret update.
//
// Returns nil or error:
//
// - ErrSecretNotFound - if secret not found in storage.
//
// - ErrSecretOutdated - if secret in storage has greater version than update.
//
// - ErrSecretWrongVersion - if diff between storage secret version and update version greater than 1.
//
// Correct update - when updateVersion - currentVersion = 1.
//
// - internal PG error.
func (pg *PgStorage) UpdateSecret(ctx context.Context, userID int64, name string, update *model.SecretUpdate) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Lock current secret
	var curVersion int64
	err = tx.QueryRowContext(ctx, _sqlLockSecret, userID, name).Scan(&curVersion)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrSecretNotFound
	case err != nil:
		return err
	}

	// Check version
	switch {
	case update.Version <= curVersion:
		return ErrSecretOutdated
	case update.Version-curVersion > 1:
		return ErrSecretWrongVersion
	}

	// Update
	if update.UpdatePayload {
		_, err = tx.ExecContext(ctx, _sqlUpdateSecret, userID, name, update.Meta, update.Version, update.PayloadRaw)
	} else {
		_, err = tx.ExecContext(ctx, _sqlUpdateSecretWithoutPayload, userID, name, update.Meta, update.Version)
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Ping - checks storage availability.
//
// Returns true or false in case of error.
func (pg *PgStorage) Ping(ctx context.Context) bool {
	if err := pg.db.PingContext(ctx); err != nil {
		pg.logger.Errorf("Failed to ping database, err: %v", err)
		return false
	}

	return true
}

// Close - closes connection with storage.
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

	for _, createTbl := range []string{_sqlCreateTableUser, _sqlCreateTableSecret} {
		_, err := pg.db.ExecContext(ctx, createTbl)
		if err != nil {
			return err
		}
	}

	return nil
}
