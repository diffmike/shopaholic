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
		log.Printf("%+v", user)
	}

	return nil
}
