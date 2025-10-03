package model

import "time"

type Transaction struct {
	ID        int       `json:"id"`          // Primary key
	AccountID int       `json:"account_id"`  // Foreign key referencing the accounts table
	Amount    float64   `json:"amount"`      // The amount of the transaction
	Method    string    `json:"method"`      // The method of the transaction (e.g., "credit", "debit")
	Status    string    `json:"status"`      // The status of the transaction (e.g., "pending", "completed")
	CreatedAt time.Time `json:"created_at"`  // Timestamp of when the transaction was created
}