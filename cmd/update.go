package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/output"
	"github.com/dpcat237/go-dsu/internal/updater"
)

var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update dependencies",
		Long:  `Add missing and remove unused modules. Update modules`,
		Run: func(cmd *cobra.Command, args []string) {
			update(cmd)
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Bool("dev", false, "Development mode")
}

func update(cmd *cobra.Command) {
	md := output.ModeProd
	if cmd.Flag("dev").Value.String() == "true" {
		md = output.ModeDev
	}
	exc, out := executor.Init()
	if out.HasError() {
		fmt.Println(out.ToString(md))
		os.Exit(1)
	}

	upd := updater.Init(exc)
	out = upd.UpdateDependencies()
	fmt.Println(out.ToString(md))
}
