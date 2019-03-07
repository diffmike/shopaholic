package cmd

import (
	"log"
)

// UserCreateCommand set of flags and command for creation
type UserListCommand struct {
	Number int `short:"n" long:"number" description:"number of users" required:"false"`

	CommonOpts
}

func (ulc *UserListCommand) Execute(args []string) error {
	log.Printf("[INFO] users list command is started")

	users, err := ulc.Store.UserList(ulc.Number)
	if err != nil {
		return err
	}

	for _, user := range users {
		balance := float64(user.Balance.Amount / 100)
		created := user.CreatedAt.Format("02.01.2006 15:04:05")
		log.Printf("ID: %s. Name: %10s. Created: %s. Balance: %6.2f", user.ID, user.Name, created, balance)
	}

	return nil
}
