package query

import (
	"enterprise_core/internal/database"
	"enterprise_core/internal/model"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetReportQuery(c *gin.Context, db database.Service, userId int) ([]model.AccountReport, error) {
	var reports []model.AccountReport

	// Fetch all accounts for the user
	rows, err := db.Query(`
		SELECT a.id, a.user_id, a.name, a.balance, a.debit, a.credit
		FROM accounts a
		WHERE a.user_id = $1
	`, userId)
	if err != nil {
		return nil, fmt.Errorf("Database error: %v", err)
	}
	defer rows.Close()

	// Iterate through each account
	for rows.Next() {
		var account model.Account
		if err := rows.Scan(&account.ID, &account.UserID, &account.Name, &account.Balance, &account.Debit, &account.Credit); err != nil {
			return nil, err
		}

		// Fetch all transactions for the current account
		transactionRows, err := db.Query(`
			SELECT t.id, t.account_id, t.amount, t.method, t.status, t.created_at
			FROM transactions t
			WHERE t.account_id = $1
		`, account.ID)
		if err != nil {
			return nil, fmt.Errorf("Database error: %v", err)
		}
		defer transactionRows.Close()

		var transactions []model.Transaction
		for transactionRows.Next() {
			var transaction model.Transaction
			if err := transactionRows.Scan(&transaction.ID, &transaction.AccountID, &transaction.Amount, &transaction.Method, &transaction.Status, &transaction.CreatedAt); err != nil {
				return nil, err
			}
			transactions = append(transactions, transaction)
		}

		// Append the account and its related transactions to the report
		reports = append(reports, model.AccountReport{
			Account:     account,
			Transactions: transactions,
		})
	}

	// Check for any errors that occurred during the iteration of the accounts
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Database error: %v", err)
	}

	return reports, nil
}
