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

	feedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "url", "user_id"}).
		AddRow(uuid.New(), time.Now().UTC(), time.Now().UTC(), "test", "https://example.com", userId)
	mock.ExpectQuery("INSERT INTO feeds").WillReturnRows(feedRows)

	feedFollowRows := sqlmock.NewRows([]string{"id", "feed_id", "user_id", "created_at", "updated_at"}).
		AddRow(uuid.New(), uuid.New(), uuid.New(), time.Now().UTC(), time.Now().UTC())
	mock.ExpectQuery("INSERT INTO feed_follows").WillReturnRows(feedFollowRows)

	server.ServeHTTP(response, request)

	assertStatus(t, response, http.StatusCreated)
	assertJSONContentType(t, response)
	assertBodyContains(t, response, `"name":"test"`)
	assertBodyContains(t, response, `"feed":{`)
	assertBodyContains(t, response, `"feed_follow":{`)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %s", err.Error())
	}
}

func TestHandleGetFeeds(t *testing.T) {
	server, db, mock := createTestServer(t)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "url", "user_id"}).
		AddRow(uuid.New(), time.Now().UTC(), time.Now().UTC(), "test", "https://example.com", uuid.New()).
		AddRow(uuid.New(), time.Now().UTC(), time.Now().UTC(), "test 2", "https://example.com/2.xml", uuid.New())
	mock.ExpectQuery("SELECT id, created_at, updated_at, name, url, user_id FROM feeds").WillReturnRows(rows)

	request, _ := http.NewRequest(http.MethodGet, "/v1/feeds", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assertStatus(t, response, http.StatusOK)
	assertJSONContentType(t, response)
	assertBodyContains(t, response, `"name":"test"`)
	assertBodyContains(t, response, `"name":"test 2"`)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %s", err.Error())
	}
}
