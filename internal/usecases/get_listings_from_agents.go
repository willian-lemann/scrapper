package usecases

import (
	"context"
	"log"
	"scrapper/config"
	"scrapper/internal/structs"
)

func GetListingsFromAgents() ([]structs.ListingItem, error) {
	DB, err := config.DatabaseConfig()

	rows, err := DB.Query(context.Background(), "SELECT id, name, price FROM listings")
	if err != nil {
		log.Fatal("Failed to retrieve listings:", err)
		return nil, err
	}
	defer rows.Close()

	var listings []structs.ListingItem

	for rows.Next() {
		var listing structs.ListingItem

		err := rows.Scan(&listings)
		if err != nil {
			log.Fatal("Failed to scan row:", err)
		}
		listings = append(listings, listing)
	}

	return listings, nil
}
