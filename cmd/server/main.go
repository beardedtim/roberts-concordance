package main

import (
	"log"
	http "mckp/roberts-concordance/http"
	monitoring "mckp/roberts-concordance/monitoring"
)

func main() {
	monitoring.Init()

	http.Create()

	err := http.Start(9999)

	if err != nil {
		log.Panic(err)
	}
}
