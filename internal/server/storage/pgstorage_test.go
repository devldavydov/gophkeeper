package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/devldavydov/gophkeeper/internal/common/model"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var logger = logrus.New() //nolint:gochecknoglobals // OK

const (
	_envTestDatabaseDsn = "TEST_DATABASE_DSN"
	_testDBTimeout      = 15 * time.Second
)

type PgStorageSuite struct {
	suite.Suite
	stg *PgStorage
}

func (pg *PgStorageSuite) SetupSuite() {
	var err error
	pg.stg, err = NewPgStorage(os.Getenv(_envTestDatabaseDsn), logger)
	require.NoError(pg.T(), err)
}

func (pg *PgStorageSuite) TearDownSuite() {
	pg.stg.Close()
}

func (pg *PgStorageSuite) TestPing() {
	ctx, cancel := context.WithTimeout(context.Background(), _testDBTimeout)
	defer cancel()

	pg.True(pg.stg.Ping(ctx))
}

func (pg *PgStorageSuite) TestCreateUser() {
	userName, userPassword := uuid.NewString(), uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), _testDBTimeout)
	defer cancel()

	pg.Run("create new user", func() {
		_, err := pg.stg.CreateUser(ctx, userName, userPassword)
		pg.NoError(err)
	})

	pg.Run("create same user twice", func() {
		_, err := pg.stg.CreateUser(ctx, userName, userPassword)
		pg.ErrorIs(err, ErrUserAlreadyExists)
	})
}

func (pg *PgStorageSuite) TestFindUser() {
	var userID int64

	userName, userPassword := uuid.NewString(), uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), _testDBTimeout)
	defer cancel()

	pg.Run("get not exist user", func() {
		_, _, err := pg.stg.FindUser(ctx, userName)
		pg.ErrorIs(err, ErrUserNotFound)
	})

	pg.Run("create new user", func() {
		var err error
		userID, err = pg.stg.CreateUser(ctx, userName, userPassword)
		pg.NoError(err)
	})

	pg.Run("get user", func() {
		uID, uPass, err := pg.stg.FindUser(ctx, userName)
		pg.NoError(err)
		pg.Equal(userID, uID)
		pg.Equal(userPassword, uPass)
	})
}

func (pg *PgStorageSuite) TestCreateSecret() {
	var userID int64
	var err error

	userName, userPassword := uuid.NewString(), uuid.NewString()
	secretName := uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), _testDBTimeout)
	defer cancel()

	pg.Run("create user", func() {
		userID, err = pg.stg.CreateUser(ctx, userName, userPassword)
		pg.NoError(err)
	})

	secret := &model.Secret{
		Type:       model.CredsSecret,
		Name:       secretName,
		Meta:       "test descr",
		Version:    1,
		PayloadRaw: []byte("test byte data"),
	}

	pg.Run("create secret", func() {
		err = pg.stg.CreateSecret(ctx, userID, secret)
		pg.NoError(err)
	})

	pg.Run("create secret twice", func() {
		err = pg.stg.CreateSecret(ctx, userID, secret)
		pg.ErrorIs(err, ErrSecretAlreadyExists)
	})
}

func (pg *PgStorageSuite) TestGetSecret() {
	var userID int64
	var err error

	userName, userPassword := uuid.NewString(), uuid.NewString()
	secretName := uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), _testDBTimeout)
	defer cancel()

	pg.Run("create user", func() {
		userID, err = pg.stg.CreateUser(ctx, userName, userPassword)
		pg.NoError(err)
	})

	pg.Run("get not exists secret", func() {
		_, err = pg.stg.GetSecret(ctx, userID, secretName)
		pg.ErrorIs(err, ErrSecretNotFound)
	})

	expSecret := &model.Secret{
		Type:       model.CredsSecret,
		Name:       secretName,
		Meta:       "test descr",
		Version:    1,
		PayloadRaw: []byte("test byte data"),
	}

	pg.Run("create secret", func() {
		err = pg.stg.CreateSecret(ctx, userID, expSecret)
		pg.NoError(err)
	})

	pg.Run("get secret", func() {
		var secret *model.Secret
		secret, err = pg.stg.GetSecret(ctx, userID, secretName)
		pg.NoError(err)
		pg.Equal(expSecret, secret)
	})
}

func (pg *PgStorageSuite) TestGetAllSecrets() {
	var userID int64
	var err error

	userName, userPassword := uuid.NewString(), uuid.NewString()
	secretName1, secretName2 := "b"+uuid.NewString(), "a"+uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), _testDBTimeout)
	defer cancel()

	pg.Run("create user", func() {
		userID, err = pg.stg.CreateUser(ctx, userName, userPassword)
		pg.NoError(err)
	})

	pg.Run("clear all secrets", func() {
		pg.NoError(pg.stg.DeleteAllSecrets(ctx))
	})

	pg.Run("get all secrets - no secrets", func() {
		_, err = pg.stg.GetAllSecrets(ctx, userID)
		pg.ErrorIs(err, ErrNoSecrets)
	})

	pg.Run("create secrets", func() {
		for _, name := range []string{secretName1, secretName2} {
			pg.NoError(
				pg.stg.CreateSecret(ctx, userID, &model.Secret{
					Type:       model.CredsSecret,
					Name:       name,
					Meta:       "test",
					Version:    1,
					PayloadRaw: []byte("123"),
				}))
		}
	})

	pg.Run("get all secrets", func() {
		var lst []model.SecretInfo
		lst, err = pg.stg.GetAllSecrets(ctx, userID)
		pg.Equal(2, len(lst))
		pg.Equal(secretName2, lst[0].Name)
		pg.Equal(secretName1, lst[1].Name)
	})
}

func (pg *PgStorageSuite) TestDeleteSecret() {
	var userID int64
	var err error

	userName, userPassword := uuid.NewString(), uuid.NewString()
	secretName := uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), _testDBTimeout)
	defer cancel()

	pg.Run("create user", func() {
		userID, err = pg.stg.CreateUser(ctx, userName, userPassword)
		pg.NoError(err)
	})

	expSecret := &model.Secret{
		Type:       model.CredsSecret,
		Name:       secretName,
		Meta:       "test descr",
		Version:    1,
		PayloadRaw: []byte("test byte data"),
	}

	pg.Run("create secret", func() {
		err = pg.stg.CreateSecret(ctx, userID, expSecret)
		pg.NoError(err)
	})

	pg.Run("get secret", func() {
		var secret *model.Secret
		secret, err = pg.stg.GetSecret(ctx, userID, secretName)
		pg.NoError(err)
		pg.Equal(expSecret, secret)
	})

	pg.Run("delete secret", func() {
		pg.NoError(pg.stg.DeleteSecret(ctx, userID, secretName))
	})

	pg.Run("get not exists secret", func() {
		_, err = pg.stg.GetSecret(ctx, userID, secretName)
		pg.ErrorIs(err, ErrSecretNotFound)
	})
}

func TestPgStorageSuite(t *testing.T) {
	// TODO: REMOVE AFTER ALL TESTS
	t.Setenv(_envTestDatabaseDsn, "postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable")
	//
	_, ok := os.LookupEnv(_envTestDatabaseDsn)
	if !ok {
		t.Skip("Test environment not set")
		return
	}
	suite.Run(t, new(PgStorageSuite))
}

func TestPgStorageCreateError(t *testing.T) {
	_, err := NewPgStorage("FooBar", logger)
	assert.Error(t, err)
}
