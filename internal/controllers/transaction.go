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

func CreateTransaction(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var transaction model.Transaction

		if err := json.NewDecoder(c.Request.Body).Decode(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		err := query.CreateTransactionQuery(db, userID.(int), &transaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Create transaction failed"})
			return
		}

		utils.ResponseWithMessage(c, "Create Transaction Successfully")
	}
}


func GetAllTransactions(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		transactions, err := query.GetTransactionsQuery(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
			return
		}
		utils.ResponseWithData(c, "Get All Transactions Successfully", transactions)
	}
}
func GetTransaction(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		transaction, err := query.GetTransactionQuery(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}
		utils.ResponseWithData(c, "Get Transaction Successfully", transaction)
	}
}

func UpdateTransaction(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var transaction model.Transaction
		if err := c.ShouldBindJSON(&transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update the transaction
		err := query.UpdateTransactionQuery(db, id, &transaction)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
			return
		}

		utils.ResponseWithMessage(c, "Update Transaction Successfully")
	}
}


func DeleteTransaction(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := query.DeleteTransactionQuery(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
			return
		}
		utils.ResponseWithMessage(c, "Delete Transaction Successfully")
	}
}
