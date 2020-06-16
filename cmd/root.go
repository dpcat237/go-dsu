package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dpcat237/go-dsu/internal/httputil"
	"github.com/dpcat237/go-dsu/internal/output"
)

var rootCmd = &cobra.Command{
	Use:   "go-dsu",
	Short: "Go Dependencies Secure Updater",
	Long:  `Go DSU - provides tools to update Go dependencies with more control than default Go modules.`,
}

//Execute builds CLI commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkPrerequisites() output.Output {
	out := output.Create("cmd.checkPrerequisites")

	if !httputil.IsConnection() {
		return out.WithErrorString("Check your Internet connection")
	}
	return out
}
