package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// provided by https://goreleaser.com/
	commit  = "none"
	date    = "unknown"
	version = "dev"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version",
		Long:  `Version details`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("go-dsu %s built from commit %s on %s\n", version, commit, date)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
