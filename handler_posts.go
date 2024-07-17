package rss_aggregator

import (
	"context"
	"github.com/mathieuhays/rss-aggregator/internal/database"
	"net/http"
)

func (a *AggregatorServer) handleGetPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	ctx := context.Background()
	posts, err := a.config.DB.GetPostsByUser(ctx, database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var public []Post
	for _, item := range posts {
		public = append(public, databasePostToPost(item))
	}

	respondWithJSON(w, http.StatusOK, public)
}
