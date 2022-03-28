package commands

import (
	"os"

	"github.com/gookit/color"
	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/git"
	"github.com/groupe-edf/watchdog/internal/server"
	"github.com/groupe-edf/watchdog/internal/server/broadcast"
	"github.com/groupe-edf/watchdog/internal/server/store"
	"github.com/groupe-edf/watchdog/pkg/authentication"
	"github.com/groupe-edf/watchdog/pkg/container"
	"github.com/groupe-edf/watchdog/pkg/event"
	"github.com/groupe-edf/watchdog/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "print the version number of watchdog",
		Long:  `all software has versions. This is watchdog's`,
		Run: func(cmd *cobra.Command, _ []string) {
			options, err := config.NewOptions(viper.GetViper())
			if err != nil {
				color.Red.Printf("unable to decode into config struct, %v", err)
				os.Exit(0)
			}
			di := container.GetContainer()
			di.Set(config.ServiceName, func(c container.Container) container.Service {
				return options
			})
			di.Provide(&event.ServiceProvider{})
			di.Provide(&logging.ServiceProvider{Options: options.Logs})
			di.Provide(&git.ServiceProvider{Options: options})
			di.Provide(&store.ServiceProvider{Options: options.Database})
			di.Provide(&authentication.ServiceProvider{})
			di.Provide(&broadcast.ServiceProvider{})
			// starting server
			logger := di.Get(logging.ServiceName).(logging.Interface)
			server := server.New(cmd.Context(), logger)
			listener, err := server.Listener()
			if err != nil {
				logger.Fatal(err)
			}
			server.RegisterEvents()
			server.Start(listener)
		},
	}
)

func init() {
	rootCommand.AddCommand(serveCmd)
}
