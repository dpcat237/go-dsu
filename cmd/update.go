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
	updateCmd.Flags().BoolP("select", "s", false, "Select direct modules to update")
	updateCmd.Flags().BoolP("test", "t", false, "Run local tests after updating each module and rollback in case of errors")
	updateCmd.Flags().BoolP("verbose", "v", false, "Print output")
}

func update(cmd *cobra.Command) {
	md := output.ModeProd
	var ind, scl, tst, vrb bool

	if cmd.Flag("dev").Value.String() == "true" {
		md = output.ModeDev
	}
	if cmd.Flag("indirect").Value.String() == "true" {
		ind = true
	}
	if cmd.Flag("select").Value.String() == "true" {
		scl = true
	}
	if cmd.Flag("test").Value.String() == "true" {
		tst = true
	}
	if cmd.Flag("verbose").Value.String() == "true" {
		vrb = true
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
	out = upd.UpdateModules(ind, scl, tst, vrb)
	fmt.Println(out.ToString(md))
}
