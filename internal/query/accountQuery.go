package query

import (
	"enterprise_core/internal/database"
	"enterprise_core/internal/model"
	"fmt"
)

// CreateAccount creates a new account

func CreateAccountQuery(db database.Service, userId int, account *model.Account) error {

	fmt.Println("CreateAccountQuery")
	_, err := db.Exec("INSERT INTO accounts (user_id, name, balance, debit, credit) VALUES ($1, $2, $3, $4, $5)", userId, account.Name, account.Balance, account.Debit, account.Credit)

	return err
}
// GetAccounts fetches all accounts
func GetAccountsQuery(db database.Service) ([]model.Account, error) {
	rows, err := db.Query("SELECT id, user_id, name, balance, debit, credit FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts [] model.Account
	for rows.Next() {
		var account  model.Account
		if err := rows.Scan(&account.ID, &account.UserID, &account.Name, &account.Balance, &account.Debit, &account.Credit); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func GetAccountQuery(db database.Service, id string) (*model.Account, error) {
	var account model.Account
	err := db.QueryRow("SELECT id, user_id, name, balance, debit, credit FROM accounts WHERE id = $1", id).Scan(&account.ID, &account.UserID, &account.Name, &account.Balance, &account.Debit, &account.Credit)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func GetOwnAccountQuery(db database.Service, id int) (*model.Account, error) {
	var account model.Account
	err := db.QueryRow("SELECT id, user_id, name, balance, debit, credit FROM accounts WHERE user_id = $1", id).Scan(&account.ID, &account.UserID, &account.Name, &account.Balance, &account.Debit, &account.Credit)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func UpdateAccountQuery(db database.Service, id string, account * model.Account) error {
	query := "UPDATE accounts SET name = $1, balance = $2, debit = $3, credit = $4 WHERE id = $5"
	_, err := db.Exec(query, account.Name, account.Balance, account.Debit, account.Credit, id)
	return err
}

func DeleteAccountQuery(db database.Service, id string) error {
	_, err := db.Exec("DELETE FROM accounts WHERE id = $1", id)
	return err
}