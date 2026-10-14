package main

import (
	"api/database"
	"api/db"
	"api/external/jquants/fetch"
	"api/handler"
	"api/router"
	"context"
	"log"
	"os"
)

func main() {

	// test code
	jquantsAPIKey := os.Getenv("JQUANTS_API_KEY")

	client, err := fetch.NewClient(jquantsAPIKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	equityMasterList, err := client.FetchEquityMasterList(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("equityMasterList:", equityMasterList)
	// test code

	databaseConnection, databaseErr := database.Open()
	if databaseErr != nil {
		log.Fatal(databaseErr)
	}

	defer databaseConnection.Close()

	queries := db.New(databaseConnection)

	stockHandler := handler.NewStockHandler(queries)

	routerErr := router.NewRouter(stockHandler).Run()

	if routerErr != nil {
		log.Fatal(routerErr)
	}

}
