package model
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	VehicleName   *string `json:"vehicleName,omitempty"`
	VehicleNumber *string `json:"vehicleNumber,omitempty"`
	Role          *string `json:"role,omitempty"`
}