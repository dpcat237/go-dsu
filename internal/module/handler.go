package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/output"
	"github.com/dpcat237/go-dsu/internal/version"
)

const (
	cmdListModules       = "list -u -m -json all"
	cmdListModulesMod    = "list -u -m -mod=mod -json all"
	cmdListSubModules    = "(cd %s && go list -m -json all)"
	cmdListSubModulesMod = "(cd %s && go list -m -mod=mod -json all)"
)

//Handler handles functions related to modules
type Handler interface {
	ListAvailable(direct, withUpdate bool) (Modules, output.Output)
	ListSubModules(pth string) (Modules, output.Output)
}

type handler struct {
	exc *executor.Executor
}

//InitHandler initializes Module handler
func InitHandler(exc *executor.Executor) *handler {
	return &handler{
		exc: exc,
	}
}

// ListAvailable list modules with available updates
func (hnd handler) ListAvailable(direct, withUpdate bool) (Modules, output.Output) {
	var mds Modules
	out := output.Create(pkg + ".ListAvailable")

	cmdStr := cmdListModulesMod
	if !version.IsModSupported() {
		cmdStr = cmdListModules
	}

	excRsp, cmdOut := hnd.exc.ExecProject(cmdStr)
	if cmdOut.HasError() {
		return mds, cmdOut
	}
	if excRsp.HasError() {
		return mds, out.WithErrorString(excRsp.StdErrorString())
	}
	if excRsp.IsEmpty() {
		return mds, out.WithErrorString("Not found any dependency")
	}

	return hnd.bytesToModules(excRsp.StdOutput, direct, withUpdate)
}

//ListSubModules return submodules (indirect modules)
func (hnd handler) ListSubModules(pth string) (Modules, output.Output) {
	out := output.Create(pkg + ".listSubModules")
	var mds Modules

	cmdStr := cmdListSubModulesMod
	if !version.IsModSupported() {
		cmdStr = cmdListSubModules
	}

	excRsp, cmdOut := hnd.exc.ExecGlobal(fmt.Sprintf(cmdStr, pth))
	if cmdOut.HasError() {
		return mds, cmdOut
	}
	if excRsp.HasError() {
		return mds, out.WithErrorString(excRsp.StdErrorString())
	}
	if excRsp.IsEmpty() {
		return mds, out.WithErrorString("Not found any dependency")
	}

	return hnd.bytesToModules(excRsp.StdOutput, true, false)
}

func (hnd handler) bytesToModules(rspBts []byte, direct, withUpdate bool) (Modules, output.Output) {
	var mds Modules
	out := output.Create(pkg + ".bytesToModules")
	dec := json.NewDecoder(bytes.NewReader(rspBts))
	for {
		var md Module
		if err := dec.Decode(&md); err != nil {
			if err == io.EOF {
				break
			}
			return mds, out.WithError(err)
		}

		if md.Main || (direct && md.Indirect) || (withUpdate && md.Update == nil) {
			continue
		}
		mds = append(mds, md)
	}

	return mds, out
}
