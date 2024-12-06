//nolint:gofmt,golint
/* Data Flow:
-----------------        -----------------        ---------------
| worker_server | -----> | worker_client | -----> | graphql_api |
-----------------        -----------------        ---------------
																												^
																												|
																												|
																										We are here!
*/
package main

import (
	"errors"
	"fmt"
	"graphql_api/env"
	"graphql_api/graph"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	// This HAS to be called in the start
	// as it sets env.grpcServerUrl to os.GetEnv(env.ENV_GRPC_SERVER_URL)
	if err := env.Init(); err != nil {
		if errors.Is(err, env.ErrEnvNotSetOrEmpty) {
			fmt.Printf("ERROR: env.Init: %v\nDidn't you forget to set %s in Dockerfile?", err, env.ENV_GRPC_SERVER_URL)
			os.Exit(1)
		}
		fmt.Printf("ERROR: UNHANDLED ERROR: %v\n", err)
		os.Exit(2)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
