package cleaner

import (
	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	cmdModTidy = "mod tidy"

	pkg = "cleaner"
)

type Cleaner struct {
	exc *executor.Executor
}

func Init(exc *executor.Executor) *Cleaner {
	return &Cleaner{
		exc: exc,
	}
}

// Cleaner adds missing and remove unused modules
func (cln Cleaner) Clean() output.Output {
	out := output.Create(pkg + ".Cleaner")

	excRsp, cmdOut := cln.exc.Exec(cmdModTidy)
	if cmdOut.HasError() {
		return cmdOut
	}
	if excRsp.HasError() {
		return out.WithErrorString(excRsp.StdErrorString())
	}

	return out.WithResponse("Mod clean")
}
