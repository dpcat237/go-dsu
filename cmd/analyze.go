package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dpcat237/go-dsu/internal/output"
	"github.com/dpcat237/go-dsu/internal/service"
)

var (
	analyzeCmd = &cobra.Command{
		Use:   "analyze",
		Short: "Analyze current dependencies",
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
	analyzeCmd.Flags().String("ossemail", "", "OSS Index email")
	analyzeCmd.Flags().String("osstoken", "", "OSS Index API")
	analyzeCmd.Flags().String("oss", "", "OSS Index API encoded token base64(email:token)")
}

func analyze(cmd *cobra.Command) {
	mod := output.ModeProd
	if cmd.Flag("dev").Changed {
		mod = output.ModeDev
	}

	fmt.Println("Analyzing prerequisites...")
	if out := checkPrerequisites(); out.HasError() {
		fmt.Println(out.ToString(mod))
		return
	}

	anz, initOut := service.InitAnalyzer(mod, extractOSSToken(cmd))
	if initOut.HasError() {
		fmt.Println(initOut.ToString(mod))
		os.Exit(1)
	}

	out := anz.AnalyzeDependencies(cmd.Flag("path").Value.String())
	fmt.Println(out.ToString(mod))
}
