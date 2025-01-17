package rss_aggregator

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mathieuhays/rss-aggregator/internal/database"
	"net/http"
	"time"
)

func (a *AggregatorServer) handlePostUsers(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name string
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Name == "" {
		respondWithError(w, http.StatusBadRequest, "name field is required")
		return
	}

	ctx := context.Background()
	user, err := a.config.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      payload.Name,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func (a *AggregatorServer) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
