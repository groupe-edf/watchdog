package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/server/container"
	"github.com/groupe-edf/watchdog/internal/server/store/bolt"
	"github.com/groupe-edf/watchdog/internal/server/store/postgres"
)

const (
	// ConnectionTimeout connection timeout
	ConnectionTimeout = 30 * time.Second
	ServiceName       = "store"
)

type ServiceProvider struct {
	Options *config.Database
}

func (service *ServiceProvider) Register(di container.Container) {
	di.Set(ServiceName, func(_ container.Container) container.Service {
		switch service.Options.Driver {
		case config.PostgresDriver:
			database, err := CreateLoopConnection(service.Options.Postgres)
			if err != nil {
				panic(err)
			}
			store, err := postgres.NewPostgresStore(database)
			if err != nil {
				panic(err)
			}
			return store
		case config.BoldDriver:
			store, err := bolt.NewBoltStore(service.Options.Bolt)
			if err != nil {
				panic(err)
			}
			return store
		}
		return nil
	})
}

func CreateLoopConnection(settings config.PostgresOptions) (*sql.DB, error) {
	var err error
	dataSourceName := buildConnectionString(settings)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(ConnectionTimeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, err
		case <-ticker.C:
			db, _ := sql.Open("postgres", dataSourceName)
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
	}
}

func buildConnectionString(settings config.PostgresOptions) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=public",
		settings.Host,
		settings.Port,
		settings.Username,
		settings.Password,
		settings.Name)
}
