package types

type Transaction struct {
	ID     int
	Amount int
	Status string
}

type Account struct {
	ID           int
	Balance      int
	Transactions []Transaction
}
