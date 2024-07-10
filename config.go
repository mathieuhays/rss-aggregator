package rss_aggregator

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mathieuhays/rss-aggregator/internal/auth"
	"github.com/mathieuhays/rss-aggregator/internal/database"
	"net/http"
)

type ApiConfig struct {
	DB *database.Queries
}

func NewApiConfig(db database.DBTX) (*ApiConfig, error) {
	return &ApiConfig{DB: database.New(db)}, nil
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAuthorization(r.Header)
		if err != nil || apiKey.Name != auth.TypeApiKey || apiKey.Value == "" {
			respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		ctx := context.Background()
		user, err := cfg.DB.GetUserByAPIKey(ctx, apiKey.Value)

		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusForbidden, "forbidden")
			return
		}

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		handler(w, r, user)
	}
}
