package updater

import (
	"fmt"

	"github.com/dpcat237/go-dsu/internal/compare"
	"github.com/dpcat237/go-dsu/internal/download"
	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/logger"
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

//Updater manages execution of update processes
type Updater struct {
	cmpHnd *compare.Handler
	dwnHnd download.Handler
	exc    *executor.Executor
	lgr    logger.Logger
	mdHnd  module.Handler
}

// Init initializes updater handler
func Init(cmpHnd *compare.Handler, dwnHnd download.Handler, exc *executor.Executor, lgr logger.Logger, mdHnd module.Handler) *Updater {
	return &Updater{
		cmpHnd: cmpHnd,
		dwnHnd: dwnHnd,
		exc:    exc,
		lgr:    lgr,
		mdHnd:  mdHnd,
	}
}

// UpdateModules clean and update dependencies
func (upd Updater) UpdateModules(opt UpdateOptions) output.Output {
	out := output.Create(pkg + ".UpdateModules")

	upd.dwnHnd.CleanTemporaryData()
	mds, mdsOut := upd.mdHnd.ListAvailable(true, true)
	if mdsOut.HasError() {
		return mdsOut
	}
	if len(mds) == 0 {
		return out.WithResponse("All dependencies up to date")
	}

	if opt.IsIndirect {
		return upd.updateAll()
	}

	if opt.IsSelect {
		if err := mds.SelectCLI(); err != nil {
			return out.WithErrorString("Error rendering modules selector")
		}
	}

	if opt.IsTests {
		if tstOut := upd.runLocalTests(); tstOut.HasError() {
			return out.WithErrorString("Updater aborted because project's tests fail")
		}
	}

	if clsOut := upd.cmpHnd.InitializeClassifiers(); out.HasError() {
		return clsOut
	}
	defer upd.dwnHnd.CleanTemporaryData()

	return upd.updateDirectModules(mds, opt)
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
func (upd Updater) updateDirectModule(m module.Module, opt UpdateOptions, vnd bool) output.Output {
	out := output.Create(pkg + ".updateModule")
	excRsp, cmdOut := upd.exc.ExecProject(fmt.Sprintf("get %s", m.NewModule()))
	if cmdOut.HasError() {
		return cmdOut
	}
	if excRsp.HasError() {
		return out.WithErrorString(excRsp.StdErrorString())
	}

	if opt.IsTests {
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

	if opt.IsVerbose {
		fmt.Printf("Updated %s  ->  %s\n", m, m.NewModule())
	}
	return out
}

// updateDirectModules updates direct modules
func (upd Updater) updateDirectModules(mds module.Modules, opt UpdateOptions) output.Output {
	out := output.Create(pkg + ".updateAll")
	vnd := upd.exc.ExistsInProject(vendorFolder)

	for _, md := range mds {
		if !opt.IsPrompt {
			if mOut := upd.updateDirectModule(md, opt, vnd); mOut.HasError() {
				return mOut
			}
			continue
		}

		dfs, dfsOut := upd.cmpHnd.AnalyzeUpdateDifferences(md)
		if dfsOut.HasError() {
			upd.lgr.Debug(dfsOut.String())
		}
		if len(dfs) == 0 {
			if mOut := upd.updateDirectModule(md, opt, vnd); mOut.HasError() {
				return mOut
			}
			continue
		}

		md.UpdateDifferences = dfs
		fmt.Println(module.Modules{md}.ToPreviewTable())
		if !upd.exc.PromptConfirmation("Update this module? (y/n):") {
			continue
		}
		if mOut := upd.updateDirectModule(md, opt, vnd); mOut.HasError() {
			return mOut
		}
	}

	if vnd && !opt.IsTests {
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
