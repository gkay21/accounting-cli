package immudbstore

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/codenotary/immudb/pkg/stdlib"
	"github.com/gkay21/kttipay/types"
)

type ImmudbStore struct {
	db *sql.DB
}

// New initializes a new MemoryStore with no accounts.
func New() (*ImmudbStore, error) {
	connStr := "immudb://immudb:immudb@127.0.0.1:3322/defaultdb?sslmode=disable"

	immudb, err := sql.Open("immudb", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = immudb.Exec(
		"CREATE TABLE IF NOT EXISTS account (id INTEGER AUTO_INCREMENT, balance INTEGER NOT NULL, PRIMARY KEY id);",
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = immudb.Exec(
		"CREATE TABLE IF NOT EXISTS transactions (id INTEGER AUTO_INCREMENT, account_id INTEGER NOT NULL, status VARCHAR[20] NOT NULL, amount INTEGER NOT NULL, PRIMARY KEY id);",
	)
	if err != nil {
		log.Fatal(err)
	}

	return &ImmudbStore{
		db: immudb,
	}, nil
}

// CreateAccount creates a new account with a zero balance and an empty transaction list.
func (is *ImmudbStore) CreateAccount() (*types.Account, error) {
	q := "INSERT INTO account(balance) VALUES (0);"

	res, err := is.db.Exec(q)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &types.Account{ID: int(id)}, nil
}

// GetAccount retrieves an account by its ID. It returns an error if the account is not found.
func (is *ImmudbStore) GetAccount(accountID int) (*types.Account, error) {
	q := "SELECT * FROM account WHERE id=$1"

	var id, balance int
	err := is.db.QueryRow(q, accountID).Scan(&id, &balance)
	if err != nil {
		return nil, err
	}

	return &types.Account{
		ID:      id,
		Balance: balance,
	}, nil
}

// UpdateAccount updates an account by its ID. It returns an error if the account is not found.
func (is *ImmudbStore) UpdateAccount(accountID int, balance int) (*types.Account, error) {
	q := "UPDATE account SET balance=$2 WHERE id=$1"

	_, err := is.db.Exec(q, accountID, balance)
	if err != nil {
		return nil, err
	}

	return &types.Account{
		ID:      accountID,
		Balance: balance,
	}, nil
}

// CreateTransaction creates a new pending transaction for a specified account, by account ID.
func (is *ImmudbStore) CreateTransaction(accountID int, amount int) (*types.Transaction, error) {
	account, err := is.GetAccount(accountID)
	if err != nil {
		return nil, err
	}

	query := "INSERT INTO transactions(account_id,status,amount) VALUES ($1,$2,$3);"
	res, err := is.db.Exec(query, account.ID, "pending", amount)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	balance, err := is.CalculateBalance(accountID)
	if err != nil {
		return nil, err
	}

	_, err = is.UpdateAccount(accountID, balance)
	if err != nil {
		return nil, err
	}

	return &types.Transaction{
		ID:     int(id),
		Amount: amount,
		Status: "pending",
	}, nil

}

// SettleTransaction settles an existing transaction
func (is *ImmudbStore) SettleTransaction(accountID int, transactionID int) error {
	q := "UPDATE transactions SET status=$3 WHERE id=$1 AND account_id=$2"

	_, err := is.db.Exec(q, transactionID, accountID, "settled")
	if err != nil {
		return err
	}

	return nil
}

// RefundTransaction refunds part or all of a pending transaction.
func (is *ImmudbStore) RefundTransaction(accountID int, transactionID int, amount int) error {
	transaction, err := is.GetTransaction(accountID, transactionID)
	if err != nil {
		return err
	}

	if transaction.Status == "settled" {
		return fmt.Errorf("transaction %d is already settled", transaction.ID)
	}
	if transaction.Amount < amount {
		return fmt.Errorf("transaction amount %d is less than refund amount %d", transaction.Amount, amount)
	}

	q := "UPDATE transactions SET amount=$2 WHERE id=$1;"

	_, err = is.db.Exec(q, transactionID, transaction.Amount-amount)
	if err != nil {
		return err
	}

	return nil

}

// GetTransaction returns the transaction by ID
func (is *ImmudbStore) GetTransaction(accountID int, transactionID int) (types.Transaction, error) {
	q := "SELECT id,status,amount FROM transactions WHERE id=$1 AND account_id=$2;"

	var transaction types.Transaction

	err := is.db.QueryRow(q, transactionID, accountID).Scan(&transaction.ID, &transaction.Status, &transaction.Amount)
	if err != nil {
		return types.Transaction{}, err
	}

	return transaction, nil
}

// ListTransactions returns the list of all transactions
func (is *ImmudbStore) ListTransactions(accountID int) ([]types.Transaction, error) {
	q := "SELECT id,status,amount FROM transactions WHERE account_id=$1;"

	var transactions []types.Transaction

	rows, err := is.db.Query(q, accountID)
	if err != nil {
		return transactions, err
	}

	for rows.Next() {
		var id, amount int
		var status string
		err := rows.Scan(&id, &status, &amount)
		if err != nil {
			return transactions, err
		}
		transactions = append(transactions, types.Transaction{
			ID:     id,
			Status: status,
			Amount: amount,
		})
	}

	return transactions, nil
}

// CalculateBalance calculates the balance for a specified account by account ID.
func (is *ImmudbStore) CalculateBalance(accountID int) (int, error) {
	q := "SELECT amount FROM transactions WHERE account_id=$1;"

	rows, err := is.db.Query(q, accountID)
	if err != nil {
		return 0, err
	}

	var balance int
	for rows.Next() {
		var amount int
		err := rows.Scan(&amount)
		if err != nil {
			return 0, err
		}
		balance += amount
	}

	return balance, nil
}
