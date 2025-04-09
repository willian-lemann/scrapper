package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"scrapper/config"
	"scrapper/internal/structs"

	"github.com/jackc/pgx/v5"
)

type ListingRepository struct {
	DB *pgx.Conn
}

func NewListingRepository(db *pgx.Conn) *ListingRepository {
	return &ListingRepository{DB: db}
}

func (r *ListingRepository) GetAll(table string) ([]structs.ListingItem, error) {
	query := `SELECT id, link, image, address, price, area, bedrooms, type, "forSale", parking, content, photos, agency, bathrooms, ref, "placeholderImage" FROM scrapped_listings`
	rows, err := r.DB.Query(context.Background(), query, pgx.QueryExecModeSimpleProtocol)

	if err != nil {
		fmt.Println("Error executing query:")
		return nil, err
	}
	defer rows.Close()
	var listings []structs.ListingItem
	for rows.Next() {
		var listing structs.ListingItem
		err := rows.Scan(&listing.Id, &listing.Link, &listing.Image, &listing.Address, &listing.Price, &listing.Area, &listing.Bedrooms, &listing.Type, &listing.ForSale, &listing.Parking, &listing.Content, &listing.Photos, &listing.Agency, &listing.Bathrooms, &listing.Ref, &listing.PlaceholderImage)
		if err != nil {
			return nil, err
		}
		listings = append(listings, listing)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return listings, nil
}

func Create(listing *structs.ListingItem, table string) error {
	supabase := config.GetSupabaseClient()

	var results []structs.ListingItem

	err := supabase.DB.From(table).Insert(listing).Execute(&results)
	if err != nil {
		return err
	}

	return nil
}

func (r *ListingRepository) GetByID(id int, table string) (*structs.ListingItem, error) {
	query := fmt.Sprintf(`SELECT id, link, image, address, price, area, bedrooms, type, "forSale", parking, content, photos, agency, bathrooms, ref, "placeholderImage" FROM %s WHERE id = $1`, table)
	row := r.DB.QueryRow(context.Background(), query, id)

	var listing structs.ListingItem
	err := row.Scan(&listing.Id, &listing.Link, &listing.Image, &listing.Address, &listing.Price, &listing.Area, &listing.Bedrooms, &listing.Type, &listing.ForSale, &listing.Parking, &listing.Content, &listing.Photos, &listing.Agency, &listing.Bathrooms, &listing.Ref, &listing.PlaceholderImage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &listing, nil
}

func (r *ListingRepository) Update(listing *structs.ListingItem) error {
	query := `
		UPDATE listings
		SET link = $1, image = $2, address = $3, price = $4, area = $5, bedrooms = $6, type = $7, 'forSale' = $8, parking = $9, content = $10, photos = $11, agency = $12, bathrooms = $13, ref = $14, 'placeholderImage')' = $15, id = $16 
		WHERE id = $18
	`
	_, err := r.DB.Exec(context.Background(), query, listing.Link, listing.Image, listing.Address, listing.Price, listing.Area, listing.Bedrooms, listing.Type, listing.ForSale, listing.Parking, listing.Content, listing.Photos, listing.Agency, listing.Bathrooms, listing.Ref, listing.PlaceholderImage, listing.Id)
	return err
}

func (r *ListingRepository) Delete(id int, table string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", table)
	_, err := r.DB.Exec(context.Background(), query, id)
	return err
}
