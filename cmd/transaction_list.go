package cmd

import (
	"log"
)

type TransactionListCommand struct {
	User string `short:"u" long:"user" description:"UID of the user" required:"true"`

	CommonOpts
}

func (tlc *TransactionListCommand) Execute(args []string) error {
	log.Printf("[INFO] transation list command is started")

	user, err := tlc.Store.Details(tlc.User)
	if err != nil {
		return err
	}

	transactions, err := tlc.Store.List(user)
	if err != nil {
		return err
	}

	for _, transaction := range transactions {
		amount := float64(transaction.Amount.Amount / 100)
		balanceWas := float64(transaction.BalanceWas.Amount / 100)
		balanceNow := float64(transaction.BalanceNow.Amount / 100)
		time := transaction.CreatedAt.Format("02.01.2006 15:04:05")
		log.Printf("ID %s. Balance %6.2f → %6.2f. 🕰 %s. %6.2f💲 at %s",
			transaction.ID, balanceWas, balanceNow, time, amount, transaction.Category.Title)
	}

	if len(transactions) == 0 {
		log.Print("There no transactions in the DataStore for the %s 🤷", user.Name)
	}

	return nil
}
