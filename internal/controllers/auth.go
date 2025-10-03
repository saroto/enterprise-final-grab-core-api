package controller

import (
	"encoding/json"
	"enterprise_core/internal/database"
	"net/http"

	"enterprise_core/internal/utils"

	"enterprise_core/internal/query"

	"github.com/gin-gonic/gin"
)


func AuthTest() gin.HandlerFunc {
	return func (c *gin.Context) {
		resp := make(map[string]string)
		resp["message"] = "Authentication healthy"

		c.JSON(http.StatusOK, resp)
	}
}

func RegisterUser(db database.Service) gin.HandlerFunc {
	return func (c *gin.Context) {

		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Name 	 string `json:"name"`
		}

		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		err := query.RegisterUserQuery(req.Email, req.Name, req.Password, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
}

func LoginUser(db database.Service) gin.HandlerFunc {
	return func (c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		token, expiryDate, err := query.LoginUserQuery(req.Email, req.Password, db)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		//sync to subsystem
		syncMessage, syncedURLs := SyncSystem(c, token) 
		utils.ResponseWithData(c, "Login successfully", gin.H{
			"user": gin.H{
				"email": req.Email,
			},
			"tokens":       gin.H{
				"access_token": token,
				"expiry_date": expiryDate,
			},
			"synced" : gin.H{
			"sync_message": syncMessage,
			"synced_urls":  syncedURLs,
			},
	
		})
	}
}

func LogoutUser(db database.Service) gin.HandlerFunc {
	return func (c *gin.Context) {

		query.LogoutQuery(c, db)

		utils.ResponseWithMessage(c, "Logout successful")
	}
}

func GetMe(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user ID from context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		data, err := query.GetMeQuery(c, db, userID.(int))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		utils.ResponseWithData(c, "Get Profile Successfully", gin.H{"user": data})
	}
}




