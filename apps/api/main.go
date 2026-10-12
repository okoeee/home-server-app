package main

import (
	"api/database"
	"api/db"
	"api/handler"
	"api/router"
	"log"
)

func main() {

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
