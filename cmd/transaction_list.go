package cmd

import (
	"log"
)

type TransactionListCommand struct {
	User string `short:"u" long:"user" description:"UID of the user" required:"true"`

	CommonOpts
}

func (tlc *TransactionListCommand) Execute(args []string) error {
	log.Printf("[INFO] transation list of %s command is started")

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
		log.Printf("ID %s. Balance was: %f, now: %f. Created: %s. Amount %f",
			transaction.ID, balanceWas, balanceNow, transaction.CreatedAt.String(), amount)
	}

	return nil
}
