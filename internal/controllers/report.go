package controller

import (
	"enterprise_core/internal/database"
	"enterprise_core/internal/query"
	"enterprise_core/internal/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)


func GetUserReport(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user ID from context (auth middleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Extract the user ID from the query parameter (if provided)
		userIdStr := c.DefaultQuery("id", fmt.Sprintf("%d", userID)) // Default to userID from context if not provided
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Fetch report data using the userId (either from query parameter or authenticated user)
		data, err := query.GetReportQuery(c, db, userId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
			return
		}

		// Send the report as response
		utils.ResponseWithData(c, "Get Report Successfully", gin.H{"report": data})
	}
}
