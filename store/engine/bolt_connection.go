package engine

import (
	"encoding/json"
	bolt "github.com/coreos/bbolt"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"log"
)

type BoltDB struct {
	db *bolt.DB
}

const (
	purchasesBucketName = "purchases"
	usersBucketName     = "users"
)

func NewBoltDB(options bolt.Options, filename string) (*BoltDB, error) {
	log.Printf("[INFO] bolt store with options %+v", options)

	result := BoltDB{}
	db, err := bolt.Open(filename, 0600, &options) // bolt.Options{Timeout: 30 * time.Second}
	result.db = db

	if err != nil {
		return nil, errors.Wrapf(err, "failed to make boltdb for %s", filename)
	}

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

// getPurchaseBucket return bucket with all purchases for the user
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

// Close boltdb store
func (b *BoltDB) Disconnect() error {
	errs := new(multierror.Error)
	err := errors.Wrapf(b.db.Close(), "can't close connection")
	errs = multierror.Append(errs, err)
	return errs.ErrorOrNil()
}

// save marshaled value to key for bucket. Should run in update tx
func (b *BoltDB) save(bkt *bolt.Bucket, key []byte, value interface{}) (err error) {
	if value == nil {
		return errors.Errorf("can't save nil value for %s", key)
	}

	jdata, jerr := json.Marshal(value)
	if jerr != nil {
		return errors.Wrap(jerr, "can't marshal data")
	}
	log.Printf("[INFO] storing from %+v to %s", value, jdata)
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
