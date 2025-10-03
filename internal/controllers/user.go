package controller

import (
	"enterprise_core/internal/database"
	"enterprise_core/internal/model"
	"enterprise_core/internal/query"
	"enterprise_core/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)


func UpdateUser(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := query.RegisterDriverQuery(db, id, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update driver"})
			return
		}
		utils.ResponseWithMessage(c, "Update User Successfully")
	}
}

func GetAllUsers(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fetch all users from the database
		users, err := query.GetAllUsersQuery(c,db)
		if err != nil {
			// Handle errors and send an appropriate response
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		// Send the response with the fetched data
		utils.ResponseWithData(c, "Get All Users Successfully", users)
	}
}
