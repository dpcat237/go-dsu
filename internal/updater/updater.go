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
	cmdGetUpdate    = "get -u -t"
	cmdListModules  = "list -u -m -mod=mod -json all"
	cmdModTidyClean = "mod tidy"

	pkg = "updater"
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

	_, excErr, cmdOut := upd.exc.ExecToString(cmdModTidyClean)
	if cmdOut.HasError() {
		return cmdOut
	}
	if excErr != "" {
		return out.WithErrorString(excErr)
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
func (upd Updater) UpdateDependencies(all bool) output.Output {
	out := output.Create(pkg + ".UpdateDependencies")

	if outCln := upd.Clean(); outCln.HasError() {
		return outCln
	}

	if all {
		return upd.updateAll()
	}

	return out
}

// listAvailable list modules with available updates
func (upd Updater) listAvailable(direct bool) (mod.Modules, output.Output) {
	var mds mod.Modules
	out := output.Create(pkg + ".listAvailable")

	excOut, excErr, cmdOut := upd.exc.ExecToBytes(cmdListModules)
	if cmdOut.HasError() {
		return mds, cmdOut
	}
	if len(excErr) > 0 {
		return mds, out.WithErrorString(string(excErr))
	}
	if len(excOut) == 0 {
		return mds, out.WithErrorString("Not found any dependency")
	}

	dec := json.NewDecoder(bytes.NewReader(excOut))
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

	excOut, excErr, cmdOut := upd.exc.ExecToString(cmdGetUpdate)
	if cmdOut.HasError() {
		return cmdOut
	}
	if excErr != "" {
		return out.WithErrorString(excErr)
	}

	return out.WithResponse(fmt.Sprintf("Successfully updated: \n %s", excOut))
}
