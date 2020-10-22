package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/groupe-edf/watchdog/internal/version"
	"github.com/spf13/cobra"
)

const (
	defaultVersion = "develop"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of watchdog",
	Long:  `All software has versions. This is watchdog's`,
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo := version.GetBuildInfo()
		if bi, ok := debug.ReadBuildInfo(); ok && buildInfo.Version == defaultVersion {
			buildInfo.Version = bi.Main.Version
		}
		fmt.Println(string(buildInfo.ToJSON()))
	},
}

func init() {
	analyzeCommand.AddCommand(versionCmd)
}
