package engine

import (
	"encoding/json"
	bolt "github.com/coreos/bbolt"
	"github.com/pkg/errors"
	"shopaholic/store"
	"sort"
)

// Create new transaction to store. Adds to transactions bucket
func (b *BoltDB) Create(transaction store.Transaction) (transactionID string, err error) {

	err = b.db.Update(func(tx *bolt.Tx) error {

		transactionBkt, e := b.makePurchaseBucket(tx, transaction.User.ID)
		if e != nil {
			return e
		}

		// check if key already in store, reject doubles
		if transactionBkt.Get([]byte(transaction.ID)) != nil {
			return errors.Errorf("key %s already in store", transaction.ID)
		}

		// serialize comment to json []byte for bolt and save
		if e = b.save(transactionBkt, []byte(transaction.ID), transaction); e != nil {
			return errors.Wrapf(e, "failed to put key %s", transaction.ID)
		}

		return nil
	})

	return transaction.ID, err
}

// List of all transactions for user
func (b *BoltDB) List(user store.User) (transactions []store.Transaction, err error) {
	transactions = []store.Transaction{}

	err = b.db.View(func(tx *bolt.Tx) error {

		bucket, e := b.getPurchaseBucket(tx, user.ID)
		if e != nil {
			return e
		}

		return bucket.ForEach(func(k, v []byte) error {
			transaction := store.Transaction{}
			if e := json.Unmarshal(v, &transaction); e != nil {
				return errors.Wrap(e, "failed to unmarshal")
			}
			transactions = append(transactions, transaction)
			return nil
		})
	})

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt.After(transactions[j].CreatedAt)
	})

	return transactions, err
}

func (b *BoltDB) Get(user store.User, transactionID string) (transaction store.Transaction, err error) {

	err = b.db.View(func(tx *bolt.Tx) error {
		bucket, e := b.getPurchaseBucket(tx, user.ID)
		if e != nil {
			return e
		}
		return b.load(bucket, []byte(transactionID), &transaction)
	})
	return transaction, err
}

func (b *BoltDB) Put(user store.User, transaction store.Transaction) error {

	if curPurchase, err := b.Get(user, transaction.ID); err == nil {
		transaction.Amount = curPurchase.Amount
		transaction.CreatedAt = curPurchase.CreatedAt
		transaction.User = curPurchase.User
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		bucket, e := b.getPurchaseBucket(tx, user.ID)
		if e != nil {
			return e
		}
		return b.save(bucket, []byte(transaction.ID), transaction)
	})
}
