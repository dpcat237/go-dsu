package updater

import (
	"fmt"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	cmdGetUpdate = "get -u -t"
	cmdModVendor = "mod vendor"
	cmdTestLocal = "test $(go list ./... | grep -v /vendor/)"

	pkg          = "updater"
	vendorFolder = "vendor"
)

type Updater struct {
	exc *executor.Executor
	hnd *module.Handler
}

// Init initializes updater handler
func Init(exc *executor.Executor, hnd *module.Handler) *Updater {
	return &Updater{
		exc: exc,
		hnd: hnd,
	}
}

// UpdateModules clean and update dependencies
func (upd Updater) UpdateModules(all, sct, tst, vrb bool) output.Output {
	out := output.Create(pkg + ".UpdateModules")

	mds, mdsOut := upd.hnd.ListAvailable(true)
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
			return out.WithErrorString("Updater aborted because project's tests fail")
		}
	}

	return upd.updateDirectModules(mds, tst, vrb)
}

func (upd Updater) runLocalTests() output.Output {
	out := output.Create(pkg + ".runLocalTests")

	_, cmdOut := upd.exc.ExecProject(cmdTestLocal)
	if cmdOut.HasError() {
		return cmdOut
	}
	return out
}

func (upd Updater) rollback(m module.Module) output.Output {
	out := output.Create(pkg + ".rollback")
	excRsp, cmdOut := upd.exc.ExecProject(fmt.Sprintf("get %s", m))
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

	excRsp, cmdOut := upd.exc.ExecProject(cmdGetUpdate)
	if cmdOut.HasError() {
		return cmdOut
	}
	if excRsp.HasError() {
		return out.WithErrorString(excRsp.StdErrorString())
	}

	return out.WithResponse(fmt.Sprintf("Successfully updated: \n %s", excRsp.StdOutputString()))
}

// updateDirectModule updates direct module
func (upd Updater) updateDirectModule(m module.Module, tst, vnd, vrb bool) output.Output {
	out := output.Create(pkg + ".updateModule")
	excRsp, cmdOut := upd.exc.ExecProject(fmt.Sprintf("get %s", m.NewModule()))
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
func (upd Updater) updateDirectModules(mds module.Modules, tst, vrb bool) output.Output {
	out := output.Create(pkg + ".updateAll")
	vnd := upd.exc.ExistsInProject(vendorFolder)

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
	excRsp, cmdOut := upd.exc.ExecProject(cmdModVendor)
	if cmdOut.HasError() {
		return cmdOut
	}
	if excRsp.HasError() {
		return out.WithErrorString(excRsp.StdErrorString())
	}
	return out
}
