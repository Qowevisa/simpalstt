package env

import (
	"errors"
	"fmt"
	"os"
)

const (
	ENV_ELASTICSEARCH_URL = "ELASTICSEARCH_URL"
	ENV_GRPC_SERVER_URL   = "WORKER_SERVER_URL"
)

var (
	ErrEnvNotSetOrEmpty = errors.New("env not set or is empty")
	ErrGrpcServerEnv    = fmt.Errorf("%s: %w", ENV_GRPC_SERVER_URL, ErrEnvNotSetOrEmpty)
	ErrElasticSearchEnv = fmt.Errorf("%s: %w", ENV_ELASTICSEARCH_URL, ErrEnvNotSetOrEmpty)
)

var (
	grpcServerUrl    string
	elasticSearchUrl string
)

// Init changes internal private variables `grpcServerUrl` and `elasticSearchUrl`
// To get the grpcServerUrl value use env.GRPCServerURL function
// To get the elasticSearchUrl value use env.ElasticSearchURL function
func Init() error {
	var url string
	var set bool
	//
	url, set = os.LookupEnv(ENV_GRPC_SERVER_URL)
	if !set || url == "" {
		return ErrGrpcServerEnv
	}
	grpcServerUrl = url
	//
	url, set = os.LookupEnv(ENV_ELASTICSEARCH_URL)
	if !set || url == "" {
		return ErrElasticSearchEnv
	}
	elasticSearchUrl = url

	return nil
}

// Basically returns os.GetEnv(ENV_GRPC_SERVER_URL)
// NOTE: It is ASSUMED that function env.Init() was called BEFORE this function
// otherwise you WILL GET EMPTY STRING
func GRPCServerURL() string {
	return grpcServerUrl
}

// Basically returns os.GetEnv(ENV_ELASTICSEARCH_URL)
// NOTE: It is ASSUMED that function env.Init() was called BEFORE this function
// otherwise you WILL GET EMPTY STRING
func ElasticSearchURL() string {
	return elasticSearchUrl
}
