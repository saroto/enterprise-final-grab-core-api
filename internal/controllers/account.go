package controller

import (
	"encoding/json"
	"enterprise_core/internal/database"
	"enterprise_core/internal/model"
	"enterprise_core/internal/utils"
	"net/http"

	"enterprise_core/internal/query"

	"github.com/gin-gonic/gin"
)

func CreateAccount(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}

		var account model.Account
		
		if err := json.NewDecoder(c.Request.Body).Decode(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		err := query.CreateAccountQuery(db,  userID.(int), &account)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Creat Account failed"})
			return
		}
		
		utils.ResponseWithMessage(c, "Create Account Successfully")
	}
}

func GetAllAccounts(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		accounts, err := query.GetAccountsQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
			return
		}
		utils.ResponseWithData(c, "Get All Accounts Successfully", accounts)
	}
}

func GetOwnAccount(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		accounts, err := query.GetOwnAccountQuery(db, userID.(int))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
			return
		}
		utils.ResponseWithData(c, "Get Own Accounts Successfully", accounts)
	}
}
func GetAccount(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		account, err := query.GetAccountQuery(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		utils.ResponseWithData(c, "Get Account Successfully", account)
	}
}

func UpdateAccount(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var account model.Account
		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := query.UpdateAccountQuery(db, id, &account)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
			return
		}
		utils.ResponseWithMessage(c, "Update Account Successfully")
	}
}

func DeleteAccount(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := query.DeleteAccountQuery(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
			return
		}
		utils.ResponseWithMessage(c, "Delete Account Successfully")
	}
}