package store

import (
	"fmt"

	immudbstore "github.com/gkay21/kttipay/store/immudb"
	memorystore "github.com/gkay21/kttipay/store/memory"
	"github.com/gkay21/kttipay/types"
)

type Store interface {
	CreateAccount() (*types.Account, error)
	GetAccount(accountID int) (*types.Account, error)
	CreateTransaction(accountID int, amount int) (*types.Transaction, error)
	SettleTransaction(accountID int, transactionID int) error
	RefundTransaction(accountID int, transactionID int, amount int) error
	ListTransactions(accountID int) ([]types.Transaction, error)
	CalculateBalance(accountID int) (int, error)
}

func New(store string) (Store, error) {
	switch store {
	case "memory":
		return memorystore.New()
	case "database":
		return immudbstore.New()
	}

	return nil, fmt.Errorf("invalid store passed in")
}
