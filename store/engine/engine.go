package engine

import (
	"shopaholic/store"
)

type Interface interface {
	Transaction
	User
	Category
}

type Connection interface {
	Disconnect() error // close/stop engine
}

type Transaction interface {
	Create(transaction store.Transaction) (transactionID string, err error) // create new transaction, avoid dups by id
	Get(user store.User, transactionID string) (store.Transaction, error)   // get transaction by id
	Put(user store.User, transaction store.Transaction) error               // update transaction, mutable parts only
	List(user store.User) ([]store.Transaction, error)                      // list transactions for user
}

type User interface {
	Register(user store.User) (userID string, err error)
	Details(userID string) (store.User, error)
	Users(number int) ([]store.User, error)
	UpdateUser(user store.User) error
}

type Category interface {
	StoreCategory(category store.Category) (categoryID string, err error)
	CountCategories() (number int, err error)
	DefaultCategory(t store.Type) (category store.Category, err error)
	FindCategoryByTitle(title string) (category store.Category, err error)
	FindCategoryByTitleAndUserID(title string, userID string) (category store.Category, err error)
}
