package previewer

import (
	"fmt"

	"github.com/dpcat237/go-dsu/internal/cleaner"
	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	pkg = "previewer"
)

type Preview struct {
	cln *cleaner.Cleaner
	exc *executor.Executor
	hnd *module.Handler
}

func Init(cln *cleaner.Cleaner, exc *executor.Executor, hnd *module.Handler) *Preview {
	return &Preview{
		cln: cln,
		exc: exc,
		hnd: hnd,
	}
}

// Preview returns available updates of direct modules
func (prv Preview) Preview() output.Output {
	out := output.Create(pkg + ".Preview")
	fmt.Println("Discovering modules...")

	mds, mdsOut := prv.hnd.ListAvailable(true)
	if mdsOut.HasError() {
		return mdsOut
	}

	if len(mds) == 0 {
		return out.WithResponse("All dependencies up to date")
	}

	for k, md := range mds {
		dfs, dfsOut := prv.hnd.AnalyzeUpdateDifferences(md)
		if dfsOut.HasError() {
			return dfsOut
		}

		if len(dfs) > 0 {
			mds[k].UpdateDifferences = dfs
		}
	}

	return out.WithResponse(mds.ToTable())
}
