package main

import (
	"api/db"
	"api/handler"
	"api/infrastructure/mysql"
	"api/router"
	"log"
)

func main() {

	databaseConnection, databaseErr := mysql.Open()
	if databaseErr != nil {
		log.Fatal(databaseErr)
	}

	defer func() {
		err := databaseConnection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	queries := db.New(databaseConnection)

	stockHandler := handler.NewStockHandler(queries)

	routerErr := router.NewRouter(stockHandler).Run()

	if routerErr != nil {
		log.Fatal(routerErr)
	}

}
