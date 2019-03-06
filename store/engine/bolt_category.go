package engine

import (
	bolt "github.com/coreos/bbolt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"shopaholic/store"
	"time"
)

func (b *BoltDB) StoreCategory(category store.Category) (categoryID string, err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		categoriesBkt := tx.Bucket([]byte(categoriesBucketName))

		if category.ID == "" {
			category.ID = uuid.New().String()
		}

		if categoriesBkt.Get([]byte(category.ID)) != nil {
			return errors.Errorf("key %s already in store", category.ID)
		}

		category.CreatedAt = time.Now()

		if e := b.save(categoriesBkt, []byte(category.ID), category); e != nil {
			return errors.Wrapf(e, "failed to put key %s", category.ID)
		}

		return nil
	})

	return category.ID, err
}

func (b *BoltDB) CountCategories() (number int, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		categoriesBkt := tx.Bucket([]byte(categoriesBucketName))
		number = categoriesBkt.Stats().KeyN
		return nil
	})
	return number, err
}

func (b *BoltDB) DefaultCategory(t store.Type) (category store.Category, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		categoriesBkt := tx.Bucket([]byte(categoriesBucketName))

		c := categoriesBkt.Cursor()
		tmpCat := store.Category{}
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			_ = b.load(categoriesBkt, k, &tmpCat)
			if tmpCat.IsDefault && tmpCat.Type == t {
				category = tmpCat
				break
			}
		}

		return nil
	})

	return category, err
}
