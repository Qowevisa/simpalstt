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
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
	"worker_client/elasticsearch"
	"worker_client/env"

	pb "worker"

	"worker_client/grpc"
)

func main() {
	if err := initializeServices(); err != nil {
		log.Fatalf("Initialization error: %v\n", err)
	}

	if err := fetchAndStoreData(); err != nil {
		log.Fatalf("Error during data fetching and storing: %v\n", err)
	}
	// Starting storage service
	// NOTE: it contains s.Serve that will block
	// any line of code written below StartStorageServer()
	grpc.StartStorageServer()
}

// initializeServices initializes all required dependencies.
func initializeServices() error {
	if err := env.Init(); err != nil {
		return fmt.Errorf("env.Init: %w", err)
	}
	if err := elasticsearch.Init(); err != nil {
		return fmt.Errorf("elasticsearch.Init: %w", err)
	}
	return nil
}

// fetchAndStoreData establishes a gRPC connection to the worker server,
// fetches the stream of data, and stores it in Elasticsearch.
func fetchAndStoreData() error {
	g, err := grpc.SafeGetConnectionToWorkerServer()
	if err != nil {
		return fmt.Errorf("grpc.SafeGetConnectionToWorkerServer: %w", err)
	}
	defer g.Conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	stream, err := g.Client.GetStreamOfData(ctx, &pb.DataFilter{})
	if err != nil {
		return fmt.Errorf("GetStreamOfData: %w", err)
	}

	for {
		data, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// This means the stream was successfully ended.
				break
			}
			return fmt.Errorf("stream.Recv: %w", err)
		}
		log.Printf("Received: %v\n", data)

		es8Conn := elasticsearch.GetEs8Connection()
		if err := elasticsearch.StoreDataInElasticsearch(es8Conn, elasticsearch.MainIndex, data); err != nil {
			log.Printf("Failed to store data in Elasticsearch: %v\n", err)
		}
	}
	return nil
}
