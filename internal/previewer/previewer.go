package previewer

import (
	"github.com/dpcat237/go-dsu/internal/cleaner"
	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/mod"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	pkg = "previewer"
)

type Preview struct {
	cln *cleaner.Cleaner
	exc *executor.Executor
	hnd *mod.Handler
}

func Init(cln *cleaner.Cleaner, exc *executor.Executor, hnd *mod.Handler) *Preview {
	return &Preview{
		cln: cln,
		exc: exc,
		hnd: hnd,
	}
}

// Preview returns available updates of direct modules
func (prv Preview) Preview() output.Output {
	out := output.Create(pkg + ".Preview")

	if outCln := prv.cln.Clean(); outCln.HasError() {
		return outCln.WithErrorPrefix("Actions done during clean up")
	}

	mds, mdsOut := prv.hnd.ListAvailable(true)
	if mdsOut.HasError() {
		return mdsOut
	}

	if len(mds) == 0 {
		return out.WithResponse("All dependencies up to date")
	}
	return out.WithResponse(mds.ToTable())
}
