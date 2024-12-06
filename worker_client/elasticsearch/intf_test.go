package elasticsearch_test

import (
	"bytes"
	"context"
	"io"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/mock"
)

type MockElasticClient struct {
	mock.Mock
}

// Index mocks the Index method
func (m *MockElasticClient) Index(index string, body *bytes.Reader, opts ...func(*esapi.IndexRequest)) (*esapi.Response, error) {
	args := m.Called(index, body, opts)
	return args.Get(0).(*esapi.Response), args.Error(1)
}

// IndexWithContext mocks the IndexWithContext method
func (m *MockElasticClient) IndexWithContext(ctx context.Context) func(*esapi.IndexRequest) {
	return func(r *esapi.IndexRequest) {
		m.Called(ctx)
	}
}

// Search mocks the Search method
func (m *MockElasticClient) Search(opts ...func(*esapi.SearchRequest)) (*esapi.Response, error) {
	args := m.Called(opts)
	return args.Get(0).(*esapi.Response), args.Error(1)
}

// SearchWithContext mocks the SearchWithContext method
func (m *MockElasticClient) SearchWithContext(ctx context.Context) func(*esapi.SearchRequest) {
	return func(r *esapi.SearchRequest) {
		m.Called(ctx)
	}
}

// SearchWithIndex mocks the SearchWithIndex method
func (m *MockElasticClient) SearchWithIndex(indexes ...string) func(*esapi.SearchRequest) {
	return func(r *esapi.SearchRequest) {
		m.Called(indexes)
	}
}

// SearchWithBody mocks the SearchWithBody method
func (m *MockElasticClient) SearchWithBody(body io.Reader) func(*esapi.SearchRequest) {
	return func(r *esapi.SearchRequest) {
		m.Called(body)
	}
}
