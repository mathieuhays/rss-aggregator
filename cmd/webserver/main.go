package main

import (
	"github.com/joho/godotenv"
	rss "github.com/mathieuhays/rss-aggregator"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	server, err := rss.NewAggregatorServer()
	if err != nil {
		log.Fatal(err)
	}

	addr := "localhost:" + port

	log.Printf("Starting local server on port %s", port)
	log.Fatal(http.ListenAndServe(addr, server))
}
