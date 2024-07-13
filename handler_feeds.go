package rss_aggregator

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mathieuhays/rss-aggregator/internal/database"
	"net/http"
	"time"
)

func (a *AggregatorServer) handlePostFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	var payload struct {
		Name string
		URL  string
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Name == "" || payload.URL == "" {
		respondWithError(w, http.StatusBadRequest, "missing data")
		return
	}

	ctx := context.Background()
	feed, err := a.config.DB.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      payload.Name,
		Url:       payload.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	feedFollow, err := a.config.DB.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}{
		Feed:       databaseFeedToFeed(feed),
		FeedFollow: databaseFeedFollowToFeedFollow(feedFollow),
	})
}

func (a *AggregatorServer) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	items, err := a.config.DB.GetFeeds(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var publicItems []Feed
	for _, item := range items {
		publicItems = append(publicItems, databaseFeedToFeed(item))
	}

	respondWithJSON(w, http.StatusOK, publicItems)
}
