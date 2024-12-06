package env

import (
	"errors"
	"fmt"
	"os"
)

const (
	ENV_GRPC_SERVER_URL = "GRPC_SERVER_URL"
)

var (
	ErrEnvNotSetOrEmpty = errors.New("env not set or is empty")
	ErrGrpcServerEnv    = fmt.Errorf("%s: %w", ENV_GRPC_SERVER_URL, ErrEnvNotSetOrEmpty)
)

var (
	grpcServerUrl string
)

func Init() error {
	url, set := os.LookupEnv(ENV_GRPC_SERVER_URL)
	if !set || url == "" {
		return ErrGrpcServerEnv
	}
	grpcServerUrl = url
	return nil
}

// Basically returns os.GetEnv(ENV_GRPC_SERVER_URL)
// NOTE: It is ASSUMED that function env.Init() was called BEFORE this function
// otherwise you WILL GET EMPTY STRING
func GRPCServerURL() string {
	return grpcServerUrl
}
