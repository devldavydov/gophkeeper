package storage

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

var logger = logrus.New()

const _envTestDatabaseDsn = "TEST_DATABASE_DSN"

type PgStorageSuite struct {
	suite.Suite
	stg *PgStorage
}
