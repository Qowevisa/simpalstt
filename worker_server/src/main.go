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
	"log"
	"os"
	"worker_server/env"
	"worker_server/grpc"
)

func main() {
	if err := env.Init(); err != nil {
		log.Printf("ERROR: env.Init: %v\n", err)
		os.Exit(1)
	}
	grpc.StartWorkerServer()
}
