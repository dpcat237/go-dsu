package updater

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dpcat237/go-dsu/internal/cleaner"
	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/mod"
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
	cln *cleaner.Cleaner
	exc *executor.Executor
	hnd *mod.Handler
}

func Init(cln *cleaner.Cleaner, exc *executor.Executor, hnd *mod.Handler) *Updater {
	return &Updater{
		cln: cln,
		exc: exc,
		hnd: hnd,
	}
}

// UpdateModules clean and update dependencies
func (upd Updater) UpdateModules(all, sct, tst, vrb bool) output.Output {
	out := output.Create(pkg + ".UpdateModules")

	if outCln := upd.cln.Clean(); outCln.HasError() {
		return outCln
	}

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
