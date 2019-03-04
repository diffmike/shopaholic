package engine

import (
	bolt "github.com/coreos/bbolt"
	"github.com/pkg/errors"
	"log"
	"shopaholic/store"
)

func (b *BoltDB) Register(user store.User) (userID string, err error) {
	log.Printf("[INFO] storing user %s into the bucket", user.Name)
	err = b.db.Update(func(tx *bolt.Tx) error {
		usersBkt := tx.Bucket([]byte(usersBucketName))

		// check if key already in store, reject doubles
		if usersBkt.Get([]byte(user.ID)) != nil {
			return errors.Errorf("key %s already in store", user.ID)
		}

		// serialize user to json []byte for bolt and save
		if e := b.save(usersBkt, []byte(user.ID), user); e != nil {
			return errors.Wrapf(e, "failed to put key %s", user.ID)
		}

		return nil
	})

	return user.ID, err
}

func (b *BoltDB) Details(userID string) (user store.User, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		usersBkt := tx.Bucket([]byte(usersBucketName))
		return b.load(usersBkt, []byte(userID), &user)
	})
	return user, err
}
