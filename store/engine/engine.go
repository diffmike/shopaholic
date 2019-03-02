package engine

import "shopaholic/store"

type Interface interface {
	Purchase
}

type Purchase interface {
	Create(comment store.Purchase) (purchaseID string, err error) // create new purchase, avoid dups by id
	Get(purchaseID string) (store.Purchase, error)                // get purchase by id
	Put(purchase store.Purchase) error                            // update purchase, mutable parts only
	Find(user store.User, sort string) ([]store.Purchase, error)  // find purchases for user
	Count(locator store.User) (int, error)                        // number of purchases for the user
	Close() error                                                 // close/stop engine
}
