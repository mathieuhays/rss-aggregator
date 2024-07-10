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

func TestHandlePostUsers(t *testing.T) {
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

	bodyReader := bytes.NewReader([]byte(`{"name":"test"}`))
	request, _ := http.NewRequest(http.MethodPost, "/v1/users", bodyReader)
	response := httptest.NewRecorder()

	id, _ := uuid.NewUUID()
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "api_key"}).
		AddRow(id, time.Now().UTC(), time.Now().UTC(), "test", "test")
	mock.ExpectQuery("INSERT INTO users").WillReturnRows(rows)

	server.ServeHTTP(response, request)

	assertStatus(t, response, http.StatusCreated)
	assertJSONContentType(t, response)
	assertBodyContains(t, response, `"name":"test"`)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectations not met: %s", err.Error())
	}
}
