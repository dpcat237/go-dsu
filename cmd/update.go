package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dpcat237/go-dsu/internal/output"
	"github.com/dpcat237/go-dsu/internal/service"
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
	updateCmd.Flags().BoolP("prompt", "p", false, "Confirm in prompt updates with changes")
	updateCmd.Flags().BoolP("select", "s", false, "Select direct modules to update")
	updateCmd.Flags().BoolP("tests", "t", false, "Run local tests after updating each module and rollback in case of errors")
	updateCmd.Flags().BoolP("verbose", "v", false, "Print output")
}

func update(cmd *cobra.Command) {
	mod := output.ModeProd
	if cmd.Flag("dev").Changed {
		mod = output.ModeDev
	}

	updOpt := updater.UpdateOptions{
		IsIndirect: cmd.Flag("indirect").Changed,
		IsPrompt:   cmd.Flag("prompt").Changed,
		IsSelect:   cmd.Flag("select").Changed,
		IsTests:    cmd.Flag("tests").Changed,
		IsVerbose:  cmd.Flag("verbose").Changed,
	}

	fmt.Println("Analyzing dependencies...")
	if out := checkPrerequisites(); out.HasError() {
		fmt.Println(out.ToString(mod))
		return
	}

	upd, initOut := service.InitUpdater(mod)
	if initOut.HasError() {
		fmt.Println(initOut.ToString(mod))
		os.Exit(1)
	}
	out := upd.UpdateModules(updOpt)
	fmt.Println(out.ToString(mod))
}
