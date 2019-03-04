package service

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
	"shopaholic/store"
	"shopaholic/store/engine"
	"time"
)

type DataStore struct {
	engine.Interface
}

func (s *DataStore) Purchase(purchase store.Purchase) (purchaseID string, err error) {
	if purchase, err = s.prepareNewPurchase(purchase); err != nil {
		return "", errors.Wrap(err, "failed to prepare comment")
	}

	return s.Interface.Create(purchase)
}

func (s *DataStore) prepareNewPurchase(purchase store.Purchase) (store.Purchase, error) {
	if purchase.ID == "" {
		purchase.ID = uuid.New().String()
	}
	if purchase.CreatedAt.IsZero() {
		purchase.CreatedAt = time.Now()
	}

	return purchase, nil
}

func (s *DataStore) Register(user store.User) (userID string, err error) {
	log.Printf("[INFO] storing user %s in the service", user.Name)
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}

	log.Printf("[INFO] storing user %+v", user)
	return s.Interface.Register(user)
}

func (s *DataStore) UserList(number int) ([]store.User, error) {
	return s.Interface.Users(number)
}
