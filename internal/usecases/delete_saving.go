package usecases

import (
	"scrapper/internal/repositories"

	"github.com/jackc/pgx/v5"
)

func DeleteSaving(db *pgx.Conn, id int) error {
	err := repositories.NewListingRepository(db).Delete(id, "scrapped_listings")
	if err != nil {
		return err
	}
	return nil
}
