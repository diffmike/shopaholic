package cmd

import (
	"log"
	"shopaholic/store"
	"shopaholic/utils"
)

type TransactionCreateCommand struct {
	Amount   float64 `short:"a" long:"amount" description:"amount of the transaction" required:"true"`
	Category string  `short:"c" long:"category" description:"category of the transaction" required:"false"`
	User     string  `short:"u" long:"user" description:"UID of the user" required:"true"`

	CommonOpts
}

func (tcc *TransactionCreateCommand) Execute(args []string) error {
	log.Printf("[INFO] new transation of %f command is started", tcc.Amount)
	user, err := tcc.Store.Details(tcc.User)
	if err != nil {
		return err
	}

	category, err := tcc.Store.FindCategoryByTitle(tcc.Category)
	if err != nil {
		return err
	}

	transaction := store.Transaction{
		User:     user,
		Category: category,
		Amount:   utils.Money{int64(tcc.Amount * 100), tcc.CommonOpts.Currency},
	}

	transactionID, err := tcc.Store.StoreTransaction(transaction)
	if err != nil {
		return err
	}

	log.Printf("[INFO] transaction was created with ID: %s", transactionID)
	return nil
}
