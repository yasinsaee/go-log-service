package app

import (
	"github.com/yasinsaee/go-log-service/internal/app/config"
	"github.com/yasinsaee/go-log-service/pkg/elastic"
)

func initElastic() {
	elastic.Init(elastic.Config{
		Addresses: []string{config.GetEnv("ELASTIC_ADDRESS", "http://localhost:9200")},
		Username:  config.GetEnv("ELASTIC_USERNAME", ""),
		Password:  config.GetEnv("ELASTIC_PASSWORD", ""),
	})
}
