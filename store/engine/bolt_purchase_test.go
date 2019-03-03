package engine

import (
	bolt "github.com/coreos/bbolt"
	"github.com/stretchr/testify/assert"
	"os"
	"shopaholic/store"
	"testing"
	"time"
)

var testDb = "test-purchase.db"

func TestBoltDB_CreateAndList(t *testing.T) {
	defer os.Remove(testDb)
	var b = prep(t)

	res, err := b.List(store.User{ID: "user1"})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, 2121, res[0].Amount)
	assert.Equal(t, "user1", res[0].User.ID)
	t.Log(res[0].ID)

	_, err = b.Create(store.Purchase{ID: res[0].ID, User: store.User{ID: "user1"}})
	assert.NotNil(t, err)
	assert.Equal(t, "key id-1 already in store", err.Error())

	_, err = b.List(store.User{ID: "user-not-found"})
	assert.EqualError(t, err, `no bucket user-not-found in store`)

	assert.NoError(t, b.Close())
}

func TestBoltDB_New(t *testing.T) {
	_, err := NewBoltDB(bolt.Options{}, BoltInstance{FileName: "/tmp/no-such-place/tmp.db"})
	assert.EqualError(t, err, "failed to make boltdb for /tmp/no-such-place/tmp.db: open /tmp/no-such-place/tmp.db: no such file or directory")
}

func prep(t *testing.T) *BoltDB {
	os.Remove(testDb)

	boltStore, err := NewBoltDB(bolt.Options{}, BoltInstance{FileName: testDb})
	assert.Nil(t, err)
	b := boltStore

	purchase := store.Purchase{
		ID:        "id-1",
		Amount:    2121,
		CreatedAt: time.Date(2017, 12, 20, 15, 18, 22, 0, time.Local),
		User:      store.User{ID: "user1", Name: "user name"},
	}
	_, err = b.Create(purchase)
	assert.Nil(t, err)

	purchase2 := store.Purchase{
		ID:        "id-2",
		Amount:    21212,
		CreatedAt: time.Date(2017, 12, 20, 15, 18, 22, 0, time.Local),
		User:      store.User{ID: "user2", Name: "second name"},
	}
	_, err = b.Create(purchase2)
	assert.Nil(t, err)

	return b
}
