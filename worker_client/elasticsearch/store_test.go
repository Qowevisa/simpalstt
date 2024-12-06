package elasticsearch_test

import (
	"testing"
	pb "worker"
	"worker_client/elasticsearch"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/mock"
)

func TestStoreDataInElasticsearch(t *testing.T) {
	// Mock Elasticsearch client
	mockClient := new(MockElasticClient)

	mockClient.On("Index", elasticsearch.MainIndex, mock.Anything, mock.Anything).Return(&esapi.Response{StatusCode: 200}, nil)

	data := &pb.Data{XId: "test"}
	err := elasticsearch.StoreDataInElasticsearch(mockClient, elasticsearch.MainIndex, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	mockClient.AssertExpectations(t)
}
