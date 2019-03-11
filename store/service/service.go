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

func (s *DataStore) StoreTransaction(transaction store.Transaction) (transactionID string, err error) {
	if utils.ToMoney(transaction.Amount).IsPositive() {
		return s.income(transaction)
	} else {
		return s.expense(transaction)
	}
}

func (s *DataStore) expense(transaction store.Transaction) (transactionID string, err error) {
	transaction.Type = store.Expense
	if transaction, err = s.prepareNewTransaction(transaction); err != nil {
		return "", errors.Wrap(err, "failed to prepare transaction")
	}
	if utils.ToMoney(transaction.Amount).IsPositive() {
		return "", errors.Wrap(err, "expense amount have to be negative")
	}
	if transaction, err = s.adjustBalanceByTransaction(transaction); err != nil {
		return "", errors.Wrap(err, "failed to adjust balance")
	}

	return s.Interface.Create(transaction)
}

func (s *DataStore) income(transaction store.Transaction) (transactionID string, err error) {
	transaction.Type = store.Income
	if transaction, err = s.prepareNewTransaction(transaction); err != nil {
		return "", errors.Wrap(err, "failed to prepare transaction")
	}
	if utils.ToMoney(transaction.Amount).IsNegative() {
		return "", errors.Wrap(err, "expense amount have to be positive")
	}
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

	transaction.User.Balance = transaction.BalanceNow
	if s.UpdateUser(transaction.User) != nil {
		return transaction, err
	}

	return transaction, nil
}

func (s *DataStore) prepareNewTransaction(transaction store.Transaction) (store.Transaction, error) {
	if transaction.ID == "" {
		transaction.ID = uuid.New().String()
	}

	if transaction.CreatedAt.IsZero() {
		transaction.CreatedAt = time.Now()
	}

	if transaction.Category.ID == "" {
		transaction.Category, _ = s.Interface.DefaultCategory(transaction.Type)
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

	userID, err = s.Interface.Register(user)
	if utils.ToMoney(user.Balance).IsPositive() {
		transaction := store.Transaction{
			User:   user,
			Type:   store.Income,
			Amount: user.Balance,
		}

		transaction.User.Balance.Amount = 0
		if _, err = s.StoreTransaction(transaction); err != nil {
			return "", err
		}
	}

	return userID, err
}

func (s *DataStore) UserList(number int) ([]store.User, error) {
	return s.Interface.Users(number)
}

func (s *DataStore) DefineType(title string) (store.Type, string) {
	firstChar := string(title[0])
	if firstChar == "+" {
		return store.Income, title[1:]
	}
	if firstChar == "-" {
		return store.Expense, title[1:]
	}

	return store.Expense, title
}
