package cmd

import (
	"fmt"
	"os"

	"github.com/dpcat237/go-dsu/internal/analyzer"

	"github.com/spf13/cobra"

	"github.com/dpcat237/go-dsu/internal/download"
	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/license"
	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
	"github.com/dpcat237/go-dsu/internal/vulnerability"
)

var (
	analyzeCmd = &cobra.Command{
		Use:   "analyze",
		Short: "Analyze current state",
		Long:  `Analyze licenses and vulnerabilities of current dependencies`,
		Run: func(cmd *cobra.Command, args []string) {
			analyze(cmd)
		},
	}
)

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().Bool("dev", false, "Development mode")
	analyzeCmd.Flags().String("path", "", "Preview project from git path. Eg. github.com/spf13/cobra")
}

func analyze(cmd *cobra.Command) {
	mod := output.ModeProd
	if cmd.Flag("dev").Value.String() == "true" {
		mod = output.ModeDev
	}

	fmt.Println("Analyzing dependencies...")
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

	licHnd := license.InitHandler(lgr)
	dwnHnd := download.InitHandler(exc, lgr)
	vlnHnd := vulnerability.InitHandler(lgr)
	hnd := module.InitHandler(exc)
	anz := analyzer.Init(dwnHnd, exc, lgr, licHnd, hnd, vlnHnd)
	out := anz.AnalyzeDependencies(cmd.Flag("path").Value.String())
	fmt.Println(out.ToString(mod))
}
