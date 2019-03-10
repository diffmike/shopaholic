package engine

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"shopaholic/store"
	"strings"
	"time"
)

func (r *RedisDB) StoreCategory(category store.Category) (categoryID string, err error) {
	if category.ID == "" {
		category.ID = uuid.New().String()
	}

	if existed := r.db.HGet(categoriesRedisTable, category.ID); existed.Val() != "" {
		return "", errors.Errorf("key %s already in store", category.ID)
	}

	category.CreatedAt = time.Now()

	if e := r.save(categoriesRedisTable, category.ID, category); e != nil {
		return "", errors.Wrapf(e, "failed to put key %s", category.ID)
	}

	return category.ID, err
}

func (r *RedisDB) CountCategories() (number int, err error) {
	result := r.db.HLen(categoriesRedisTable)

	return int(result.Val()), result.Err()
}

func (r *RedisDB) DefaultCategory(t store.Type) (category store.Category, err error) {
	keys, err := r.db.HKeys(categoriesRedisTable).Result()
	if err != nil {
		return category, err
	}

	tmpCat := store.Category{}
	for _, k := range keys {
		err = r.load(categoriesRedisTable, k, &tmpCat)
		if tmpCat.IsDefault && tmpCat.Type == t {
			category = tmpCat
			break
		}
	}

	return category, err
}

func (r *RedisDB) FindCategoryByTitle(title string) (category store.Category, err error) {
	keys, err := r.db.HKeys(categoriesRedisTable).Result()
	if err != nil {
		return category, err
	}

	title = strings.ToLower(title)
	tmpCat := store.Category{}
	for _, k := range keys {
		err = r.load(categoriesRedisTable, k, &tmpCat)
		if strings.ToLower(tmpCat.Title) == title {
			category = tmpCat
			break
		}
	}

	return category, err
}
