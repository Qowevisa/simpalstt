package elasticsearch

import (
	"bytes"
	"context"
	"io"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// ElasticClient defines the minimal interface required for Elasticsearch operations
type ElasticClient interface {
	Index(index string, body *bytes.Reader, opts ...func(*esapi.IndexRequest)) (*esapi.Response, error)
	// Replaces es8.Index.WithContext
	IndexWithContext(v context.Context) func(*esapi.IndexRequest)

	// Immitates es8.Search function for unit-testing
	Search(o ...func(*esapi.SearchRequest)) (*esapi.Response, error)
	// Replaces es8.Search.WithContext
	SearchWithContext(v context.Context) func(*esapi.SearchRequest)
	// Replaces es8.Search.WithContext
	SearchWithIndex(v ...string) func(*esapi.SearchRequest)
	// Replaces es8.Search.WithContext
	SearchWithBody(v io.Reader) func(*esapi.SearchRequest)
}

// We need ElasticClient.IndexWithContext and other function that clearly
// just immitates real es8 API (e.g. elasticsearch.Client.Index.WithContext function)
// because when we want to unit-test it we just can't write that interface ElasticClient
// can suddenly implements Index function AND 'have' Index field that implements
// WithContext function
// Interfaces can't have fields
