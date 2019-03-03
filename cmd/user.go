package cmd

import (
	"log"
)

// UserCommand set of flags and command for creation
type UserCommand struct {
	Name string `short:"n" long:"name" description:"new user name" required:"true"`
}

func (uc *UserCommand) Execute(args []string) error {
	log.Printf("[INFO] user %s creating in progress", uc.Name)

	log.Printf("[INFO] user %s was created", uc.Name)
	return nil
}
