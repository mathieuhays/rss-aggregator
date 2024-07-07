package rss_aggregator

import (
	"github.com/mathieuhays/rss-aggregator/internal/database"
)

type ApiConfig struct {
	DB *database.Queries
}

func NewApiConfig(db database.DBTX) (*ApiConfig, error) {
	return &ApiConfig{DB: database.New(db)}, nil
}
