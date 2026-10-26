package main

import (
	"api/db"
	"api/handler"
	"api/infrastructure/mysql"
	"api/router"
	"api/usecase/dailyprice"
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
			log.Print(err)
		}
	}()

	queries := db.New(databaseConnection)

	stockHandler := handler.NewStockHandler(queries)
	dailyPriceRepository := mysql.NewDailyPriceRepository(queries)
	calculator := dailyprice.NewDailyPriceChangeCalculator(dailyPriceRepository)
	dailyPriceChangeHandler := handler.NewDailyPriceChangeHandler(calculator)
	volumeRatioCalculator := dailyprice.NewDailyVolumeRatioCalculator(dailyPriceRepository)
	dailyVolumeRatioHandler := handler.NewDailyVolumeRatioHandler(volumeRatioCalculator)

	routerErr := router.NewRouter(
		stockHandler,
		dailyPriceChangeHandler,
		dailyVolumeRatioHandler,
	).Run()

	if routerErr != nil {
		log.Fatal(routerErr)
	}

}
