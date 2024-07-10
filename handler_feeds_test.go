package rss_aggregator

import (
	"bytes"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandlePostFeeds(t *testing.T) {
	server, db, mock := createTestServer(t)
	defer db.Close()

	userApiKey := "testapikey"
	bodyReader := bytes.NewReader([]byte(`{"name":"test", "url":"https://example.com"}`))
	request, _ := http.NewRequest(http.MethodPost, "/v1/feeds", bodyReader)
	request.Header.Set("Authorization", "ApiKey "+userApiKey)
	response := httptest.NewRecorder()

	userId := uuid.New()
	userRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "api_key"}).
		AddRow(userId, time.Now().UTC(), time.Now().UTC(), "test", userApiKey)
	mock.ExpectQuery("SELECT id, created_at, updated_at, name, api_key FROM users").WillReturnRows(userRows)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "url", "user_id"}).
		AddRow(uuid.New(), time.Now().UTC(), time.Now().UTC(), "test", "https://example.com", userId)
	mock.ExpectQuery("INSERT INTO feeds").WillReturnRows(rows)

	server.ServeHTTP(response, request)

	assertStatus(t, response, http.StatusCreated)
	assertJSONContentType(t, response)
	assertBodyContains(t, response, `"name":"test"`)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %s", err.Error())
	}
}
