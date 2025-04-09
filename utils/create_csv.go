package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"scrapper/internal/structs"
)

func CreateCSV(listings []structs.ListingItem) {
	file, err := os.Create("listings.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Link", "Image", "Address", "Price", "Area", "Bedrooms", "Bathrooms", "Type", "ForSale", "Parking", "Content"}
	err = writer.Write(header)
	if err != nil {
		log.Fatal("Cannot write header to file", err)
	}

	// // Write the data rows
	// for _, listing := range listings {
	// 	row := []string{
	// 		listing.Link,
	// 		listing.Image,
	// 		listing.Address,
	// 		listing.Price,
	// 		listing.Area,
	// 		listing.Bedrooms,
	// 		listing.Bathrooms,
	// 		listing.Type,
	// 		strconv.FormatBool(listing.ForSale),
	// 		listing.Parking,
	// 		listing.Content,
	// 		strconv.Itoa(len(listing.Photos)), // Convert the length of listing.Photos to a string
	// 	}
	// 	err = writer.Write(row)
	// 	if err != nil {
	// 		log.Fatal("Cannot write row to file", err)
	// 	}
	// }

	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("CSV file created successfully")
}
