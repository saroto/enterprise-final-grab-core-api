package query

import (
	"enterprise_core/internal/database"
	"enterprise_core/internal/model"
	"fmt"
	"log"
)
func RegisterDriverQuery(db database.Service, id string, user *model.User) error {
	query := `
		UPDATE users 
		SET email = $1, name = $2, vehicleName = $3, vehicleNumber = $4, role = $5 
		WHERE id = $6
	`
	_, err := db.Exec(query, user.Email, user.Name, user.VehicleName, user.VehicleNumber, user.Role, id)

	// Check if there is an error
	if err != nil {
		log.Printf("Failed to update user %s: %v", id, err) // Print error to console
		return fmt.Errorf("error updating user %s: %w", id, err) // Return a wrapped error
	}

	return nil
}


