package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	Addresses []string
	Username  string
	Password  string
}

type Client struct {
	ES *elasticsearch.Client
}

var ClientInstance *Client

func Init(cfg Config) (*Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elastic client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := es.Info(es.Info.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elastic: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elastic connection error: %s", res.String())
	}

	log.Println("✅ Connected to Elasticsearch successfully")

	ClientInstance = &Client{ES: es}
	return ClientInstance, nil
}

type LogEntry struct {
	Level string `json:"level"`

	Message string `json:"message"`

	Service string `json:"service"`

	Module string `json:"module,omitempty"`

	RequestID string `json:"request_id,omitempty"`

	UserID string `json:"user_id,omitempty"`

	Host string `json:"host,omitempty"`

	Extra map[string]interface{} `json:"extra,omitempty"`

	Error string `json:"error,omitempty"`

	Timestamp time.Time `json:"timestamp"`
}

func (c *Client) IndexDocument(index string, doc interface{}) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:   index,
		Body:    bytes.NewReader(body),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), c.ES)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index error: %s", res.String())
	}

	return nil
}

type SearchQuery struct {
	Field string
	Value string
	Size  int
}

func (c *Client) SearchDocuments(index string, query SearchQuery) ([]map[string]interface{}, error) {
	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				query.Field: query.Value,
			},
		},
		"size": query.Size,
	}

	body, err := json.Marshal(searchBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search body: %w", err)
	}

	res, err := c.ES.Search(
		c.ES.Search.WithContext(context.Background()),
		c.ES.Search.WithIndex(index),
		c.ES.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elastic search error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	hits, ok := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, nil
	}

	var docs []map[string]interface{}
	for _, h := range hits {
		source := h.(map[string]interface{})["_source"].(map[string]interface{})
		docs = append(docs, source)
	}

	return docs, nil
}
