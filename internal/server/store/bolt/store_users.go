package bolt

import (
	"encoding/json"
	"strings"

	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/server/query"
	bolt "go.etcd.io/bbolt"
)

const (
	UsersBucket = "users"
)

func (store *BoltStore) DeleteBucket(bucket string) error {
	err := store.database.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucket))
	})
	return err
}

func (store *BoltStore) DeleteUsers(q *query.Query) error {
	return store.DeleteBucket(UsersBucket)
}

func (store *BoltStore) FindUserByEmail(email string) (*models.User, error) {
	var user *models.User
	err := store.database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UsersBucket))
		if bucket == nil {
			return models.ErrUserNotFound
		}
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			err := json.Unmarshal(value, &user)
			if err != nil {
				return err
			}
			if strings.EqualFold(user.Email, email) {
				break
			}
		}
		if user == nil {
			return models.ErrUserNotFound
		}
		return nil
	})
	return user, err
}

func (store *BoltStore) SaveUser(user *models.User) (*models.User, error) {
	err := store.database.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(UsersBucket))
		if bucket == nil {
			return models.ErrUserNotFound
		}
		encoded, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(user.ID.String()), encoded)
	})
	return user, err
}

func (store *BoltStore) SaveOrUpdateUser(data *models.User) (*models.User, error) {
	user, err := store.FindUserByEmail(data.Email)
	if err != nil {
		return nil, err
	}
	user.LastLogin = data.LastLogin
	user.Username = data.Username
	return store.SaveUser(user)
}
