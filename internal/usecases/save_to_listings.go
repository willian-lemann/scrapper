package usecases

import (
	"context"
	"scrapper/config"
	"scrapper/internal/repositories"
)

func SaveToListings() {
	db, err := config.DatabaseConfig()
	if err != nil {
		panic(err)
	}
	defer db.Close(context.Background())

	scrappedListings, err := repositories.NewListingRepository(db).GetAll("scrapped_listings")
	if err != nil {
		panic(err)
	}

	for _, listing := range scrappedListings {
		err := repositories.Create(&listing, "listings")
		if err != nil {
			panic(err)
		}
	}
}
