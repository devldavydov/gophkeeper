package server

import (
	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"github.com/devldavydov/gophkeeper/internal/server/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type GrpcServerSuite struct {
	suite.Suite
	testSrv   *GrpcServer
	testClt   pb.GophKeeperServiceClient
	stg       storage.Storage
	logger    *logrus.Logger
	fTeardown func()
}
