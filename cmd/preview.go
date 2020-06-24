package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dpcat237/go-dsu/internal/output"
	"github.com/dpcat237/go-dsu/internal/service"
)

var (
	previewCmd = &cobra.Command{
		Use:   "preview",
		Short: "Preview updates",
		Long:  `Preview available updates of direct modules with changes`,
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

	fmt.Println("Analyzing dependencies...")
	if out := checkPrerequisites(); out.HasError() {
		fmt.Println(out.ToString(mod))
		return
	}

	prw, initOut := service.InitPreviewer(mod)
	if initOut.HasError() {
		fmt.Println(initOut.ToString(mod))
		os.Exit(1)
	}
	out := prw.Preview(cmd.Flag("path").Value.String())
	fmt.Println(out.ToString(mod))
}
