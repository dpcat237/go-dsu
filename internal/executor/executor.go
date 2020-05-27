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

// ExecToBytes executes requested command and returns STD output and error as bytes.
func (exc Executor) ExecToBytes(atr string) ([]byte, []byte, output.Output) {
	out := output.Create(pkg + ".ExecToBytes")

	cmdOut, cmdErr, err := exc.exec(atr)
	if err != nil {
		return []byte{}, []byte{}, out.WithError(err)
	}
	return cmdOut.Bytes(), cmdErr.Bytes(), out
}

// ExecToString executes requested command and returns STD output and error as string.
func (exc Executor) ExecToString(atr string) (string, string, output.Output) {
	out := output.Create(pkg + ".ExecToString")

	cmdOut, cmdErr, err := exc.exec(atr)
	if err != nil {
		return "", "", out.WithError(err)
	}
	return cmdOut.String(), cmdErr.String(), out
}

func (exc Executor) exec(atr string) (bytes.Buffer, bytes.Buffer, error) {
	var cmdOut, cmdErr bytes.Buffer
	cmdStr := fmt.Sprintf("(cd %s/ && %s %s)", exc.projectPath, exc.goPath, atr)

	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmd.Env = os.Environ()
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	return cmdOut, cmdErr, cmd.Run()
}
