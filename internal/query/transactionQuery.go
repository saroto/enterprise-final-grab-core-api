package query

import (
	"enterprise_core/internal/database"
	"enterprise_core/internal/model"
)

// func CreateTransactionQuery(db database.Service, userId int, transaction *model.Transaction) error {
// 	_, err := db.Exec("INSERT INTO transactions (account_id, amount, method, status, created_at) VALUES ($1, $2, $3, $4, $5)", transaction.AccountID, transaction.Amount, transaction.Method, transaction.Status, transaction.CreatedAt)
// 	return err
// }

func CreateTransactionQuery(db database.Service, userId int, transaction *model.Transaction) error {
	// Start a transaction
	tx, err := db.Begin("BEGIN")
	if err != nil {
		return err
	}

	// Ensure rollback in case of failure
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get the account ID and role based on user ID
	var accountID int
	var role string
	err = tx.QueryRow("SELECT id, role FROM users WHERE id = $1", userId).Scan(&accountID, &role)
	if err != nil {
		return err
	}

	// Assign the retrieved account ID to the transaction
	transaction.AccountID = accountID

	// Determine the final transaction amount based on role
	amount := transaction.Amount
	if role == "user" {
		amount = -amount // Subtract if the user is a regular user
	}

	// Insert the transaction into the transactions table
	_, err = tx.Exec(
		"INSERT INTO transactions (account_id, amount, method, status, created_at) VALUES ($1, $2, $3, $4, $5)",
		transaction.AccountID, transaction.Amount, transaction.Method, transaction.Status, transaction.CreatedAt,
	)
	if err != nil {
		return err
	}

	// Update the account balance based on role
	_, err = tx.Exec(
		"UPDATE accounts SET balance = balance + $1 WHERE id = $2",
		amount, accountID,
	)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	return err
}

func GetTransactionsQuery(db database.Service) ([]model.Transaction, error) {
	rows, err := db.Query("SELECT id, account_id, amount, method, status, created_at FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.AccountID, &transaction.Amount, &transaction.Method, &transaction.Status, &transaction.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func GetTransactionQuery(db database.Service, id string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := db.QueryRow("SELECT id, account_id, amount, method, status, created_at FROM transactions WHERE id = $1", id).
		Scan(&transaction.ID, &transaction.AccountID, &transaction.Amount, &transaction.Method, &transaction.Status, &transaction.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func UpdateTransactionQuery(db database.Service, id string, transaction *model.Transaction) error {
	query := `
		UPDATE transactions 
		SET account_id = $1, amount = $2, method = $3, status = $4, created_at = $5 
		WHERE id = $6
	`
	_, err := db.Exec(query, transaction.AccountID, transaction.Amount, transaction.Method, transaction.Status, transaction.CreatedAt, id)
	return err
}

func DeleteTransactionQuery(db database.Service, id string) error {
	_, err := db.Exec("DELETE FROM transactions WHERE id = $1", id)
	return err
}
