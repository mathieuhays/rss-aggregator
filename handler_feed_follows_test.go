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

func TestHandlePostFeedFollows(t *testing.T) {
	server, db, mock := createTestServer(t)
	defer db.Close()

	userApiKey := "testapikey"
	feedID := uuid.New()
	bodyReader := bytes.NewReader([]byte(`{"feed_id":"` + feedID.String() + `"}`))
	request, _ := http.NewRequest(http.MethodPost, "/v1/feed_follows", bodyReader)
	request.Header.Set("Authorization", "ApiKey "+userApiKey)
	response := httptest.NewRecorder()

	userId := uuid.New()
	userRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "api_key"}).
		AddRow(userId, time.Now().UTC(), time.Now().UTC(), "test", userApiKey)
	mock.ExpectQuery("SELECT id, created_at, updated_at, name, api_key FROM users").WillReturnRows(userRows)

	feedFollowRows := sqlmock.NewRows([]string{"id", "feed_id", "user_id", "created_at", "updated_at"}).
		AddRow(uuid.New(), feedID, userId, time.Now().UTC(), time.Now().UTC())
	mock.ExpectQuery("INSERT INTO feed_follows").WillReturnRows(feedFollowRows)

	server.ServeHTTP(response, request)

	assertStatus(t, response, http.StatusCreated)
	assertJSONContentType(t, response)
	assertBodyContains(t, response, `"feed_id":"`+feedID.String()+`"`)
	assertBodyContains(t, response, `"user_id":"`+userId.String()+`"`)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %s", err.Error())
	}
}

func TestHandleGetFeedFollows(t *testing.T) {
	server, db, mock := createTestServer(t)
	defer db.Close()

	userApiKey := "testapikey"
	request, _ := http.NewRequest(http.MethodGet, "/v1/feed_follows", nil)
	request.Header.Set("Authorization", "ApiKey "+userApiKey)
	response := httptest.NewRecorder()

	userId := uuid.New()
	userRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "api_key"}).
		AddRow(userId, time.Now().UTC(), time.Now().UTC(), "test", userApiKey)
	mock.ExpectQuery("SELECT id, created_at, updated_at, name, api_key FROM users").WillReturnRows(userRows)

	rows := sqlmock.NewRows([]string{"id", "feed_id", "user_id", "created_at", "updated_at"}).
		AddRow(uuid.New(), uuid.New(), uuid.New(), time.Now().UTC(), time.Now().UTC()).
		AddRow(uuid.New(), uuid.New(), uuid.New(), time.Now().UTC(), time.Now().UTC())
	mock.ExpectQuery("SELECT id, feed_id, user_id, created_at, updated_at FROM feed_follows").WillReturnRows(rows)

	server.ServeHTTP(response, request)

	assertStatus(t, response, http.StatusOK)
	assertJSONContentType(t, response)
	assertBodyContains(t, response, `"feed_id":`)
	assertBodyContains(t, response, `"user_id":`)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %s", err.Error())
	}
}

func TestHandleDeleteFeedFollow(t *testing.T) {
	server, db, mock := createTestServer(t)
	defer db.Close()

	userApiKey := "testapikey"
	feedFollowID := uuid.New()
	request, _ := http.NewRequest(http.MethodDelete, "/v1/feed_follows/"+feedFollowID.String(), nil)
	request.Header.Set("Authorization", "ApiKey "+userApiKey)
	response := httptest.NewRecorder()

	userId := uuid.New()
	userRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "api_key"}).
		AddRow(userId, time.Now().UTC(), time.Now().UTC(), "test", userApiKey)
	mock.ExpectQuery("SELECT id, created_at, updated_at, name, api_key FROM users").WillReturnRows(userRows)

	rows := sqlmock.NewRows([]string{"id", "feed_id", "user_id", "created_at", "updated_at"}).
		AddRow(feedFollowID, uuid.New(), userId, time.Now().UTC(), time.Now().UTC())
	mock.ExpectQuery("SELECT id, feed_id, user_id, created_at, updated_at FROM feed_follows").WillReturnRows(rows)

	mock.ExpectExec("DELETE FROM feed_follows").WillReturnResult(sqlmock.NewResult(1, 1))

	server.ServeHTTP(response, request)

	assertStatus(t, response, http.StatusOK)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %s", err.Error())
	}
}
