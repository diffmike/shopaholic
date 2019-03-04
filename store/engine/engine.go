package engine

import "shopaholic/store"

type Interface interface {
	Purchase
	User
}

type Purchase interface {
	Create(comment store.Purchase) (purchaseID string, err error)   // create new purchase, avoid dups by id
	Get(user store.User, purchaseID string) (store.Purchase, error) // get purchase by id
	Put(user store.User, purchase store.Purchase) error             // update purchase, mutable parts only
	List(user store.User) ([]store.Purchase, error)                 // list purchases for user
	Close() error                                                   // close/stop engine
}

type User interface {
	Register(user store.User) (userID string, err error)
	Details(userID string) (store.User, error)
	Users(number int) ([]store.User, error)
}
