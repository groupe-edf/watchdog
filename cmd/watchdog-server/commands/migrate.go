package commands

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/server/store"
	"github.com/groupe-edf/watchdog/internal/server/store/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// Register file loader
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run migration scripts",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	migrateDropCmd = &cobra.Command{
		Use:   "drop",
		Short: "Run down migrations",
		Run: func(_ *cobra.Command, _ []string) {
			migrater := NewMigrater()
			defer func() {
				if _, err := migrater.Close(); err != nil {
					log.Fatal(err)
				}
			}()
			if err := migrater.Drop(); err != nil {
				if err != migrate.ErrNoChange {
					log.Fatal(err)
				}
			}
		},
	}
	migrateDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Run down migrations",
		Run: func(_ *cobra.Command, _ []string) {
			migrater := NewMigrater()
			defer func() {
				if _, err := migrater.Close(); err != nil {
					log.Fatal(err)
				}
			}()
			if err := migrater.Steps(-1); err != nil {
				if err != migrate.ErrNoChange {
					log.Fatal(err)
				}
			}
		},
	}
	migrateUpCmd = &cobra.Command{
		Use:   "up",
		Short: "Run up migrations",
		Run: func(_ *cobra.Command, _ []string) {
			migrater := NewMigrater()
			defer func() {
				if _, err := migrater.Close(); err != nil {
					log.Fatal(err)
				}
			}()
			if err := migrater.Up(); err != nil {
				if err != migrate.ErrNoChange {
					log.Fatal(err)
				}
			}
		},
	}
)

// NewMigrater return new migrater
func NewMigrater() *migrate.Migrate {
	options, err := config.NewOptions(viper.GetViper())
	if err != nil {
		log.Panic(err)
	}
	database, err := store.CreateLoopConnection(options.Database.Postgres)
	if err != nil {
		log.Panic(err)
	}
	migrater, err := postgres.NewMigrater(database)
	if err != nil {
		log.Panic(err)
	}
	return migrater
}

func init() {
	migrateCmd.AddCommand(migrateDropCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	rootCommand.AddCommand(migrateCmd)
}
