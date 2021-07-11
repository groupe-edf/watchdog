package commands

import (
	"fmt"
	"runtime/debug"

	"github.com/groupe-edf/watchdog/internal/version"
	"github.com/spf13/cobra"
)

const (
	defaultVersion = "develop"
)

var (
	short      bool
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of watchdog",
		Long:  `All software has versions. This is watchdog's`,
		Run: func(_ *cobra.Command, _ []string) {
			buildInfo := version.GetBuildInfo()
			if bi, ok := debug.ReadBuildInfo(); ok && buildInfo.Version == defaultVersion {
				buildInfo.Version = bi.Main.Version
			}
			if short {
				fmt.Println(buildInfo.Version)
			} else {
				fmt.Println(string(buildInfo.ToJSON()))
			}
		},
	}
)

func init() {
	rootCommand.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&short, "short", false, "print just the version number.")
}
