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
	md := output.ModeProd
	if cmd.Flag("dev").Value.String() == "true" {
		md = output.ModeDev
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
	out = upd.Clean()
	fmt.Println(out.ToString(md))
}
