package main

import (
	"api/router"
	"log"
)

func main() {

	err := router.NewRouter().Run()

	if err != nil {
		log.Fatal(err)
	}

}
