package elasticsearch

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"worker_client/env"

	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// RealElasticClient wraps the real Elasticsearch client
type RealElasticClient struct {
	Client *elasticsearch8.Client
}

// Index implements the ElasticClient:1
func (rec RealElasticClient) Index(index string, body *bytes.Reader, opts ...func(*esapi.IndexRequest)) (*esapi.Response, error) {
	return rec.Client.Index(index, body, opts...)
}

// IndexWithContext implements the ElasticClient:2
func (rec RealElasticClient) IndexWithContext(ctx context.Context) func(*esapi.IndexRequest) {
	return rec.Client.Index.WithContext(ctx)
}

// Search implements the ElasticClient:3
func (rec RealElasticClient) Search(opts ...func(*esapi.SearchRequest)) (*esapi.Response, error) {
	return rec.Client.Search(opts...)
}

// SearchWithContext implements the ElasticClient:4
func (rec RealElasticClient) SearchWithContext(v context.Context) func(*esapi.SearchRequest) {
	return rec.Client.Search.WithContext(v)
}

// SearchWithIndex implements the ElasticClient:5
func (rec RealElasticClient) SearchWithIndex(v ...string) func(*esapi.SearchRequest) {
	return rec.Client.Search.WithIndex(v...)
}

// SearchWithBody implements the ElasticClient:6
func (rec RealElasticClient) SearchWithBody(v io.Reader) func(*esapi.SearchRequest) {
	return rec.Client.Search.WithBody(v)
}

var es8 *elasticsearch8.Client

func Init() error {
	es8URL := env.ElasticSearchURL()
	es8Connection, err := elasticsearch8.NewClient(elasticsearch8.Config{
		Addresses: []string{es8URL},
	})
	if err != nil {
		panic(err)
	}
	resp, err := es8Connection.Indices.Exists([]string{MainIndex})
	if err != nil {
		return fmt.Errorf("es8Connection.Indices.Exists: %w", err)
	}
	if resp.StatusCode != 200 {
		if err := createIndex(es8Connection); err != nil {
			return fmt.Errorf("createIndex: %w", err)
		}
	}
	es8 = es8Connection
	// Checking health of elasticsearch
	res, err := es8Connection.Cluster.Health()
	if err != nil {
		return fmt.Errorf("es8Connection.Cluster.Health: %w", err)
	}
	defer res.Body.Close()
	log.Printf("Elasticsearch cluster health: %v\n", res)
	return nil
}

// Basically returns package variable elasticsearch.es8
// NOTE: It is ASSUMED that function elasticsearch.Init() was called BEFORE this function
// otherwise you WILL GET `nil`!
func GetEs8Connection() *RealElasticClient {
	return &RealElasticClient{
		Client: es8,
	}
}
