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
		Short: "Update modules",
		Long:  `Add missing and remove unused modules. Update direct modules`,
		Run: func(cmd *cobra.Command, args []string) {
			update(cmd)
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Bool("dev", false, "Development mode")
	updateCmd.Flags().BoolP("indirect", "i", false, "Update all direct and indirect modules")
}

func update(cmd *cobra.Command) {
	md := output.ModeProd
	var ind bool

	if cmd.Flag("dev").Value.String() == "true" {
		md = output.ModeDev
	}
	if cmd.Flag("indirect").Value.String() == "true" {
		ind = true
	}

	if out := checkPrerequisites(); out.HasError() {
		fmt.Println(out.ToString(md))
		return
	}

	exc, out := executor.Init()
	if out.HasError() {
		fmt.Println(out.ToString(md))
		os.Exit(1)
	}

	upd := updater.Init(exc)
	out = upd.UpdateDependencies(ind)
	fmt.Println(out.ToString(md))
}
