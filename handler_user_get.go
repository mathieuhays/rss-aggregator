package rss_aggregator

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mathieuhays/rss-aggregator/internal/auth"
	"net/http"
)

func (a *AggregatorServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAuthorization(r.Header)
	if err != nil || apiKey.Name != auth.TypeApiKey || apiKey.Value == "" {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	ctx := context.Background()
	user, err := a.config.DB.GetUser(ctx, apiKey.Value)

	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
