package engine

import (
	"github.com/stretchr/testify/assert"
	"os"
	"shopaholic/store"
	"shopaholic/utils"
	"testing"
	"time"
)

func TestRedisDB_DefaultCategory(t *testing.T) {
	defer os.Remove(testDb)
	r := prepRedis(t)

	category, err := r.DefaultCategory(store.Income)
	assert.Equal(t, true, category.IsDefault)
	assert.Equal(t, store.Income, category.Type)
	assert.Nil(t, err)

	category, err = r.DefaultCategory(store.Expense)
	assert.Equal(t, true, category.IsDefault)
	assert.Equal(t, store.Expense, category.Type)
	assert.Nil(t, err)

	r.db.FlushAll()
}

func TestRedisDB_CountCategories(t *testing.T) {
	defer os.Remove(testDb)
	r := prepRedis(t)

	count, err := r.CountCategories()
	assert.Equal(t, 2, count, "2 default categories have to be created on init")
	assert.Nil(t, err)

	r.db.FlushAll()
}

func prepRedis(t *testing.T) *RedisDB {
	redisClient, err := NewRedisClient(15)
	assert.Nil(t, err)

	transaction := store.Transaction{
		ID:        "id-1",
		Amount:    utils.Money{2121, "usd"},
		CreatedAt: time.Date(2017, 12, 20, 15, 18, 22, 0, time.Local),
		User:      store.User{ID: "user1", Name: "user name"},
	}
	_, err = redisClient.Create(transaction)
	assert.Nil(t, err)

	transaction2 := store.Transaction{
		ID:        "id-2",
		Amount:    utils.Money{21212, "usd"},
		CreatedAt: time.Date(2017, 12, 20, 15, 18, 22, 0, time.Local),
		User:      store.User{ID: "user2", Name: "second name"},
	}
	_, err = redisClient.Create(transaction2)
	assert.Nil(t, err)

	return redisClient
}
