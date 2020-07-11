package previewer

import (
	"fmt"
	"sync"

	"github.com/schollz/progressbar/v3"

	"github.com/dpcat237/go-dsu/internal/compare"
	"github.com/dpcat237/go-dsu/internal/download"
	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
)

const pkg = "previewer"

//Preview handles changes preview processes
type Preview struct {
	cmpHnd *compare.Handler
	dwnHnd download.Handler
	exc    *executor.Executor
	lgr    logger.Logger
	mdHnd  module.Handler
}

//Init initializes Preview
func Init(cmpHnd *compare.Handler, dwnHnd download.Handler, exc *executor.Executor, lgr logger.Logger, mdHnd module.Handler) *Preview {
	return &Preview{
		cmpHnd: cmpHnd,
		dwnHnd: dwnHnd,
		exc:    exc,
		lgr:    lgr,
		mdHnd:  mdHnd,
	}
}

// Preview returns available updates of direct modules
func (hnd Preview) Preview(pth string) output.Output {
	out := output.Create(pkg + ".Preview")
	bar := progressbar.Default(100)

	hnd.dwnHnd.CleanTemporaryData()
	fmt.Println("Analyzing dependencies...")
	hnd.addProgress(bar, 5)

	if pthOut := hnd.updateProjectPath(pth); pthOut.HasError() {
		return pthOut
	}

	mds, mdsOut := hnd.mdHnd.ListAvailable(false, true)
	if mdsOut.HasError() {
		return mdsOut
	}

	if len(mds) == 0 {
		return out.WithResponse("All dependencies up to date")
	}
	hnd.addProgress(bar, 5)

	if clsOut := hnd.cmpHnd.InitializeClassifiers(); clsOut.HasError() {
		return clsOut
	}

	return hnd.processPreview(mds, bar)
}

func (hnd Preview) addProgress(bar *progressbar.ProgressBar, num int) {
	if err := bar.Add(num); err != nil {
		hnd.lgr.Debug(err.Error())
	}
}

func (hnd Preview) analyzeModuleGoroutine(md *module.Module, wg *sync.WaitGroup, bar *progressbar.ProgressBar, each int) {
	dfs, dfsOut := hnd.cmpHnd.AnalyzeUpdateDifferences(*md)
	if dfsOut.HasError() {
		if dfsOut.IsToManyRequests() {
			hnd.lgr.Fatal(dfsOut.String())
		}
		hnd.lgr.Debug(dfsOut.String())
	}

	if len(dfs) > 0 {
		md.UpdateDifferences = dfs
	}

	hnd.addProgress(bar, each)
	wg.Done()
}

func (hnd Preview) processPreview(mds module.Modules, bar *progressbar.ProgressBar) output.Output {
	out := output.Create(pkg + ".processPreview")
	defer hnd.dwnHnd.CleanTemporaryData()

	var wg sync.WaitGroup
	tt := len(mds)
	each := 90 / tt
	for k := range mds {
		wg.Add(1)
		go hnd.analyzeModuleGoroutine(&mds[k], &wg, bar, each)
	}

	wg.Wait()
	hnd.addProgress(bar, 90-each*tt)
	tbl := module.NewTable()

	return out.WithResponse(tbl.GeneratePreviewTable(mds))
}

func (hnd Preview) updateProjectPath(mdPth string) output.Output {
	out := output.Create(pkg + ".updateProjectPath")
	if mdPth == "" {
		return out
	}

	dir, dirOut := hnd.dwnHnd.DownloadModule(mdPth)
	if dirOut.HasError() {
		return dirOut
	}
	hnd.exc.UpdateProjectPath(dir)

	return out
}
