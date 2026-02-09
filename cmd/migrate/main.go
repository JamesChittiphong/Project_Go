package main

import (
	"Backend_Go/internal/config"
	"Backend_Go/internal/entities"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// load env
	if err := godotenv.Load(".env.config"); err != nil {
		log.Println("Warning: .env.config not found, relying on system env")
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting Migration Backfill...")

	// 1. Backfill Dealer Status
	// If IsApproved is true, set Status = 'approved'.
	// If IsApproved is false, set Status = 'pending' (and ensure IsApproved is false).
	// We only touch records where Status is empty to avoid overwriting new data.
	result := db.Model(&entities.Dealer{}).
		Where("status = '' OR status IS NULL").
		Updates(map[string]interface{}{
			"status": db.Raw("CASE WHEN is_approved = true THEN 'approved' ELSE 'pending' END"),
		})
	if result.Error != nil {
		log.Printf("Error backfilling dealers: %v\n", result.Error)
	} else {
		fmt.Printf("Backfilled %d dealers.\n", result.RowsAffected)
	}

	// 2. Backfill Car Status
	// Existing cars should be 'approved' so they remain visible?
	// Or 'pending'?
	// Strategy: Set all existing cars to 'approved' to avoid hiding everything on production update.
	// New cars created after this code deploy will be 'pending' by default in code.
	resultCars := db.Model(&entities.Car{}).
		Where("status = '' OR status IS NULL").
		Update("status", "approved")

	if resultCars.Error != nil {
		log.Printf("Error backfilling cars: %v\n", resultCars.Error)
	} else {
		fmt.Printf("Backfilled %d cars to 'approved'.\n", resultCars.RowsAffected)
	}

	fmt.Println("Migration Complete.")
}
