package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/license"
	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
	"github.com/dpcat237/go-dsu/internal/previewer"
)

var (
	previewCmd = &cobra.Command{
		Use:   "preview",
		Short: "Preview updates",
		Long:  `Preview available updates of direct modules`,
		Run: func(cmd *cobra.Command, args []string) {
			preview(cmd)
		},
	}
)

func init() {
	rootCmd.AddCommand(previewCmd)
	previewCmd.Flags().Bool("dev", false, "Development mode")
	previewCmd.Flags().String("path", "", "Preview project from git path. Eg. github.com/spf13/cobra")
}

func preview(cmd *cobra.Command) {
	mod := output.ModeProd
	if cmd.Flag("dev").Value.String() == "true" {
		mod = output.ModeDev
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

	exc, excOut := executor.Init(lgr)
	if excOut.HasError() {
		fmt.Println(excOut.ToString(mod))
		os.Exit(1)
	}

	licHnd, licHndOut := license.InitHandler(lgr)
	if licHndOut.HasError() {
		fmt.Println(licHndOut.ToString(mod))
		os.Exit(1)
	}

	hnd := module.InitHandler(exc, lgr, licHnd)
	upd := previewer.Init(exc, lgr, hnd)
	out := upd.Preview(cmd.Flag("path").Value.String())
	fmt.Println(out.ToString(mod))
}
