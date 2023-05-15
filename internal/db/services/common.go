package services

import "github.com/pauloo27/aryzona/internal/db"

func upsert[T any](entity T) error {
	// EAFP: easy to ask for forgiveness than permission,
	// let's try to update the user and if it fails, create it

	aff, err := db.DB.Update(entity)

	if aff == 1 || err != nil {
		return err
	}

	_, err = db.DB.Insert(entity)

	return err
}
