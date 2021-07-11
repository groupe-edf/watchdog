package cmd

import (
	"context"
	"os"
	"runtime/debug"

	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/server"
	"github.com/spf13/cobra"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Print the version number of watchdog",
		Long:  `All software has versions. This is watchdog's`,
		Run: func(cmd *cobra.Command, args []string) {
			debug.SetTraceback("crash")
			logger := logging.New(logging.Options{
				LogsLevel:        "INFO",
				LogsOutput:       os.Stdout,
				LogsReportCaller: false,
			})
			server := server.New(logger)
			listener, err := server.Listener()
			if err != nil {
				logger.Fatal(err)
			}
			ctx := context.Background()
			server.Run(ctx, listener)
		},
	}
)

func init() {
	rootCommand.AddCommand(serveCmd)
}
