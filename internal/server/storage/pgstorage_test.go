package storage

import (
	"context"
	"os"
	"testing"
	"time"

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

func (pg *PgStorageSuite) SetupTest() {
	var err error
	pg.stg, err = NewPgStorage(os.Getenv(_envTestDatabaseDsn), logger)
	require.NoError(pg.T(), err)
}

func (pg *PgStorageSuite) TearDownTest() {
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
	var userID int

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

func TestPgStorageSuite(t *testing.T) {
	// t.Setenv(_envTestDatabaseDsn, "postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable")
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
