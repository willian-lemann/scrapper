package main

import (
	"fmt"
	"scrapper/internal/scrappers"
	"scrapper/internal/usecases"
	"scrapper/utils"
	"sync"
	"time"
)

func main() {
	utils.TruncateTable("scrapped_listings")

	var scraperWG sync.WaitGroup
	scraperWG.Add(1)

	start := time.Now()

	go scrappers.ExecuteJefersonAlba(&scraperWG)
	scrappers.ExecuteAuxPredial()

	scraperWG.Wait()

	elapsed := time.Since(start)

	utils.TruncateTable("listings")

	usecases.SaveToListings()

	fmt.Println("Scraping finished")
	fmt.Println("It took:", elapsed.Seconds())
}
