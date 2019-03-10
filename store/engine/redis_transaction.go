package engine

import (
	"encoding/json"
	"github.com/pkg/errors"
	"shopaholic/store"
	"sort"
)

func (r *RedisDB) Create(transaction store.Transaction) (transactionID string, err error) {
	if existed := r.db.HGet(transactionsTable(transaction.User.ID), transaction.ID); existed.Val() != "" {
		return "", errors.Errorf("key %s already in store", transaction.ID)
	}

	if e := r.save(transactionsTable(transaction.User.ID), transaction.ID, transaction); e != nil {
		return "", errors.Wrapf(e, "failed to put key %s", transaction.ID)
	}

	return transaction.ID, err
}

func (r *RedisDB) List(user store.User) (transactions []store.Transaction, err error) {
	transactions = []store.Transaction{}

	data := r.db.HGetAll(transactionsTable(user.ID))
	if data.Err() != nil {
		return transactions, data.Err()
	}

	for _, v := range data.Val() {
		transaction := store.Transaction{}
		if e := json.Unmarshal([]byte(v), &transaction); e != nil {
			return transactions, errors.Wrap(e, "failed to unmarshal")
		}
		transactions = append(transactions, transaction)
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt.After(transactions[j].CreatedAt)
	})

	return transactions, err
}

func (r *RedisDB) Get(user store.User, transactionID string) (transaction store.Transaction, err error) {
	err = r.load(transactionsTable(user.ID), transactionID, &transaction)
	return transaction, err
}

func (r *RedisDB) Put(user store.User, transaction store.Transaction) error {
	if curPurchase, err := r.Get(user, transaction.ID); err == nil {
		transaction.Amount = curPurchase.Amount
		transaction.CreatedAt = curPurchase.CreatedAt
		transaction.User = curPurchase.User
	}

	return r.save(transactionsTable(user.ID), transaction.ID, transaction)
}
