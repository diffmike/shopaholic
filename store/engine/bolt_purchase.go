package engine

import (
	"encoding/json"
	bolt "github.com/coreos/bbolt"
	"github.com/pkg/errors"
	"shopaholic/store"
)

// Create new purchase to store. Adds to purchases bucket
func (b *BoltDB) Create(purchase store.Purchase) (purchaseID string, err error) {

	err = b.db.Update(func(tx *bolt.Tx) error {

		purchaseBkt, e := b.makePurchaseBucket(tx, purchase.User.ID)
		if e != nil {
			return e
		}

		// check if key already in store, reject doubles
		if purchaseBkt.Get([]byte(purchase.ID)) != nil {
			return errors.Errorf("key %s already in store", purchase.ID)
		}

		// serialize comment to json []byte for bolt and save
		if e = b.save(purchaseBkt, []byte(purchase.ID), purchase); e != nil {
			return errors.Wrapf(e, "failed to put key %s", purchase.ID)
		}

		return nil
	})

	return purchase.ID, err
}

// List of all purchases for user
func (b *BoltDB) List(user store.User) (purchases []store.Purchase, err error) {
	purchases = []store.Purchase{}

	err = b.db.View(func(tx *bolt.Tx) error {

		bucket, e := b.getPurchaseBucket(tx, user.ID)
		if e != nil {
			return e
		}

		return bucket.ForEach(func(k, v []byte) error {
			purchase := store.Purchase{}
			if e := json.Unmarshal(v, &purchase); e != nil {
				return errors.Wrap(e, "failed to unmarshal")
			}
			purchases = append(purchases, purchase)
			return nil
		})
	})

	return purchases, err
}

func (b *BoltDB) Get(user store.User, purchaseID string) (purchase store.Purchase, err error) {

	err = b.db.View(func(tx *bolt.Tx) error {
		bucket, e := b.getPurchaseBucket(tx, user.ID)
		if e != nil {
			return e
		}
		return b.load(bucket, []byte(purchaseID), &purchase)
	})
	return purchase, err
}

func (b *BoltDB) Put(user store.User, purchase store.Purchase) error {

	if curPurchase, err := b.Get(user, purchase.ID); err == nil {
		purchase.Amount = curPurchase.Amount
		purchase.CreatedAt = curPurchase.CreatedAt
		purchase.User = curPurchase.User
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		bucket, e := b.getPurchaseBucket(tx, user.ID)
		if e != nil {
			return e
		}
		return b.save(bucket, []byte(purchase.ID), purchase)
	})
}
