package env_test

import (
	"os"
	"testing"
	"worker_client/env"
)

func TestInit(t *testing.T) {
	os.Setenv(env.ENV_ELASTICSEARCH_URL, "http://localhost:9200")
	os.Setenv(env.ENV_GRPC_SERVER_URL, "http://localhost:50051")
	defer os.Unsetenv(env.ENV_ELASTICSEARCH_URL)
	defer os.Unsetenv(env.ENV_GRPC_SERVER_URL)

	err := env.Init()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if env.ElasticSearchURL() != "http://localhost:9200" {
		t.Errorf("Unexpected Elasticsearch URL: %s", env.ElasticSearchURL())
	}

	if env.GRPCServerURL() != "http://localhost:50051" {
		t.Errorf("Unexpected gRPC Server URL: %s", env.GRPCServerURL())
	}
}

func TestInitMissingEnv(t *testing.T) {
	os.Unsetenv(env.ENV_ELASTICSEARCH_URL)
	os.Unsetenv(env.ENV_GRPC_SERVER_URL)

	err := env.Init()
	if err == nil {
		t.Fatal("Expected error for missing environment variables, got nil")
	}
}
