package cmd

import (
	bolt "github.com/coreos/bbolt"
	"log"
	"shopaholic/store/engine"
	"shopaholic/store/service"
	"time"
)

// UserCreateCommand set of flags and command for creation
type UserListCommand struct {
	Number int `short:"n" long:"number" description:"number of users" required:"false"`
}

func (ulc *UserListCommand) Execute(args []string) error {
	log.Printf("[INFO] users list command is started")
	b, err := engine.NewBoltDB(bolt.Options{Timeout: 30 * time.Second}, "shopaholic.db")
	if err != nil {
		return err
	}

	s := &service.DataStore{Interface: b}

	users, err := s.UserList(ulc.Number)
	if err != nil {
		return err
	}

	for _, user := range users {
		log.Printf("%+v", user)
	}

	return nil
}
