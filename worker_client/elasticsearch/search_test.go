package elasticsearch_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	pb "worker"
	"worker_client/elasticsearch"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/mock"
)

func TestSearchDataBy(t *testing.T) {
	mockClient := new(MockElasticClient)

	mockResponse := &esapi.Response{
		StatusCode: 201,
		Body: io.NopCloser(bytes.NewReader([]byte(`
			{
				"hits": {
					"hits": [
						{
							"_id": "1",
							"_source": {
								"Id": "1",
								"categories": {
									"subcategory": "test-category"
								},
								"title": {
									"ro": "No ti tle t eest",
									"ru": "test-title"
								},
								"type": "article",
								"posted": 1234567890
							},
							"sort": [1234567890]
						}
					]
				}
			}
		`))),
	}

	mockClient.On("Search", mock.Anything).Return(mockResponse, nil)

	req := &pb.DataSearch{
		Title: "test-title",
		Limit: 10,
	}
	ctx := context.Background()

	res, err := elasticsearch.SearchDataBy(mockClient, req, ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(res.Hits.Hits) != 1 {
		t.Fatalf("Expected 1 hit, got %d", len(res.Hits.Hits))
	}
	if res.Hits.Hits[0].Source.Title.Ru != "test-title" {
		t.Fatalf("Expected 'Заголовок RU', got '%s'", res.Hits.Hits[0].Source.Title.Ru)
	}

	mockClient.AssertExpectations(t)
}
