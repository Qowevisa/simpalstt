package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	pb "worker"
)

func encodeQuery(query map[string]interface{}) *bytes.Reader {
	data, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	return bytes.NewReader(data)
}

type ElasticsearchData struct {
	ID     string `json:"_id"`
	Source struct {
		Id         string `json:"Id"`
		Categories struct {
			Subcategory string `json:"subcategory"`
		} `json:"categories"`
		Title struct {
			Ro string `json:"ro"`
			Ru string `json:"ru"`
		} `json:"title"`
		Type   string  `json:"type"`
		Posted float64 `json:"posted"`
	} `json:"_source"`
	Sort []interface{} `json:"sort"`
}

type TermsBucket struct {
	Key      string `json:"key"`
	DocCount int64  `json:"doc_count"`
}

type SubcategoryCountsAggregation struct {
	DocCountErrorUpperBound int64         `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int64         `json:"sum_other_doc_count"`
	Buckets                 []TermsBucket `json:"buckets"`
}

type ElasticsearchResponse struct {
	Hits struct {
		Hits []ElasticsearchData `json:"hits"`
	} `json:"hits"`
	Aggregations struct {
		SubcategoryCounts SubcategoryCountsAggregation `json:"subcategory_counts"`
	} `json:"aggregations"`
}

func SearchDataBy(es8 ElasticClient, req *pb.DataSearch, ctx context.Context) (*ElasticsearchResponse, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Title,
				"fields": []string{"title.ru", "title.ro"},
			},
		},
		"size": req.Limit,
		"sort": []map[string]interface{}{
			{"posted": "asc"},
		},
	}

	if req.PageToken != "" {
		var searchAfter float64
		if err := json.Unmarshal([]byte(req.PageToken), &searchAfter); err != nil {
			return nil, fmt.Errorf("invalid PageToken: %w", err)
		}
		query["search_after"] = []interface{}{searchAfter}
	}

	res, err := es8.Search(
		es8.SearchWithContext(ctx),
		es8.SearchWithIndex(MainIndex),
		es8.SearchWithBody(encodeQuery(query)),
	)
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %w", err)
	}
	defer res.Body.Close()

	var esResult ElasticsearchResponse
	if err := json.NewDecoder(res.Body).Decode(&esResult); err != nil {
		return nil, fmt.Errorf("error decoding Elasticsearch response: %w", err)
	}
	return &esResult, nil
}

func SearchAggregatedDataBy(es8 ElasticClient, req *pb.AggregatedDataSearch, ctx context.Context) (*ElasticsearchResponse, error) {
	query := map[string]interface{}{
		"size": 0,
		"aggs": map[string]interface{}{"subcategory_counts": map[string]interface{}{
			"terms": map[string]string{
				"field": "categories.subcategory",
			},
		}},
	}
	if req.Categories != nil && req.Categories.Subcategory != "" {
		query["query"] = map[string]interface{}{
			"term": map[string]interface{}{
				"categories.subcategory": req.Categories.Subcategory,
			},
		}
	}

	res, err := es8.Search(
		es8.SearchWithContext(context.Background()),
		es8.SearchWithIndex(MainIndex),
		es8.SearchWithBody(encodeQuery(query)),
	)
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %w", err)
	}
	// this if construction prevents mock unit tests to segfault itself
	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}

	var esResult ElasticsearchResponse
	if err := json.NewDecoder(res.Body).Decode(&esResult); err != nil {
		return nil, fmt.Errorf("error decoding Elasticsearch response: %w", err)
	}
	return &esResult, nil
}
