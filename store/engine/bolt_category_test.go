package engine

import (
	"github.com/stretchr/testify/assert"
	"os"
	"shopaholic/store"
	"testing"
)

func TestBoltDB_DefaultCategory(t *testing.T) {
	defer os.Remove(testDb)
	b := prepBolt(t)

	category, err := b.DefaultCategory(store.Income)
	assert.Equal(t, true, category.IsDefault)
	assert.Equal(t, store.Income, category.Type)
	assert.Nil(t, err)

	category, err = b.DefaultCategory(store.Expense)
	assert.Equal(t, true, category.IsDefault)
	assert.Equal(t, store.Expense, category.Type)
	assert.Nil(t, err)
}

func TestBoltDB_CountCategories(t *testing.T) {
	defer os.Remove(testDb)
	b := prepBolt(t)

	count, err := b.CountCategories()
	assert.Equal(t, 2, count, "2 default categories have to be created on init")
	assert.Nil(t, err)
}
