package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gkay21/kttipay/store"
)

func printChoices() {
	fmt.Print(`Choose an option:
1. Create Account
2. Create Pending Transaction
3. Settle Transaction
4. Refund Transaction
5. Calculate Available Balance
6. List Transactions
7. Exit
Enter choice: `)
}

func main() {
	var storeInput string
	flag.StringVar(&storeInput, "store", "memory", "choose where to store data ['memory', 'database'] (default: memory)")
	flag.Parse()

	s, err := store.New(storeInput)
	if err != nil {

	}

	for {
		printChoices()
		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			_, err := s.CreateAccount()
			if err != nil {
				fmt.Printf("something went wrong: %s\n", err.Error())
				continue
			}
			fmt.Println("Successfully created account")
		case 2:
			var accountID, amount int
			fmt.Print("Enter account ID: ")
			fmt.Scan(&accountID)
			fmt.Print("Enter transaction amount (in cents): ")
			fmt.Scan(&amount)
			transaction, err := s.CreateTransaction(accountID, amount)
			if err != nil {
				fmt.Printf("something went wrong: %s\n", err.Error())
				continue
			}
			fmt.Printf("Successfully created transaction %d\n", transaction.ID)
		case 3:
			var accountID, transactionID int
			fmt.Print("Enter account ID: ")
			fmt.Scan(&accountID)
			fmt.Print("Enter transaction ID: ")
			fmt.Scan(&transactionID)
			err := s.SettleTransaction(accountID, transactionID)
			if err != nil {
				fmt.Printf("something went wrong: %s\n", err.Error())
				continue
			}
			fmt.Printf("Successfully settled transaction %d\n", transactionID)
		case 4:
			var accountID, transactionID, refundAmount int
			fmt.Print("Enter account ID: ")
			fmt.Scan(&accountID)
			fmt.Print("Enter transaction ID: ")
			fmt.Scan(&transactionID)
			fmt.Print("Enter refund amount (in cents): ")
			fmt.Scan(&refundAmount)
			err := s.RefundTransaction(accountID, transactionID, refundAmount)
			if err != nil {
				fmt.Printf("something went wrong: %s\n", err.Error())
				continue
			}
			fmt.Printf("Successfully refunded %d for transaction %d\n", refundAmount, transactionID)
		case 5:
			var accountID int
			fmt.Print("Enter account ID: ")
			fmt.Scan(&accountID)
			balance, err := s.CalculateBalance(accountID)
			if err != nil {
				fmt.Printf("something went wrong: %s\n", err.Error())
				continue
			}
			fmt.Printf("Balance: %d\n", balance)
		case 6:
			var accountID int
			fmt.Print("Enter account ID: ")
			fmt.Scan(&accountID)
			transactions, err := s.ListTransactions(accountID)
			if err != nil {
				fmt.Printf("something went wrong: %s\n", err.Error())
				continue
			}
			for _, t := range transactions {
				fmt.Printf("%+v\n", t)
			}
		case 7:
			fmt.Println("Exiting.")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Try again.")
		}
		fmt.Println()
	}
}
