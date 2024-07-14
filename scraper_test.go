package rss_aggregator

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/mathieuhays/rss-aggregator/internal/database"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchFeed(t *testing.T) {
	t.Run("valid fetch", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			source := `
<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>Boot.dev Blog</title>
    <link>https://blog.boot.dev/</link>
    <description>Recent content on Boot.dev Blog</description>
    <generator>Hugo</generator>
    <language>en-us</language>
    <lastBuildDate>Wed, 10 Jul 2024 00:00:00 +0000</lastBuildDate>
    <atom:link href="https://blog.boot.dev/index.xml" rel="self" type="application/rss+xml" />
    <item>
      <title>The Boot.dev Beat. July 2024</title>
      <link>https://blog.boot.dev/news/bootdev-beat-2024-07/</link>
      <pubDate>Wed, 10 Jul 2024 00:00:00 +0000</pubDate>
      <guid>https://blog.boot.dev/news/bootdev-beat-2024-07/</guid>
      <description>One million lessons. Well, to be precise, you have all completed 1,122,050 lessons just in June.</description>
    </item>
    <item>
      <title>The Boot.dev Beat. June 2024</title>
      <link>https://blog.boot.dev/news/bootdev-beat-2024-06/</link>
      <pubDate>Wed, 05 Jun 2024 00:00:00 +0000</pubDate>
      <guid>https://blog.boot.dev/news/bootdev-beat-2024-06/</guid>
      <description>ThePrimeagen&amp;rsquo;s new Git course is live. A new boss battle is on the horizon, and we&amp;rsquo;ve made massive speed improvements to the site.</description>
    </item>
</channel>
</rss>`
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(source))
		}))
		defer ts.Close()

		testFeed := database.Feed{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Name:          "test",
			Url:           ts.URL,
			UserID:        uuid.New(),
			LastFetchedAt: sql.NullTime{},
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		dbQueries := database.New(db)

		scraper := &Scraper{
			httpClient: http.Client{},
			db:         dbQueries,
		}

		feedRows := mock.NewRows([]string{"id", "created_at", "updated_at", "name", "url", "user_id", "last_fetched_at"}).
			AddRow(testFeed.ID, testFeed.CreatedAt, time.Now(), testFeed.Name, ts.URL, testFeed.UserID, time.Now())
		mock.ExpectQuery("UPDATE feeds").WillReturnRows(feedRows)

		rssFeed, err := scraper.fetchFeed(testFeed)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		if len(rssFeed.Channel.Items) != 2 {
			t.Errorf("expected 2 items, got %d instead", len(rssFeed.Channel.Items))
		}

		if rssFeed.Channel.Items[0].Title != "The Boot.dev Beat. July 2024" {
			t.Errorf("wrong item title. Expected 'The Boot.dev Beat. July 2024'. Got %s", rssFeed.Channel.Items[0].Title)
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectations not met: %s", err.Error())
		}
	})

	t.Run("fetch url not found", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))

		testFeed := database.Feed{
			ID:            uuid.New(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Name:          "test",
			Url:           ts.URL,
			UserID:        uuid.New(),
			LastFetchedAt: sql.NullTime{},
		}

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		dbQueries := database.New(db)

		scraper := &Scraper{
			httpClient: http.Client{},
			db:         dbQueries,
		}

		feedRows := mock.NewRows([]string{"id", "created_at", "updated_at", "name", "url", "user_id", "last_fetched_at"}).
			AddRow(testFeed.ID, testFeed.CreatedAt, time.Now(), testFeed.Name, ts.URL, testFeed.UserID, time.Now())
		mock.ExpectQuery("UPDATE feeds").WillReturnRows(feedRows)

		_, err = scraper.fetchFeed(testFeed)
		if err == nil {
			t.Fatal("fetch should have failed but returned with no error")
		}

		if !strings.HasPrefix(err.Error(), "feed unavailable") {
			t.Errorf("wrong error detected. Expected: feed unavailable. got %s", err.Error())
		}

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectations not met: %s", err.Error())
		}
	})
}
