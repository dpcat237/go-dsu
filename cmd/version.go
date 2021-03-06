package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "0.9.1"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version",
		Long:  `Version details`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("go-dsu %s\n", version)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
