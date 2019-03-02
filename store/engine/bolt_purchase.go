package engine

import (
	"encoding/json"
	"github.com/hashicorp/go-multierror"
	"shopaholic/store"

	bolt "github.com/coreos/bbolt"
	log "github.com/go-pkgz/lgr"
	"github.com/pkg/errors"
)

type BoltDB struct {
	db *bolt.DB
}

const (
	purchasesBucketName = "purchases"
	usersBucketName     = "users"
)

type BoltInstance struct {
	FileName string // full path to boltdb
}

func NewBoltDB(options bolt.Options, instance BoltInstance) (*BoltDB, error) {
	log.Printf("[INFO] bolt store with options %+v", options)

	result := BoltDB{}
	db, err := bolt.Open(instance.FileName, 0600, &options) // bolt.Options{Timeout: 30 * time.Second}
	result.db = db

	if err != nil {
		return nil, errors.Wrapf(err, "failed to make boltdb for %s", instance.FileName)
	}

	// make top-level buckets
	topBuckets := []string{purchasesBucketName, usersBucketName}
	err = db.Update(func(tx *bolt.Tx) error {
		for _, bktName := range topBuckets {
			if _, e := tx.CreateBucketIfNotExists([]byte(bktName)); e != nil {
				return errors.Wrapf(e, "failed to create top level bucket %s", bktName)
			}
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to create top level bucket)")
	}

	return &result, nil
}

// Create new purchase to store. Adds to purchases bucket
func (b *BoltDB) Create(purchase store.Purchase) (purchaseID string, err error) {

	err = b.db.Update(func(tx *bolt.Tx) error {

		postBkt, e := b.makePurchaseBucket(tx, purchase.User.ID)
		if e != nil {
			return e
		}

		// check if key already in store, reject doubles
		if postBkt.Get([]byte(purchase.ID)) != nil {
			return errors.Errorf("key %s already in store", purchase.ID)
		}

		// serialize comment to json []byte for bolt and save
		if e = b.save(postBkt, []byte(purchase.ID), purchase); e != nil {
			return errors.Wrapf(e, "failed to put key %s", purchase.ID)
		}

		return nil
	})

	return purchase.ID, err
}

// Find returns all purchases for user
func (b *BoltDB) Find(user store.User) (purchases []store.Purchase, err error) {
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

// Close boltdb store
func (b *BoltDB) Close() error {
	errs := new(multierror.Error)
	err := errors.Wrapf(b.db.Close(), "can't close site")
	errs = multierror.Append(errs, err)
	return errs.ErrorOrNil()
}

// getPurchaseBucket return bucket with all comments for postURL
func (b *BoltDB) getPurchaseBucket(tx *bolt.Tx, userID string) (*bolt.Bucket, error) {
	postsBkt := tx.Bucket([]byte(purchasesBucketName))
	if postsBkt == nil {
		return nil, errors.Errorf("no bucket %s", purchasesBucketName)
	}
	res := postsBkt.Bucket([]byte(userID))
	if res == nil {
		return nil, errors.Errorf("no bucket %s in store", userID)
	}
	return res, nil
}

// makePurchaseBucket create new bucket for userID as a key. This bucket holds all purchases for the user.
func (b *BoltDB) makePurchaseBucket(tx *bolt.Tx, userID string) (*bolt.Bucket, error) {
	postsBkt := tx.Bucket([]byte(purchasesBucketName))
	if postsBkt == nil {
		return nil, errors.Errorf("no bucket %s", purchasesBucketName)
	}
	res, err := postsBkt.CreateBucketIfNotExists([]byte(userID))
	if err != nil {
		return nil, errors.Wrapf(err, "no bucket %s in store", userID)
	}
	return res, nil
}

func (b *BoltDB) getUserBucket(tx *bolt.Tx, userID string) (*bolt.Bucket, error) {
	usersBkt := tx.Bucket([]byte(usersBucketName))
	userIDBkt, e := usersBkt.CreateBucketIfNotExists([]byte(userID)) // get bucket for userID
	if e != nil {
		return nil, errors.Wrapf(e, "can't get bucket %s", userID)
	}
	return userIDBkt, nil
}

// save marshaled value to key for bucket. Should run in update tx
func (b *BoltDB) save(bkt *bolt.Bucket, key []byte, value interface{}) (err error) {
	if value == nil {
		return errors.Errorf("can't save nil value for %s", key)
	}
	jdata, jerr := json.Marshal(value)
	if jerr != nil {
		return errors.Wrap(jerr, "can't marshal comment")
	}
	if err = bkt.Put(key, jdata); err != nil {
		return errors.Wrapf(err, "failed to save key %s", key)
	}
	return nil
}

func (b *BoltDB) load(bkt *bolt.Bucket, key []byte, res interface{}) error {
	value := bkt.Get(key)
	if value == nil {
		return errors.Errorf("no value for %s", key)
	}

	if err := json.Unmarshal(value, &res); err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}
	return nil
}
