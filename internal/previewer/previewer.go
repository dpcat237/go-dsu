package previewer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/schollz/progressbar/v3"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	cmdModDownload = "go mod download -json"

	pkg = "previewer"
)

//Preview handles changes preview processes
type Preview struct {
	exc   *executor.Executor
	lgr   *logger.Logger
	mdHnd *module.Handler
}

//Init initializes Preview
func Init(exc *executor.Executor, lgr *logger.Logger, mdHnd *module.Handler) *Preview {
	return &Preview{
		exc:   exc,
		lgr:   lgr,
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

	var wg sync.WaitGroup
	tt := len(mds)
	each := 90 / tt
	for k := range mds {
		wg.Add(1)
		go func(md *module.Module, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
			defer wg.Done()
			dfs, dfsOut := hnd.mdHnd.AnalyzeUpdateDifferences(*md)
			if dfsOut.HasError() {
				hnd.lgr.Debug(dfsOut.String())
			}

			if len(dfs) > 0 {
				md.UpdateDifferences = dfs
			}

			if err := bar.Add(each); err != nil {
				hnd.lgr.Debug(err.Error())
			}
		}(&mds[k], &wg, bar)
	}

	wg.Wait()
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
