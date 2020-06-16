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
	cmdPermissions = "(cd %s && go list -m -mod=mod -json)"

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
		return rsp, out.WithErrorString(fmt.Sprintf("Error %s executing %s ", err.Error(), cmdStr))
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
		return rsp, out.WithErrorString(fmt.Sprintf("Error %s executing %s ", err.Error(), cmdStr))
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

// FolderAccessible verifies that provided folder is accessible and allow commands execution
func (exc Executor) FolderAccessible(pth string) bool {
	if pth == "" {
		return false
	}

	if _, err := os.Stat(pth); os.IsNotExist(err) {
		return false
	}

	if _, out := exc.ExecGlobal(fmt.Sprintf(cmdPermissions, pth)); out.HasError() {
		return false
	}
	return true
}

// UpdateProjectPath defines projects path
func (exc *Executor) UpdateProjectPath(prjPath string) {
	exc.projectPath = prjPath
}
