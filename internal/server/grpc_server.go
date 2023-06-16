package server

import (
	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"github.com/devldavydov/gophkeeper/internal/server/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GrpcServer represents gRPC server.
type GrpcServer struct {
	pb.UnimplementedGophKeeperServiceServer
	stg    storage.Storage
	logger *logrus.Logger
}

// NewGrpcServer creates new GRPCServer object.
func NewGrpcServer(
	stg storage.Storage,
	tlsCredentials credentials.TransportCredentials,
	logger *logrus.Logger,
) (*grpc.Server, *GrpcServer) {
	opts := []grpc.ServerOption{
		// grpc.UnaryInterceptor(interceptor.NewTrustedSubnetInterceptor(trustedSubnet, []string{"/grpc.MetricService/UpdateMetrics"}).Handle),
	}

	opts = append([]grpc.ServerOption{grpc.Creds(tlsCredentials)}, opts...)

	grpcSrv := grpc.NewServer(opts...)
	srv := &GrpcServer{stg: stg, logger: logger}
	pb.RegisterGophKeeperServiceServer(grpcSrv, srv)
	return grpcSrv, srv
}
