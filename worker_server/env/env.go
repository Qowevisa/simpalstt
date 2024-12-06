package env

import (
	"errors"
	"fmt"
	"os"
)

const (
	ENV_JSON_FILE = "JSON_FILE"
)

var (
	ErrEnvNotSetOrEmpty = errors.New("env not set or is empty")
	ErrJSONFile         = fmt.Errorf("%s: %w", ENV_JSON_FILE, ErrEnvNotSetOrEmpty)
)

var (
	jsonFilepath string
)

// Init changes internal private variable jsonFilepath
// To get the value use env.JSONFilepath function
func Init() error {
	filepath, set := os.LookupEnv(ENV_JSON_FILE)
	if !set || filepath == "" {
		return ErrJSONFile
	}
	jsonFilepath = filepath
	return nil
}

// Basically returns os.GetEnv(ENV_GRPC_SERVER_URL)
// NOTE: It is ASSUMED that function env.Init() was called BEFORE this function
// otherwise you WILL GET EMPTY STRING
func JSONFilepath() string {
	return jsonFilepath
}
