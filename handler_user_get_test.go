package rss_aggregator

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleGetUser(t *testing.T) {
	t.Run("valid user", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		config, err := NewApiConfig(db)
		if err != nil {
			t.Fatal(err)
		}

		server, err := NewAggregatorServer(config)
		if err != nil {
			t.Fatal(err)
		}

		testApiKey := "testapikey"
		request, _ := http.NewRequest(http.MethodGet, "/v1/users", nil)
		request.Header.Set("Authorization", "ApiKey "+testApiKey)
		response := httptest.NewRecorder()

		id, _ := uuid.NewUUID()
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "api_key"}).
			AddRow(id, time.Now().UTC(), time.Now().UTC(), "test", testApiKey)
		mock.ExpectQuery("SELECT id, created_at, updated_at, name, api_key FROM users").WillReturnRows(rows)

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
		assertJSONContentType(t, response)
		assertBodyContains(t, response, `"name":"test"`)

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectations not met: %s", err.Error())
		}
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		config, err := NewApiConfig(db)
		if err != nil {
			t.Fatal(err)
		}

		server, err := NewAggregatorServer(config)
		if err != nil {
			t.Fatal(err)
		}

		testApiKey := "testapikey"
		request, _ := http.NewRequest(http.MethodGet, "/v1/users", nil)
		request.Header.Set("Authorization", "ApiKey "+testApiKey)
		response := httptest.NewRecorder()

		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "api_key"})
		mock.ExpectQuery("SELECT id, created_at, updated_at, name, api_key FROM users").WillReturnRows(rows)

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusForbidden)
		assertJSONContentType(t, response)

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectations not met: %s", err.Error())
		}
	})
}
