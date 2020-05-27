package updater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/mod"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	cmdGetUpdate   = "get -u -t"
	cmdListModules = "list -u -m -mod=mod -json all"
	cmdModTidy     = "mod tidy"
	cmdModVendor   = "mod vendor"

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
		out.SetPid(cmdOut.GetPid())
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

// UpdateDependencies clean and update dependencies
func (upd Updater) UpdateDependencies(all, sct bool) output.Output {
	out := output.Create(pkg + ".UpdateDependencies")

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

	return upd.updateDirect(mds)
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

// updateDirect updates direct modules
func (upd Updater) updateDirect(mds mod.Modules) output.Output {
	out := output.Create(pkg + ".updateAll")

	updateModule := func(m mod.Module, verbose bool) output.Output {
		out := output.Create(pkg + ".updateModule")
		excRsp, cmdOut := upd.exc.Exec(fmt.Sprintf("get %s", m.NewModule()))
		if cmdOut.HasError() {
			return cmdOut
		}
		if excRsp.HasError() {
			return out.WithErrorString(excRsp.StdErrorString())
		}

		if verbose {
			fmt.Printf("Updated %s  ->  %s\n", m, m.NewModule())
		}
		return out
	}

	for _, m := range mds {
		if mOut := updateModule(m, true); mOut.HasError() {
			return mOut
		}
	}

	// Add updated modules to vendor folder if it exists
	if upd.exc.IsProjectFileExists(vendorFolder) {
		excRsp, cmdOut := upd.exc.Exec(cmdModVendor)
		if cmdOut.HasError() {
			return cmdOut
		}
		if excRsp.HasError() {
			return out.WithErrorString(excRsp.StdErrorString())
		}
	}

	return out.WithResponse("Updated successfully")
}
