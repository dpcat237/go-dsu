package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dpcat237/go-dsu/internal/cleaner"
	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/license"
	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
	"github.com/dpcat237/go-dsu/internal/updater"
)

var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Updater modules",
		Long:  `Add missing and remove unused modules. Updater direct modules`,
		Run: func(cmd *cobra.Command, args []string) {
			update(cmd)
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Bool("dev", false, "Development mode")
	updateCmd.Flags().BoolP("indirect", "i", false, "Updater all direct and indirect modules")
	updateCmd.Flags().BoolP("select", "s", false, "Select direct modules to update")
	updateCmd.Flags().BoolP("test", "t", false, "Run local tests after updating each module and rollback in case of errors")
	updateCmd.Flags().BoolP("verbose", "v", false, "Print output")
}

func update(cmd *cobra.Command) {
	mod := output.ModeProd
	var ind, scl, tst, vrb bool

	if cmd.Flag("dev").Value.String() == "true" {
		mod = output.ModeDev
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
		fmt.Println(out.ToString(mod))
		return
	}

	lgr, lgrOut := logger.Init(mod)
	if lgrOut.HasError() {
		fmt.Println(lgrOut.ToString(mod))
		os.Exit(1)
	}

	exc, out := executor.Init(lgr)
	if out.HasError() {
		fmt.Println(out.ToString(mod))
		os.Exit(1)
	}

	licHnd, licHndOut := license.InitHandler(lgr)
	if licHndOut.HasError() {
		fmt.Println(licHndOut.ToString(mod))
		os.Exit(1)
	}

	cln := cleaner.Init(exc)
	hnd := module.InitHandler(exc, lgr, licHnd)
	upd := updater.Init(cln, exc, hnd)
	out = upd.UpdateModules(ind, scl, tst, vrb)
	fmt.Println(out.ToString(mod))
}
