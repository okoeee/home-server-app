package main

import (
	"api/external/jquants/fetch"
	"context"
	"log"
	"os"
	"time"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("usage: go run ./cmd/fetchdailyprices YYYY-MM-DD")
	}

	date, err := time.Parse(time.DateOnly, os.Args[1])
	if err != nil {
		log.Fatalf("date format error: %v", err)
	}

	apiKey := os.Getenv("JQUANTS_API_KEY")

	client, err := fetch.NewClient(apiKey)
	if err != nil {
		log.Fatalf("fetch client error: %v", err)
	}

	ctx := context.Background()

	dailyPrices, err := client.FetchDailyPrices(
		ctx,
		date,
	)
	if err != nil {
		log.Fatalf("fetch daily prices error: %v", err)
	}

	log.Printf("daily prices: date=%s count=%d", date.Format(time.DateOnly), len(dailyPrices))

	if len(dailyPrices) > 0 {
		log.Printf("sample: %v", dailyPrices[0])
	}

}
