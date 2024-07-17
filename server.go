package rss_aggregator

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const ResponseTimeFormat = time.RFC3339

type AggregatorServer struct {
	// probably a store at some point
	http.Handler
	config *ApiConfig
}

func NewAggregatorServer(config *ApiConfig) (*AggregatorServer, error) {
	s := new(AggregatorServer)
	s.config = config

	router := http.NewServeMux()
	router.Handle("/v1/healthz", http.HandlerFunc(handlerReadiness))
	router.Handle("/v1/err", http.HandlerFunc(handlerErr))

	router.Handle("GET /v1/users", config.middlewareAuth(s.handleGetUser))
	router.Handle("POST /v1/users", http.HandlerFunc(s.handlePostUsers))

	router.Handle("GET /v1/feeds", http.HandlerFunc(s.handleGetFeeds))
	router.Handle("POST /v1/feeds", config.middlewareAuth(s.handlePostFeeds))

	router.Handle("GET /v1/feed_follows", config.middlewareAuth(s.handleGetFeedFollows))
	router.Handle("POST /v1/feed_follows", config.middlewareAuth(s.handlePostFeedFollows))
	router.Handle("DELETE /v1/feed_follows/{feedFollowID}", config.middlewareAuth(s.handleDeleteFeedFollow))

	router.Handle("GET /v1/posts", config.middlewareAuth(s.handleGetPosts))

	s.Handler = router

	return s, nil
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{Status: "ok"})
}

func handlerErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Printf("Responding with server error (%d): %s", code, message)
	}

	payload := struct {
		Error string `json:"error"`
	}{Error: message}

	respondWithJSON(w, code, payload)
}
