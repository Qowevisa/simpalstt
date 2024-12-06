package env_test

import (
	"os"
	"testing"
	"worker_server/env"
)

func TestInit_ValidEnv(t *testing.T) {
	os.Setenv(env.ENV_JSON_FILE, "/path/to/file.json")
	defer os.Unsetenv(env.ENV_JSON_FILE)

	if err := env.Init(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if env.JSONFilepath() != "/path/to/file.json" {
		t.Errorf("Expected filepath to be '/path/to/file.json', got %v", env.JSONFilepath())
	}
}

func TestInit_MissingEnv(t *testing.T) {
	os.Unsetenv(env.ENV_JSON_FILE)

	if err := env.Init(); err == nil {
		t.Fatal("Expected error, got nil")
	}
}
