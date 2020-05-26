package updater

import (
	"fmt"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	cmdClean  = "mod tidy"
	cmdUpdate = "get -u -t"

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

// UpdateDependencies clean and update dependencies
func (upd Updater) UpdateDependencies() output.Output {
	out := output.Create(pkg + ".updateDependencies")

	// Add missing and remove unused modules
	_, excErr, cmdOut := upd.exc.Exec(cmdClean)
	if cmdOut.HasError() {
		return cmdOut
	}
	if excErr != "" {
		return out.WithErrorString(excErr)
	}

	// Update modules
	excOut, excErr, cmdOut := upd.exc.Exec(cmdUpdate)
	if cmdOut.HasError() {
		return cmdOut
	}
	if excErr != "" {
		return out.WithErrorString(excErr)
	}

	return out.WithResponse(fmt.Sprintf("Successfully updated: \n %s", excOut))
}
