package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dpcat237/go-dsu/internal/cleaner"
	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/output"
)

var (
	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Clean modules",
		Long:  `Adds missing and remove unused modules`,
		Run: func(cmd *cobra.Command, args []string) {
			clean(cmd)
		},
	}
)

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().Bool("dev", false, "Development mode")
}

func clean(cmd *cobra.Command) {
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

	exc, out := executor.Init(lgr)
	if out.HasError() {
		fmt.Println(out.ToString(mod))
		os.Exit(1)
	}

	cln := cleaner.Init(exc)
	out = cln.Clean()
	fmt.Println(out.ToString(mod))
}
