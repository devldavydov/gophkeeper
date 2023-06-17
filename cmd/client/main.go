// Package represents main application for client.
package main

import (
	"context"
	"fmt"
	"log"

	gkTLS "github.com/devldavydov/gophkeeper/internal/common/tls"
	"github.com/devldavydov/gophkeeper/internal/common/token"
	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error { //nolint:unparam // REMOVE
	// Test code
	tlsCredentials, _ := gkTLS.LoadCACert("/home/devldavydov/go/gophkeeper/tls/ca-cert.pem", "127.0.0.1")
	conn, _ := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(tlsCredentials))
	clnt := pb.NewGophKeeperServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tokenString, err := clnt.UserLogin(ctx, &pb.User{Login: "foo", Password: "bar"})
	fmt.Println(tokenString)
	fmt.Println(err)

	md := metadata.New(map[string]string{token.HeaderName: tokenString.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	_, err = clnt.Ping(ctx, &pb.Empty{})
	fmt.Println(err)
	//

	return nil
}
