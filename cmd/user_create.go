package cmd

import (
	"log"
	"shopaholic/store"
	"shopaholic/utils"
)

// UserCreateCommand set of flags and command for creation
type UserCreateCommand struct {
	Name    string `short:"n" long:"name" description:"new user name" required:"true"`
	Balance int64  `short:"b" long:"balance" description:"new user balance" default:"0"`

	CommonOpts
}

func (ucc *UserCreateCommand) Execute(args []string) error {
	log.Printf("[INFO] user %s creating command is started", ucc.Name)
	user := store.User{
		Name:    ucc.Name,
		Balance: utils.Money{ucc.Balance, ucc.CommonOpts.Currency},
	}

	userID, err := ucc.Store.Register(user)
	if err != nil {
		return err
	}

	log.Printf("[INFO] user %s was created with ID: %s", ucc.Name, userID)
	return nil
}
