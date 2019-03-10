package engine

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"shopaholic/store"
)

func (r *RedisDB) Register(user store.User) (userID string, err error) {
	log.Printf("[INFO] storing user %s into the table", user.Name)

	if existed := r.db.HGet(usersRedisTable, user.ID); existed.Val() != "" {
		return "", errors.Errorf("key %s already in store", user.ID)
	}

	if e := r.save(usersRedisTable, user.ID, user); e != nil {
		return "", errors.Wrapf(e, "failed to put key %s", user.ID)
	}

	return user.ID, err
}

func (r *RedisDB) UpdateUser(user store.User) error {
	return r.save(usersRedisTable, user.ID, user)
}

func (r *RedisDB) Details(userID string) (user store.User, err error) {
	err = r.load(usersRedisTable, userID, &user)
	return user, err
}

func (r *RedisDB) Users(number int) (users []store.User, err error) {
	users = []store.User{}

	data := r.db.HGetAll(usersRedisTable)
	if data.Err() != nil {
		return users, data.Err()
	}

	for _, v := range data.Val() {
		user := store.User{}
		if e := json.Unmarshal([]byte(v), &user); e != nil {
			return users, errors.Wrap(e, "failed to unmarshal")
		}
		users = append(users, user)
	}

	return users, err
}
