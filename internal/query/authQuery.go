package query

import (
	"database/sql"
	"enterprise_core/internal/database"
	"fmt"
	"time"

	"enterprise_core/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("your-secret-key")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func RegisterUserQuery(email string, name string, password string, db database.Service) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (email, name, password) VALUES ($1, $2, $3)", email, name, string(hashedPassword))
	return err
}

func LoginUserQuery(email, password string, db database.Service) (string, time.Time, error) {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE email=$1", email).Scan(&hashedPassword)
	if err != nil {
		return "", time.Time{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid credentials")
	}

	return GetToken(email, db)
}

func LogoutQuery(c *gin.Context,db database.Service) error {
	userID, exists := c.Get("userID")
		if !exists {
			return fmt.Errorf("Unauthorized") 
		}
	// Execute the update query to revoke the token for the user
	_, err := db.Exec("UPDATE tokens SET revoked=TRUE WHERE user_id=$1", userID)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %v", err)
	}

	return nil
}
func GetToken(email string, db database.Service) (string, time.Time, error) {
	var userID int
	var token string
	var expiry time.Time

	err := db.QueryRow("SELECT id FROM users WHERE email=$1", email).Scan(&userID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get user_id: %v", err)
	}

	err = db.QueryRow("SELECT token, created_at FROM tokens WHERE user_id=$1 AND revoked=FALSE ORDER BY created_at DESC LIMIT 1", userID).
	Scan(&token, &expiry)

	if err == nil && expiry.After(time.Now()) {
		return token, expiry, nil
	}

	expiryDate := time.Now().Add(24 * time.Hour) 
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		// UserID: email,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryDate),
		},
	})

	signedToken, err := newToken.SignedString(secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %v", err)
	}

	_, err = db.Exec("INSERT INTO tokens (user_id, token, created_at) VALUES ($1, $2, $3)", userID, signedToken, time.Now())
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to store token: %v", err)
	}

	return signedToken, expiryDate, nil
}

func GetMeQuery(c *gin.Context,db database.Service, userId int) (*model.User, error)  {

	var user model.User

	println("userId", userId)

	err := db.QueryRow("SELECT id, email, name, vehicleName, vehicleName, role FROM users WHERE id = $1", userId).
		Scan(&user.ID, &user.Email, &user.Name, &user.VehicleName, &user.VehicleNumber, &user.Role)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("User not found")
			}
			return nil, fmt.Errorf("Database error: %v", err)
		}
		
	return &user, nil
}

func GetSSOQuery(c *gin.Context, db database.Service, userId int) (*model.User, string, time.Time, error) {

	// Variables to store user details
	var user model.User
	var token string
	var expiry time.Time

	// Get user details from the database by userId
	err := db.QueryRow("SELECT id, email, name, role, vehicleNumber, vehicleName FROM users WHERE id = $1", userId).
		Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.VehicleNumber, &user.VehicleName)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", time.Time{}, fmt.Errorf("User not found")
		}
		return nil, "", time.Time{}, fmt.Errorf("Database error: %v", err)
	}

	// Call GetToken to retrieve or generate the token
	token, expiry, err = GetToken(user.Email, db)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("Failed to get token: %v", err)
	}

	// Return the user details along with the token and its expiry time
	return &user, token, expiry, nil
}

func GetAllUsersQuery(c *gin.Context, db database.Service) ([]*model.User, error) {
	var users []*model.User

	// Query to get all users
	rows, err := db.Query("SELECT id, email, name, vehicleName, vehicleNumber, role FROM users")
	if err != nil {
		return nil, fmt.Errorf("Database error: %v", err)
	}
	defer rows.Close()

	// Iterate over rows and scan them into the users slice
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.VehicleName, &user.VehicleNumber, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("Error scanning user data: %v", err)
		}
		users = append(users, &user)
	}

	// Check for any errors encountered while iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %v", err)
	}

	// Return the list of users
	return users, nil
}
