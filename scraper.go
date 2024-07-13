package rss_aggregator

import (
	"context"
	"database/sql"
	"encoding/xml"
	"github.com/mathieuhays/rss-aggregator/internal/database"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Scraper struct {
	httpClient http.Client
	config     *ApiConfig
}

func NewScraper(config *ApiConfig, interval time.Duration) (*Scraper, error) {
	scraper := &Scraper{
		httpClient: http.Client{Timeout: time.Second * 5},
		config:     config,
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

func (s *Scraper) refreshFeeds() {
	log.Println("refreshFeeds triggered")
	defer log.Println("refreshFeeds done")

	ctx := context.Background()
	feeds, err := s.config.DB.GetNextFeedsToFetch(ctx, 10)
	if err != nil {
		log.Printf("error retrieving feeds: %s", err.Error())
		return
	}

	var wg sync.WaitGroup

	for _, dbFeed := range feeds {
		wg.Add(1)
		go func(feed database.Feed) {
			defer wg.Done()
			request, err := http.NewRequest(http.MethodGet, feed.Url, nil)
			if err != nil {
				log.Printf("request error for %s. error: %s", feed.Url, err.Error())
				return
			}

			response, err := s.httpClient.Do(request)
			if err != nil {
				log.Printf("error fetching %s. error: %s", feed.Url, err.Error())
				return
			}
			defer response.Body.Close()

			err = s.config.DB.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
				LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
				ID:            feed.ID,
			})
			if err != nil {
				log.Printf("could not mark as fetched %s. error: %s", feed.Url, err.Error())
				return
			}

			log.Printf("%s marked as fetched", feed.Url)

			if response.StatusCode > 299 {
				log.Printf("request for %s came back with code %d. aborting...", feed.Url, response.StatusCode)
				return
			}

			data, err := io.ReadAll(response.Body)
			if err != nil {
				log.Printf("error reading request body for %s. error: %s", feed.Url, err.Error())
				return
			}

			items, err := extractRSSItems(string(data))
			if err != nil {
				log.Printf("error extracting items for %s. error: %s", feed.Url, err.Error())
				return
			}

			log.Printf("%d items found for %s", len(items), feed.Url)
		}(dbFeed)
	}

	wg.Wait()
}

type RSSItem struct {
	Title          string `xml:"title"`
	Link           string `xml:"link"`
	PublishingDate string `xml:"pubDate"`
	Guid           string `xml:"guid"`
	Description    string `xml:"description"`
}

type RSSChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Generator     string    `xml:"generator"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []RSSItem `xml:"item"`
}

func extractRSSItems(input string) ([]RSSItem, error) {
	var result struct {
		Channel RSSChannel `xml:"channel"`
	}

	err := xml.Unmarshal([]byte(input), &result)
	if err != nil {
		return []RSSItem{}, err
	}

	return result.Channel.Items, nil
}
