package commands

import "github.com/spf13/cobra"

var (
	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "create backup archive",
		Long:  ``,
		Run: func(cmd *cobra.Command, _ []string) {

		},
	}
)

func init() {
	rootCommand.AddCommand(backupCmd)
}
