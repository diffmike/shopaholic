package cmd

import (
	bolt "github.com/coreos/bbolt"
	"log"
	"shopaholic/store"
	"shopaholic/store/engine"
	"shopaholic/store/service"
	"time"
)

// UserCommand set of flags and command for creation
type UserCommand struct {
	Name    string `short:"n" long:"name" description:"new user name" required:"true"`
	Balance int    `short:"b" long:"balance" description:"new user balance" required:"false"`
}

func (uc *UserCommand) Execute(args []string) error {
	log.Printf("[INFO] user %s creating command is started", uc.Name)
	b, err := engine.NewBoltDB(bolt.Options{Timeout: 30 * time.Second}, "/tmp/shopaholic.db")
	if err != nil {
		return err
	}

	s := &service.DataStore{Interface: b}

	user := store.User{
		Name:    uc.Name,
		Balance: uc.Balance,
	}

	userID, err := s.Register(user)
	if err != nil {
		return err
	}

	log.Printf("[INFO] user %s was created with ID: %s", uc.Name, userID)
	return nil
}
