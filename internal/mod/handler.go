package mod

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	cmdListModules = "list -u -m -mod=mod -json all"
)

type Handler struct {
	exc *executor.Executor
}

func InitHandler(exc *executor.Executor) *Handler {
	return &Handler{
		exc: exc,
	}
}

// ListAvailable list modules with available updates
func (hnd Handler) ListAvailable(direct bool) (Modules, output.Output) {
	var mds Modules
	out := output.Create(pkg + ".listAvailable")

	excRsp, cmdOut := hnd.exc.Exec(cmdListModules)
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
		var m Module
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
