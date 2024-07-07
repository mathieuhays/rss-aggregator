package rss_aggregator

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAggregatorHandlers(t *testing.T) {
	server, err := NewAggregatorServer(&ApiConfig{})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("readiness handler", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/v1/healthz", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertJSONContentType(t, response)
		assertBodyContains(t, response, "\"ok\"")
	})

	t.Run("error handler", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/v1/err", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusInternalServerError)
		assertJSONContentType(t, response)
		assertBodyContains(t, response, "\"Internal")
	})
}
