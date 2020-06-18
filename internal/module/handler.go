package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	cmdListModules    = "list -u -m -mod=mod -json all"
	cmdListSubModules = "(cd %s && go list -m -mod=mod -json all)"
)

//Handler handles functions related to modules
type Handler struct {
	exc *executor.Executor
}

//InitHandler initializes Module handler
func InitHandler(exc *executor.Executor) *Handler {
	return &Handler{
		exc: exc,
	}
}

// ListAvailable list modules with available updates
func (hnd Handler) ListAvailable(direct bool) (Modules, output.Output) {
	var mds Modules
	out := output.Create(pkg + ".ListAvailable")

	excRsp, cmdOut := hnd.exc.ExecProject(cmdListModules)
	if cmdOut.HasError() {
		return mds, cmdOut
	}
	if excRsp.HasError() {
		return mds, out.WithErrorString(excRsp.StdErrorString())
	}
	if excRsp.IsEmpty() {
		return mds, out.WithErrorString("Not found any dependency")
	}

	return hnd.bytesToModules(excRsp.StdOutput, direct, true)
}

//ListSubModules return submodules (indirect modules)
func (hnd Handler) ListSubModules(pth string) (Modules, output.Output) {
	out := output.Create(pkg + ".listSubModules")
	var mds Modules

	excRsp, cmdOut := hnd.exc.ExecGlobal(fmt.Sprintf(cmdListSubModules, pth))
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

func (hnd Handler) bytesToModules(rspBts []byte, direct, withUpdate bool) (Modules, output.Output) {
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
