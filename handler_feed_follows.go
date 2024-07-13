package rss_aggregator

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mathieuhays/rss-aggregator/internal/database"
	"net/http"
	"time"
)

func (a *AggregatorServer) handlePostFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	var payload struct {
		FeedID string `json:"feed_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if payload.FeedID == "" {
		respondWithError(w, http.StatusBadRequest, "missing data")
		return
	}

	ctx := context.Background()
	feedFollow, err := a.config.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    uuid.MustParse(payload.FeedID),
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feedFollow))
}

func (a *AggregatorServer) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	ctx := context.Background()
	items, err := a.config.DB.GetFeedFollowsByUserId(ctx, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var publicItems []FeedFollow
	for _, item := range items {
		publicItems = append(publicItems, databaseFeedFollowToFeedFollow(item))
	}

	respondWithJSON(w, http.StatusOK, publicItems)
}

func (a *AggregatorServer) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID, err := uuid.Parse(r.PathValue("feedFollowID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid parameter")
		return
	}

	ctx := context.Background()
	feedFollow, err := a.config.DB.GetFeedFollow(ctx, feedFollowID)
	if err != nil || feedFollow.UserID != user.ID {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	err = a.config.DB.DeleteFeedFollow(ctx, feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
