package main

import (
	"api/db"
	"api/domain/model"
	"api/external/jquants/fetch"
	"api/external/jquants/mapper"
	jquantsmodel "api/external/jquants/model"
	"api/infrastructure/mysql"
	"api/usecase/dailyprice"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	if len(os.Args) != 3 {
		log.Fatal("usage: go run ./cmd/backfilldailyprices YYYY-MM-DD YYYY-MM-DD")
	}

	from, err := time.Parse("2006-01-02", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	to, err := time.Parse("2006-01-02", os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	if from.After(to) {
		log.Fatal("invalid date")
	}

	apiKey := os.Getenv("JQUANTS_API_KEY")

	ctx := context.Background()

	client, err := fetch.NewClient(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	listedTradingCalendar, err := fetchTradingDate(ctx, client, from, to)
	if err != nil {
		log.Fatal(err)
	}

	tradingDates, err := mapper.ToTradingDates(listedTradingCalendar)
	if err != nil {
		log.Fatal(err)
	}

	databaseConnection, err := mysql.Open()
	if err != nil {
		log.Fatalf("open database connection error: %w", err)
	}
	defer func() {
		err := databaseConnection.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	for _, tradingDate := range tradingDates.Date {
		jQuantsDailyPrice, err := fetchDailyPrices(ctx, client, tradingDate)
		if err != nil {
			log.Fatal(err)
		}

		dailyPrice, err := mapper.ToDailyPrices(jQuantsDailyPrice)
		if err != nil {
			log.Fatal(err)
		}

		err = save(ctx, databaseConnection, dailyPrice)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(
			"daily prices imported: date=%s, count=%d",
			tradingDate.Format(time.DateOnly),
			len(dailyPrice),
		)
	}

}

func fetchTradingDate(
	ctx context.Context,
	client *fetch.JQuantsClient,
	from, to time.Time,
) (jquantsmodel.ListedTradingCalendar, error) {

	result, err := client.FetchTradingDays(ctx, from, to)
	if err != nil {
		return jquantsmodel.ListedTradingCalendar{}, fmt.Errorf("error fetching trading days: %w", err)
	}

	return result, nil

}

func fetchDailyPrices(
	ctx context.Context,
	client *fetch.JQuantsClient,
	date time.Time,
) ([]jquantsmodel.DailyPrice, error) {

	jQuantsDailyPrices, err := client.FetchDailyPrices(ctx, date)
	if err != nil {
		return []jquantsmodel.DailyPrice{}, fmt.Errorf("fetch daily prices error: %w", err)
	}

	return jQuantsDailyPrices, nil

}

func save(
	ctx context.Context,
	databaseConnection *sql.DB,
	dailyPrices []model.DailyPrice,
) error {

	transaction, err := databaseConnection.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin database transaction error: %w", err)
	}
	defer func() {
		_ = transaction.Rollback()
	}()

	queries := db.New(transaction)
	dailyPriceRepository := mysql.NewDailyPriceRepository(queries)

	importer := dailyprice.NewDailyPriceImporter(dailyPriceRepository)

	err = importer.Import(ctx, dailyPrices)
	if err != nil {
		return fmt.Errorf("import daily prices error: %w", err)
	}

	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("commit database transaction error: %w", err)
	}

	return nil

}
