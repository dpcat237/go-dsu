package updater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/mod"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	cmdGetUpdate   = "get -u -t"
	cmdListModules = "list -u -m -mod=mod -json all"
	cmdModTidy     = "mod tidy"
	cmdModVendor   = "mod vendor"
	cmdTestLocal   = "test $(go list ./... | grep -v /vendor/)"

	pkg          = "updater"
	vendorFolder = "vendor"
)

type Updater struct {
	exc *executor.Executor
}

func Init(exc *executor.Executor) *Updater {
	return &Updater{
		exc: exc,
	}
}

// Clean adds missing and remove unused modules
func (upd Updater) Clean() output.Output {
	out := output.Create(pkg + ".Clean")

	excRsp, cmdOut := upd.exc.Exec(cmdModTidy)
	if cmdOut.HasError() {
		return cmdOut
	}
	if excRsp.HasError() {
		return out.WithErrorString(excRsp.StdErrorString())
	}

	return out.WithResponse("Mod clean")
}

// Preview returns available updates of direct modules
func (upd Updater) Preview() output.Output {
	out := output.Create(pkg + ".Preview")

	if outCln := upd.Clean(); outCln.HasError() {
		return outCln.WithErrorPrefix("Actions done during clean up")
	}

	mds, mdsOut := upd.listAvailable(true)
	if mdsOut.HasError() {
		return mdsOut
	}

	if len(mds) == 0 {
		return out.WithResponse("All dependencies up to date")
	}
	return out.WithResponse(mds.ToTable())
}

// UpdateModules clean and update dependencies
func (upd Updater) UpdateModules(all, sct, tst, vrb bool) output.Output {
	out := output.Create(pkg + ".UpdateModules")

	if outCln := upd.Clean(); outCln.HasError() {
		return outCln
	}

	mds, mdsOut := upd.listAvailable(true)
	if mdsOut.HasError() {
		return mdsOut
	}
	if len(mds) == 0 {
		return out.WithResponse("All dependencies up to date")
	}

	if all {
		return upd.updateAll()
	}

	if sct {
		if err := mds.SelectCLI(); err != nil {
			return out.WithErrorString("Error rendering modules selector")
		}
	}

	if tst {
		if tstOut := upd.runLocalTests(); tstOut.HasError() {
			return out.WithErrorString("Update aborted because project's tests fail")
		}
	}

	return upd.updateDirectModules(mds, tst, vrb)
}

// listAvailable list modules with available updates
func (upd Updater) listAvailable(direct bool) (mod.Modules, output.Output) {
	var mds mod.Modules
	out := output.Create(pkg + ".listAvailable")

	excRsp, cmdOut := upd.exc.Exec(cmdListModules)
	if cmdOut.HasError() {
		return mds, cmdOut
	}
	if excRsp.HasError() {
		return mds, out.WithErrorString(excRsp.StdErrorString())
	}
	if excRsp.IsEmpty() {
		return mds, out.WithErrorString("Not found any dependency")
	}

	dec := json.NewDecoder(bytes.NewReader(excRsp.StdOutput))
	for {
		var m mod.Module
		if err := dec.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return mds, out.WithError(err)
		}

		if m.Main || (direct && m.Indirect) || m.Update == nil {
			continue
		}
		mds = append(mds, m)
	}

	return mds, out
}

func (upd Updater) runLocalTests() output.Output {
	out := output.Create(pkg + ".runLocalTests")

	_, cmdOut := upd.exc.Exec(cmdTestLocal)
	if cmdOut.HasError() {
		return cmdOut
	}
	return out
}

func (upd Updater) hasLicense() output.Output {
	out := output.Create(pkg + ".runLocalTests")

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	return out
}

func (upd Updater) rollback(m mod.Module) output.Output {
	out := output.Create(pkg + ".rollback")
	excRsp, cmdOut := upd.exc.Exec(fmt.Sprintf("get %s", m))
	if cmdOut.HasError() {
		return cmdOut
	}
	if excRsp.HasError() {
		return out.WithErrorString(excRsp.StdErrorString())
	}
	return out
}

// updateAll updates direct and indirect modules
func (upd Updater) updateAll() output.Output {
	out := output.Create(pkg + ".updateAll")

	excRsp, cmdOut := upd.exc.Exec(cmdGetUpdate)
	if cmdOut.HasError() {
		return cmdOut
	}
	if excRsp.HasError() {
		return out.WithErrorString(excRsp.StdErrorString())
	}

	return out.WithResponse(fmt.Sprintf("Successfully updated: \n %s", excRsp.StdOutputString()))
}

// updateDirectModule updates direct module
func (upd Updater) updateDirectModule(m mod.Module, tst, vnd, vrb bool) output.Output {
	out := output.Create(pkg + ".updateModule")
	excRsp, cmdOut := upd.exc.Exec(fmt.Sprintf("get %s", m.NewModule()))
	if cmdOut.HasError() {
		return cmdOut
	}
	if excRsp.HasError() {
		return out.WithErrorString(excRsp.StdErrorString())
	}

	if tst {
		if vnd {
			if out := upd.updateVendor(); out.HasError() {
				return out
			}
		}

		if tstOut := upd.runLocalTests(); tstOut.HasError() {
			fmt.Printf("Kept %s because with %s failed tests \n", m, m.NewModule())
			return upd.rollback(m)
		}
	}

	if vrb {
		fmt.Printf("Updated %s  ->  %s\n", m, m.NewModule())
	}
	return out
}

// updateDirectModules updates direct modules
func (upd Updater) updateDirectModules(mds mod.Modules, tst, vrb bool) output.Output {
	out := output.Create(pkg + ".updateAll")
	vnd := upd.exc.IsProjectFileExists(vendorFolder)

	for _, m := range mds {
		if mOut := upd.updateDirectModule(m, tst, vnd, vrb); mOut.HasError() {
			return mOut
		}
	}

	if vnd && !tst {
		if out := upd.updateVendor(); out.HasError() {
			return out
		}
	}

	return out.WithResponse("Updated successfully")
}

func (upd Updater) updateVendor() output.Output {
	out := output.Create(pkg + ".updateVendor")
	excRsp, cmdOut := upd.exc.Exec(cmdModVendor)
	if cmdOut.HasError() {
		return cmdOut
	}
	if excRsp.HasError() {
		return out.WithErrorString(excRsp.StdErrorString())
	}
	return out
}
