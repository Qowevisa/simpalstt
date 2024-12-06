package grpc

import (
	"errors"
	"fmt"
	"graphql_api/env"
	"os"
	"sync"
	pb "worker"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrConnectionWasNotInitialized = errors.New("Connection to gRPC Storage server was not initialized")
)

type GRPCConnection struct {
	Conn     *grpc.ClientConn
	Client   pb.StorageClient
	ClientMu sync.Mutex
}

func Connect() (*GRPCConnection, error) {
	gRPCServerURL := os.Getenv(env.ENV_GRPC_SERVER_URL)
	conn, err := grpc.NewClient(gRPCServerURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// TODO: add error checks with errors.Is()
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}
	c := pb.NewStorageClient(conn)
	savedConn := &GRPCConnection{
		Conn:   conn,
		Client: c,
	}
	return savedConn, nil
}

// Should return *grpc.gRPCConnection and nil if connection was established
// Otherwise returns nil and grpc.ErrConnectionWasNotInitialized
// You can retry by calling grpc.Conenct() again
// As GRPCConnection.Conn is *grpc.ClientConn please don't forget about
// `defer g.Conn.Close()`
func SafeGetConnection() (*GRPCConnection, error) {
	conn, err := Connect()
	if err != nil {
		return nil, fmt.Errorf("Connect: %w", err)
	}
	if conn == nil {
		return nil, ErrConnectionWasNotInitialized
	}
	return conn, nil
}
