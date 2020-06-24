package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	pkg = "executor"
)

//Executor executes CLI commands
type Executor struct {
	goPath      string
	lgr         *logger.Logger
	projectPath string
}

// Init initializes CLI commands executor
func Init(lgr *logger.Logger) (*Executor, output.Output) {
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
	exc.lgr = lgr
	exc.projectPath = prjPath

	return &exc, out
}

// ExecGlobal executes requested command and returns Response and output.Output
func (exc Executor) ExecGlobal(cmdStr string) (Response, output.Output) {
	out := output.Create(fmt.Sprintf("%s.%s '%s'", pkg, "ExecGlobal", cmdStr))
	var rsp Response
	var cmdOut, cmdErr bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmd.Env = os.Environ()
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	if err := cmd.Run(); err != nil {
		return rsp, out.WithErrorString(fmt.Sprintf("Error executing %s with output: %s%s", cmdStr, cmdErr.Bytes(), cmdOut.Bytes()))
	}
	rsp.StdOutput = cmdOut.Bytes()
	rsp.StdError = cmdErr.Bytes()

	return rsp, out
}

// ExecProject executes requested command in project's folder and returns Response, and output.Output
func (exc Executor) ExecProject(atr string) (Response, output.Output) {
	out := output.Create(fmt.Sprintf("%s.%s '%s'", pkg, "ExecProject", atr))
	var rsp Response
	var cmdOut, cmdErr bytes.Buffer
	cmdStr := fmt.Sprintf("(cd %s/ && %s %s)", exc.projectPath, exc.goPath, atr)

	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmd.Env = os.Environ()
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	if err := cmd.Run(); err != nil {
		return rsp, out.WithErrorString(fmt.Sprintf("Error executing %s with output: %s%s", cmdStr, cmdErr.Bytes(), cmdOut.Bytes()))
	}
	rsp.StdOutput = cmdOut.Bytes()
	rsp.StdError = cmdErr.Bytes()

	return rsp, out
}

// ExistsInProject checks if file/folder of project exists
func (exc Executor) ExistsInProject(pth string) bool {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", exc.projectPath, pth)); os.IsNotExist(err) {
		return false
	}
	return true
}

//PromptConfirmation display prompt for confirmation
func (exc Executor) PromptConfirmation(text string) bool {
	var rsp string
	fmt.Println(text)
	if _, err := fmt.Scanf("%s", &rsp); err != nil {
		exc.lgr.Debug("Error displaying prompt: " + err.Error())
	}
	return exc.validatePromptConfirmation(rsp)
}

// UpdateProjectPath defines projects path
func (exc *Executor) UpdateProjectPath(prjPath string) {
	exc.projectPath = prjPath
}

func (exc Executor) promptConfirmationCorrect() bool {
	var rsp string
	fmt.Println("To proceed, enter y or n:")
	if _, err := fmt.Scanf("%s", &rsp); err != nil {
		exc.lgr.Debug("Error displaying prompt: " + err.Error())
	}
	return exc.validatePromptConfirmation(rsp)
}

func (exc Executor) validatePromptConfirmation(rsp string) bool {
	if rsp == "y" || rsp == "yes" {
		return true
	}
	if rsp == "n" || rsp == "no" {
		return false
	}
	return exc.promptConfirmationCorrect()
}
