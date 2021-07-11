package bolt

import (
	"encoding/json"

	"github.com/groupe-edf/watchdog/internal/models"
	"github.com/groupe-edf/watchdog/internal/server/query"
	bolt "go.etcd.io/bbolt"
)

const (
	SettingsBucket = "settings"
)

func (store *BoltStore) GetSettings(q *query.Query) ([]*models.Setting, error) {
	var settings = make([]*models.Setting, 0)
	err := store.database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(SettingsBucket))
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			var setting *models.Setting
			err := json.Unmarshal(value, &setting)
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}
			settings = append(settings, setting)
		}
		return nil
	})
	return settings, err
}
