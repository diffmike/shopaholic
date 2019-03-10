package engine

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"shopaholic/store"
)

type RedisDB struct {
	db *redis.Client
}

const (
	transactionsRedisTable = "transactions"
	usersRedisTable        = "users"
	categoriesRedisTable   = "categories"
)

func NewRedisClient(db int) (*RedisDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       db, // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	r := RedisDB{db: client}
	if err = r.storeDefaultCategories(); err != nil {
		return nil, errors.Wrap(err, "failed to store default categories")
	}

	return &r, err
}

// Close redis store
func (r *RedisDB) Disconnect() error {
	errs := new(multierror.Error)
	err := errors.Wrapf(r.db.Close(), "can't close connection")
	errs = multierror.Append(errs, err)
	return errs.ErrorOrNil()
}

// save marshaled value to key into the table
func (r *RedisDB) save(table string, key string, value interface{}) (err error) {
	if value == nil {
		return errors.Errorf("can't save nil value for %s", key)
	}

	jdata, jerr := json.Marshal(value)
	if jerr != nil {
		return errors.Wrap(jerr, "can't marshal data")
	}
	if result := r.db.HSet(table, key, jdata); result.Err() != nil {
		return errors.Wrapf(result.Err(), "failed to save key %s", key)
	}

	return nil
}

func (r *RedisDB) load(table string, key string, res interface{}) error {
	value := r.db.HGet(table, key)
	if value.Err() != nil {
		return errors.Errorf("no value for %s", key)
	}

	bytes, err := value.Bytes()
	if err != nil {
		return errors.Errorf("failed to transform to bytes %s", err.Error())
	}

	if err := json.Unmarshal(bytes, &res); err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}
	return nil
}

func (r *RedisDB) storeDefaultCategories() error {
	s := Category(r)
	if num, _ := s.CountCategories(); num > 0 {
		return nil
	}

	category := store.Category{
		Title:     "Products",
		IsDefault: true,
		Type:      store.Expense,
	}
	if _, err := s.StoreCategory(category); err != nil {
		return errors.Wrapf(err, "failed to save default category")
	}

	category = store.Category{
		Title:     "Salary",
		IsDefault: true,
		Type:      store.Income,
	}
	if _, err := s.StoreCategory(category); err != nil {
		return errors.Wrapf(err, "failed to save default category")
	}

	return nil
}

func transactionsTable(userID string) string {
	return transactionsRedisTable + ":" + userID
}
