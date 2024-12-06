package grpc

import (
	"errors"
	"fmt"
	"os"
	"sync"
	pb "worker"
	"worker_client/env"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrConnectionWasNotInitialized = errors.New("Connection to gRPC Storage server was not initialized")
)

type GRPC_Worker_ClientConnection struct {
	Conn     *grpc.ClientConn
	Client   pb.WorkerClient
	ClientMu sync.Mutex
}

func ConnectToWorkerServer() (*GRPC_Worker_ClientConnection, error) {
	gRPCServerURL := os.Getenv(env.ENV_GRPC_SERVER_URL)
	conn, err := grpc.NewClient(gRPCServerURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// TODO: add error checks with errors.Is()
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}
	c := pb.NewWorkerClient(conn)
	savedConn := &GRPC_Worker_ClientConnection{
		Conn:   conn,
		Client: c,
	}
	return savedConn, nil
}

// Should return *grpc.GRPC_Worker_ClientConnection and nil if connection
// was established
// Otherwise returns nil and grpc.ErrConnectionWasNotInitialized
// You can retry by calling grpc.ConenctToWorkerServer() again
// As GRPC_Worker_ClientConnection.Conn is *grpc.ClientConn please don't forget
// about `defer g.Conn.Close()`
func SafeGetConnectionToWorkerServer() (*GRPC_Worker_ClientConnection, error) {
	conn, err := ConnectToWorkerServer()
	if err != nil {
		return nil, fmt.Errorf("Connect: %w", err)
	}
	if conn == nil {
		return nil, ErrConnectionWasNotInitialized
	}
	return conn, nil
}
