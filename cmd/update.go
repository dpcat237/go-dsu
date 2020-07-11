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
		Short: "Updater dependencies",
		Long:  `Update dependencies`,
		Run: func(cmd *cobra.Command, args []string) {
			update(cmd)
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Bool("dev", false, "Development mode")
	updateCmd.Flags().BoolP("all", "", false, "Updater all without verifications")
	updateCmd.Flags().BoolP("prompt", "p", false, "Confirm in prompt updates with changes")
	updateCmd.Flags().BoolP("select", "s", false, "Select modules to update")
	updateCmd.Flags().BoolP("tests", "t", false, "Run local tests after updating each module and rollback in case of errors")
	updateCmd.Flags().BoolP("verbose", "v", false, "Print output")
}

func update(cmd *cobra.Command) {
	mod := output.ModeProd
	if cmd.Flag("dev").Changed {
		mod = output.ModeDev
	}

	updOpt := updater.UpdateOptions{
		IsAll:     cmd.Flag("all").Changed,
		IsPrompt:  cmd.Flag("prompt").Changed,
		IsSelect:  cmd.Flag("select").Changed,
		IsTests:   cmd.Flag("tests").Changed,
		IsVerbose: cmd.Flag("verbose").Changed,
	}

	fmt.Println("Analyzing prerequisites...")
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
