package cmd

import (
	"fmt"

	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/dpcat237/go-dsu/internal/cleaner"
	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/mod"
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
}

func preview(cmd *cobra.Command) {
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

	cln := cleaner.Init(exc)
	hnd := mod.InitHandler(exc)
	upd := previewer.Init(cln, exc, hnd)
	out = upd.Preview()
	fmt.Println(out.ToString(md))
}
