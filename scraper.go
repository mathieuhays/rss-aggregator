package rss_aggregator

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/mathieuhays/rss-aggregator/internal/database"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Scraper struct {
	httpClient http.Client
	db         *database.Queries
}

type RSSItem struct {
	Title          string `xml:"title"`
	Link           string `xml:"link"`
	PublishingDate string `xml:"pubDate"`
	Guid           string `xml:"guid"`
	Description    string `xml:"description"`
}

type RSSFeed struct {
	Channel struct {
		Title         string    `xml:"title"`
		Link          string    `xml:"link"`
		Description   string    `xml:"description"`
		Generator     string    `xml:"generator"`
		Language      string    `xml:"language"`
		LastBuildDate string    `xml:"lastBuildDate"`
		Items         []RSSItem `xml:"item"`
	} `xml:"channel"`
}

func NewScraper(db *database.Queries, interval time.Duration) (*Scraper, error) {
	scraper := &Scraper{
		httpClient: http.Client{Timeout: time.Second * 5},
		db:         db,
	}

	go scraper.realLoop(interval)

	return scraper, nil
}

func (s *Scraper) realLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		s.refreshFeeds()
	}
}

func (s *Scraper) fetchFeed(feed database.Feed) (*RSSFeed, error) {
	request, err := http.NewRequest(http.MethodGet, feed.Url, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("request error: %s", err.Error()))
	}

	response, err := s.httpClient.Do(request)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("request execution error: %s", err.Error()))
	}
	defer response.Body.Close()

	ctx := context.Background()
	_, err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		ID:            feed.ID,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not mark as fetched: %s", err.Error()))
	}

	if response.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("feed unavailable. responded with status %d", response.StatusCode))
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing body: %s", err.Error()))
	}

	var result RSSFeed

	err = xml.Unmarshal(data, &result)
	if err != nil {
		return &RSSFeed{}, err
	}

	return &result, nil
}

func (s *Scraper) refreshFeeds() {
	log.Println("refreshFeeds triggered")
	defer log.Println("refreshFeeds done")

	ctx := context.Background()
	feeds, err := s.db.GetNextFeedsToFetch(ctx, 10)
	if err != nil {
		log.Printf("error retrieving feeds: %s", err.Error())
		return
	}

	var wg sync.WaitGroup

	for _, dbFeed := range feeds {
		wg.Add(1)
		go func(feed database.Feed) {
			defer wg.Done()
			rssFeed, err := s.fetchFeed(feed)
			if err != nil {
				log.Printf("fetch error for %s: %s", feed.Url, err.Error())
				return
			}

			log.Printf("%s fetched successfully. %d items found", feed.Url, len(rssFeed.Channel.Items))
		}(dbFeed)
	}

	wg.Wait()
}
