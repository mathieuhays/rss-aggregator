package rss_aggregator

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleGetPosts(t *testing.T) {
	server, db, mock := createTestServer(t)
	defer db.Close()

	userApiKey := "testapikey"
	request, _ := http.NewRequest(http.MethodGet, "/v1/posts", nil)
	request.Header.Set("Authorization", "ApiKey "+userApiKey)
	response := httptest.NewRecorder()

	userId := uuid.New()
	userRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "api_key"}).
		AddRow(userId, time.Now().UTC(), time.Now().UTC(), "test", userApiKey)
	mock.ExpectQuery("SELECT id, created_at, updated_at, name, api_key FROM users").WillReturnRows(userRows)

	foundPostRows := sqlmock.NewRows([]string{"id", "created_at", "published_at", "title", "url", "description", "published_at", "feed_Id"}).
		AddRow(uuid.New(), time.Now(), time.Now(), "test", "https://test.com", sql.NullString{
			String: "",
			Valid:  true,
		}, sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}, uuid.New()).
		AddRow(uuid.New(), time.Now(), time.Now(), "test 2", "https://test.com/2", sql.NullString{
			String: "test",
			Valid:  true,
		}, sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}, uuid.New())
	mock.ExpectQuery("SELECT p.id, p.created_at, p.updated_at, p.title, p.url, p.description, p.published_at, p.feed_id FROM posts p").WillReturnRows(foundPostRows)

	server.ServeHTTP(response, request)

	assertStatus(t, response, http.StatusOK)
	assertJSONContentType(t, response)
	assertBodyContains(t, response, `"title":"test"`)
	assertBodyContains(t, response, `"url":"https://test.com"`)
	assertBodyContains(t, response, `"title":"test 2"`)
	assertBodyContains(t, response, `"url":"https://test.com/2`)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %s", err.Error())
	}
}
