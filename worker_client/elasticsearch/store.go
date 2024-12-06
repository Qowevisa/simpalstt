package elasticsearch

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"
	pb "worker"

	"google.golang.org/protobuf/encoding/protojson"
)

// StoreDataInElasticsearch stores a Protobuf message in Elasticsearch
func StoreDataInElasticsearch(es8 ElasticClient, index string, data *pb.Data) error {
	// Convert Protobuf message to JSON
	jsonData, err := protojson.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling Protobuf to JSON: %w", err)
	}

	// Index the data in Elasticsearch
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := es8.Index(
		index,
		bytes.NewReader(jsonData),
		es8.IndexWithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("error indexing data in Elasticsearch: %w", err)
	}
	// this if construction prevents mock unit tests to segfault itself
	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	log.Printf("Document indexed in Elasticsearch: %s\n", data)
	return nil
}
