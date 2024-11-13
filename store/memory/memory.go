package memorystore

import (
	"fmt"

	"github.com/gkay21/kttipay/types"
)

type MemoryStore struct {
	Accounts []*types.Account
}

// New initializes a new MemoryStore with no accounts.
func New() (*MemoryStore, error) {
	return &MemoryStore{
		Accounts: []*types.Account{},
	}, nil
}

func (ms *MemoryStore) ValidateAccount(accountID int) error {
	if accountID < 1 || accountID > len(ms.Accounts) {
		return fmt.Errorf("account with ID %d not found", accountID)
	}

	return nil
}

func (ms *MemoryStore) ValidateTransaction(accountID int, transactionID int) error {
	if err := ms.ValidateAccount(accountID); err != nil {
		return err
	}

	if transactionID < 1 || transactionID > len(ms.Accounts[accountID-1].Transactions) {
		return fmt.Errorf("account with ID %d does not have transaction with ID %d", accountID, transactionID)
	}

	return nil
}

// CreateAccount creates a new account with a zero balance and an empty transaction list.
func (ms *MemoryStore) CreateAccount() (*types.Account, error) {
	account := &types.Account{
		Balance:      0,
		Transactions: []types.Transaction{},
	}
	ms.Accounts = append(ms.Accounts, account)
	return account, nil
}

// GetAccount retrieves an account by its ID. It returns an error if the account is not found.
func (ms *MemoryStore) GetAccount(accountID int) (*types.Account, error) {
	if err := ms.ValidateAccount(accountID); err != nil {
		return nil, err
	}
	return ms.Accounts[accountID], nil
}

// CreateTransaction creates a new pending transaction for a specified account, by account ID.
func (ms *MemoryStore) CreateTransaction(accountID int, amount int) (*types.Transaction, error) {
	if err := ms.ValidateAccount(accountID); err != nil {
		return nil, err
	}

	transaction := &types.Transaction{
		Amount: amount,
		Status: "pending",
	}

	ms.Accounts[accountID-1].Transactions = append(ms.Accounts[accountID-1].Transactions, *transaction)
	return transaction, nil
}

// SettleTransaction settles an existing transaction
func (ms *MemoryStore) SettleTransaction(accountID int, transactionID int) error {
	if err := ms.ValidateTransaction(accountID, transactionID); err != nil {
		return err
	}

	ms.Accounts[accountID-1].Transactions[transactionID-1].Status = "settled"
	return nil
}

// RefundTransaction refunds part or all of a pending transaction
func (ms *MemoryStore) RefundTransaction(accountID int, transactionID int, amount int) error {
	if err := ms.ValidateTransaction(accountID, transactionID); err != nil {
		return err
	}

	transaction := ms.Accounts[accountID-1].Transactions[transactionID-1]

	if transaction.Status != "pending" {
		return fmt.Errorf("unable to refund non-pending transactions")
	}

	if transaction.Amount < amount {
		return fmt.Errorf("refund amount exceeds pending transaction amount")
	}

	ms.Accounts[accountID-1].Transactions[transactionID-1].Amount -= amount

	return nil
}

// ListTransactions returns the list of all transactions
func (ms *MemoryStore) ListTransactions(accountID int) ([]types.Transaction, error) {
	if err := ms.ValidateAccount(accountID); err != nil {
		return nil, err
	}

	return ms.Accounts[accountID-1].Transactions, nil
}

// CalculateBalance calculates the balance for a specified account by account ID.
func (ms *MemoryStore) CalculateBalance(accountID int) (int, error) {
	if err := ms.ValidateAccount(accountID); err != nil {
		return 0, err
	}

	balance := 0
	for _, transaction := range ms.Accounts[accountID-1].Transactions {
		balance += transaction.Amount
	}

	ms.Accounts[accountID-1].Balance = balance
	return balance, nil
}
