package previewer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/schollz/progressbar/v3"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	cmdModDownload = "go mod download -json"

	pkg = "previewer"
)

type Preview struct {
	exc   *executor.Executor
	mdHnd *module.Handler
}

func Init(exc *executor.Executor, mdHnd *module.Handler) *Preview {
	return &Preview{
		exc:   exc,
		mdHnd: mdHnd,
	}
}

// Preview returns available updates of direct modules
func (hnd Preview) Preview(pth string) output.Output {
	out := output.Create(pkg + ".Preview")
	bar := progressbar.Default(100)

	if err := bar.Add(5); err != nil {
		return out.WithError(err)
	}

	if pth != "" {
		if pthOut := hnd.updateProjectPath(pth); pthOut.HasError() {
			return pthOut
		}
	}

	fmt.Println("Discovering modules...")
	mds, mdsOut := hnd.mdHnd.ListAvailable(true)
	if mdsOut.HasError() {
		return mdsOut
	}

	if len(mds) == 0 {
		return out.WithResponse("All dependencies up to date")
	}

	if err := bar.Add(5); err != nil {
		return out.WithError(err)
	}

	tt := len(mds)
	each := 90 / tt
	for k, md := range mds {
		dfs, dfsOut := hnd.mdHnd.AnalyzeUpdateDifferences(md)
		if dfsOut.HasError() {
			return dfsOut
		}

		if len(dfs) > 0 {
			mds[k].UpdateDifferences = dfs
		}

		if err := bar.Add(each); err != nil {
			return out.WithError(err)
		}
	}
	if err := bar.Add(90 - each*tt); err != nil {
		return out.WithError(err)
	}

	return out.WithResponse(mds.ToTable())
}

func (hnd *Preview) updateProjectPath(path string) output.Output {
	out := output.Create(pkg + ".updateProjectPath")
	dwnRsp, dwnOut := hnd.exc.ExecGlobal(fmt.Sprintf("%s %s", cmdModDownload, path))
	if dwnOut.HasError() {
		return dwnOut
	}

	var mdDwn module.Module
	dec := json.NewDecoder(bytes.NewReader(dwnRsp.StdOutput))
	if err := dec.Decode(&mdDwn); err != nil {
		return out.WithError(err)
	}
	hnd.exc.UpdateProjectPath(mdDwn.Dir)

	return out
}
