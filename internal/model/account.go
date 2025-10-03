package model

type Account struct {
	ID      uint    `json:"id"`
	UserID  uint    `json:"user_id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	Debit   float64 `json:"debit"`
	Credit  float64 `json:"credit"`
}


// AccountReport represents the combined report of accounts and transactions for a user.
type AccountReport struct {
	Account     Account     `json:"account"`
	Transactions []Transaction `json:"transactions"`
}