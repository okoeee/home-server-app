package main

import (
	"api/db"
	"api/external/openmeteo"
	openmeteofetch "api/external/openmeteo/fetch"
	"api/external/prometheus"
	"api/handler"
	"api/infrastructure/mysql"
	"api/router"
	"api/usecase/dailyprice"
	"api/usecase/weather"
	"log"
	"os"
	"time"
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

	weatherClient := openmeteofetch.NewClient()
	forecastReader := openmeteo.NewForecastReader(weatherClient, time.Now)
	forecastCache := weather.NewForecastCache(forecastReader, time.Now)
	weatherHandler := handler.NewWeatherHandler(forecastCache)
	prometheusURL := os.Getenv("PROMETHEUS_URL")
	if prometheusURL == "" {
		prometheusURL = "http://localhost:9090"
	}
	metricsReader := prometheus.NewMetricsReader(prometheus.NewClient(prometheusURL), time.Now)
	serverMetricsHandler := handler.NewServerMetricsHandler(metricsReader)

	routerErr := router.NewRouter(
		stockHandler,
		dailyPriceChangeHandler,
		dailyVolumeRatioHandler,
		weatherHandler,
		serverMetricsHandler,
	).Run()

	if routerErr != nil {
		log.Fatal(routerErr)
	}

}
