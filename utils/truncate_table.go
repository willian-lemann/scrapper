package utils

import (
	"context"
	"fmt"
	"log"
	"scrapper/config"
)

func TruncateTable(table string) {
	fmt.Println("Clearing database...")

	db, err := config.DatabaseConfig()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to retrieve listings: %v", err)
	}

	query := fmt.Sprintf("DELETE FROM %s", table)
	_, err = db.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}

	db.Close(context.Background())

	fmt.Printf("%s truncated.\n", table)
}
