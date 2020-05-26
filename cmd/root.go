package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-dsu",
	Short: "Go Dependencies Secure Updater",
	Long:  `Go DSU - provides tools to update Go dependencies with more control than default Go modules.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
