package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/dpcat237/go-dsu/internal/output"
)

const pkg = "executor"

type Executor struct {
	goPath      string
	projectPath string
}

func Init() (*Executor, output.Output) {
	out := output.Create(pkg + ".init")
	var exc Executor

	goExecPath, err := exec.LookPath("go")
	if err != nil {
		return &exc, out.WithError(err)
	}
	prjPath, err := os.Getwd()
	if err != nil {
		return &exc, out.WithError(err)
	}

	exc.goPath = goExecPath
	exc.projectPath = prjPath

	return &exc, out
}

// Exec executes requested command and returns STD output and error.
func (exc Executor) Exec(atr string) (string, string, output.Output) {
	out := output.Create(pkg + ".exec")
	var cmdOut, cmdErr bytes.Buffer
	cmdStr := fmt.Sprintf("(cd %s/ && %s %s)", exc.projectPath, exc.goPath, atr)

	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmd.Env = os.Environ()
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	if err := cmd.Run(); err != nil {
		return "", "", out.WithError(err)
	}
	return cmdOut.String(), cmdErr.String(), out
}
