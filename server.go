package rss_aggregator

import (
	"encoding/json"
	"log"
	"net/http"
)

type AggregatorServer struct {
	// probably a store at some point
	http.Handler
}

func NewAggregatorServer() (*AggregatorServer, error) {
	s := new(AggregatorServer)

	router := http.NewServeMux()
	router.Handle("/v1/healthz", http.HandlerFunc(handlerReadiness))
	router.Handle("/v1/err", http.HandlerFunc(handlerErr))

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
