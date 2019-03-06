package service

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
	"shopaholic/store"
	"shopaholic/store/engine"
	"shopaholic/utils"
	"time"
)

type DataStore struct {
	engine.Interface
}

func (s *DataStore) Expense(transaction store.Transaction) (transactionID string, err error) {
	if transaction, err = s.prepareNewTransaction(transaction); err != nil {
		return "", errors.Wrap(err, "failed to prepare transaction")
	}
	if utils.ToMoney(transaction.Amount).IsPositive() {
		return "", errors.Wrap(err, "expense amount have to be negative")
	}
	transaction.Type = store.Expense
	if transaction, err = s.adjustBalanceByTransaction(transaction); err != nil {
		return "", errors.Wrap(err, "failed to adjust balance")
	}

	return s.Interface.Create(transaction)
}

func (s *DataStore) Income(transaction store.Transaction) (transactionID string, err error) {
	if transaction, err = s.prepareNewTransaction(transaction); err != nil {
		return "", errors.Wrap(err, "failed to prepare transaction")
	}
	if utils.ToMoney(transaction.Amount).IsNegative() {
		return "", errors.Wrap(err, "expense amount have to be positive")
	}
	transaction.Type = store.Income
	if transaction, err = s.adjustBalanceByTransaction(transaction); err != nil {
		return "", errors.Wrap(err, "failed to adjust balance")
	}

	return s.Interface.Create(transaction)
}

func (s *DataStore) adjustBalanceByTransaction(transaction store.Transaction) (store.Transaction, error) {
	transaction.BalanceWas = transaction.User.Balance
	balanceWas, err := utils.ToMoney(transaction.User.Balance).Add(utils.ToMoney(transaction.Amount))
	if err != nil {
		return transaction, err
	}
	transaction.BalanceNow = utils.FromMoney(balanceWas)

	return transaction, nil
}

func (s *DataStore) prepareNewTransaction(transaction store.Transaction) (store.Transaction, error) {
	if transaction.ID == "" {
		transaction.ID = uuid.New().String()
	}
	if transaction.CreatedAt.IsZero() {
		transaction.CreatedAt = time.Now()
	}

	return transaction, nil
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
