package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
)

const (
	MainIndex = "articles"
)

func createIndex(es *elasticsearch8.Client) error {
	body := map[string]interface{}{
		"settings": map[string]interface{}{
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"my_russian_analyzer": map[string]interface{}{
						"type": "russian",
					},
					"my_romanian_analyzer": map[string]interface{}{
						"type": "romanian",
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"properties": map[string]interface{}{
						"ru": map[string]interface{}{
							"type":     "text",
							"analyzer": "my_russian_analyzer",
						},
						"ro": map[string]interface{}{
							"type":     "text",
							"analyzer": "my_romanian_analyzer",
						},
					},
				},
				"categories": map[string]interface{}{
					"properties": map[string]interface{}{
						"subcategory": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"type": map[string]interface{}{
					"type": "keyword",
				},
				"posted": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	res, err := es.Indices.Create(MainIndex,
		es.Indices.Create.WithBody(bytes.NewReader(data)),
		es.Indices.Create.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("Indices.Create: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Error creating index: %s", res.String())
	}

	log.Printf("Index '%s' created with Russian and Romanian analyzers!\n", MainIndex)
	return nil
}
