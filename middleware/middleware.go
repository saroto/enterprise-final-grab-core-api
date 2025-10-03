package middleware

import (
	"bytes"
	"enterprise_core/internal/database"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Claims struct for JWT claims
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// Secret key to validate JWT token
var secretKey = []byte("your-secret-key")

// AuthMiddleware is the JWT authentication middleware for Gin
func AuthMiddleware(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("authHeader: ", authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)

		isRevoked, err := IsTokenRevoked(tokenString, db)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			c.Abort()
			return
		}

		if isRevoked {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token revoked"})
			c.Abort()
			return
		}

		c.Next()
	}
}
func IsTokenRevoked(token string, db database.Service) (bool, error) {
	var revoked bool
	err := db.QueryRow("SELECT revoked FROM tokens WHERE token = $1", token).Scan(&revoked)
	if err != nil {
		return false, fmt.Errorf("failed to check token revocation: %v", err)
	}
	return revoked, nil
}

func LogRequestMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("===== Incoming Request =====")
        fmt.Println("Method:", c.Request.Method)
        fmt.Println("URL:", c.Request.URL.Path)
        fmt.Println("Headers:")
        
        for key, values := range c.Request.Header {
            fmt.Printf("%s: %s\n", key, values)
        }

        // Read Body (if JSON)
        if c.Request.Method == "POST" || c.Request.Method == "PUT" {
            bodyBytes, _ := io.ReadAll(c.Request.Body) // Read the body
            fmt.Println("Body:", string(bodyBytes))

            // Reset the body so it can be read again in the actual handler
            c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
        }

        fmt.Println("===========================")

        c.Next() // Continue to the next middleware or handler
    }
}
