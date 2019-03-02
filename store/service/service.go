package service

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"shopaholic/store"
	"shopaholic/store/engine"
	"time"
)

type DataStore struct {
	engine.Interface
}

func (s *DataStore) Create(purchase store.Purchase) (purchaseID string, err error) {
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
