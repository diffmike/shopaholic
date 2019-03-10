package engine

import (
	"github.com/pkg/errors"
	"shopaholic/store"
)

// Create new transaction to store. Adds to transactions bucket
func (r *RedisDB) Create(transaction store.Transaction) (transactionID string, err error) {
	if existed := r.db.HGet(transactionsTable(transaction.User.ID), transaction.ID); existed.Val() != "" {
		return "", errors.Errorf("key %s already in store", transaction.ID)
	}

	if e := r.save(transactionsTable(transaction.User.ID), transaction.ID, transaction); e != nil {
		return "", errors.Wrapf(e, "failed to put key %s", transaction.ID)
	}

	return transaction.ID, err
}
