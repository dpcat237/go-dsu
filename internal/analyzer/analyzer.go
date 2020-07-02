package analyzer

import (
	"fmt"
	"sync"

	"github.com/schollz/progressbar/v3"

	"github.com/dpcat237/go-dsu/internal/download"
	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/license"
	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
	"github.com/dpcat237/go-dsu/internal/vulnerability"
)

const pkg = "analyzer"

//Analyze analyzes current dependencies
type Analyze struct {
	dwnHnd download.Handler
	exc    *executor.Executor
	lgr    logger.Logger
	licHnd license.Handler
	mdHnd  module.Handler
	vlnHnd vulnerability.Handler
}

//Init initializes Analyzer
func Init(dwnHnd download.Handler, exc *executor.Executor, lgr logger.Logger, licHnd license.Handler, mdHnd module.Handler, vlnHnd vulnerability.Handler) *Analyze {
	return &Analyze{
		dwnHnd: dwnHnd,
		exc:    exc,
		lgr:    lgr,
		licHnd: licHnd,
		mdHnd:  mdHnd,
		vlnHnd: vlnHnd,
	}
}

//AnalyzeDependencies analyzes state of current dependencies
func (hnd Analyze) AnalyzeDependencies(pth string) output.Output {
	out := output.Create(pkg + ".AnalyzeDependencies")
	bar := progressbar.Default(100)

	hnd.dwnHnd.CleanTemporaryData()
	hnd.addProgress(bar, 5)

	if pthOut := hnd.updateProjectPath(pth); pthOut.HasError() {
		return pthOut
	}

	mds, mdsOut := hnd.mdHnd.ListAvailable(true, false)
	if mdsOut.HasError() {
		return mdsOut
	}

	if len(mds) == 0 {
		return out.WithResponse("All dependencies up to date")
	}
	hnd.addProgress(bar, 5)

	if licOut := hnd.licHnd.InitializeClassifier(); out.HasError() {
		return licOut
	}

	return hnd.processAnalyzeDependencies(mds, bar)
}

func (hnd Analyze) addProgress(bar *progressbar.ProgressBar, num int) {
	if err := bar.Add(num); err != nil {
		hnd.lgr.Debug(err.Error())
	}
}

func (hnd Analyze) analyzeModule(md *module.Module) output.Output {
	out := output.Create(pkg + ".analyzeModule")

	if md.Dir == "" || !hnd.dwnHnd.FolderAccessible(md.Dir) {
		dir, dirOut := hnd.dwnHnd.DownloadModule(md.String())
		if dirOut.HasError() {
			return dirOut
		}
		if dir == "" {
			return dirOut.WithErrorString(fmt.Sprintf("%s.%s Empty dir after download module %s", pkg, "updateDir", md))
		}
		md.Dir = dir
	}

	md.License = hnd.licHnd.FindLicense(md.Dir)
	mdVlns, mdVlnsOut := hnd.vlnHnd.ModuleVulnerabilities(md.String())
	if mdVlnsOut.HasError() {
		return mdVlnsOut
	}
	md.Vulnerabilities = mdVlns

	subMds, subMdsOut := hnd.analyzeModuleDependencies(*md)
	if subMdsOut.HasError() {
		return subMdsOut
	}
	md.Dependencies = subMds

	return out
}

func (hnd Analyze) analyzeModuleGoroutine(md *module.Module, wg *sync.WaitGroup, bar *progressbar.ProgressBar, each int) {
	if out := hnd.analyzeModule(md); out.HasError() {
		hnd.lgr.Debug(out.String())
	}
	hnd.addProgress(bar, each)
	wg.Done()
}

func (hnd Analyze) analyzeModuleDependencies(md module.Module) ([]module.Module, output.Output) {
	out := output.Create(pkg + ".analyzeModuleDependencies")
	var mds []module.Module
	subMds, mdsOut := hnd.mdHnd.ListSubModules(md.Dir)
	if mdsOut.HasError() {
		return mds, mdsOut
	}

	for _, subMd := range subMds {
		if out := hnd.analyzeModule(&subMd); out.HasError() {
			hnd.lgr.Debug(out.String())
		}
		mds = append(mds, subMd)
	}

	return mds, out
}

func (hnd Analyze) processAnalyzeDependencies(mds module.Modules, bar *progressbar.ProgressBar) output.Output {
	out := output.Create(pkg + ".processAnalyzeDependencies")
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

	return out.WithResponse(mds.ToAnalyzeTable())
}

func (hnd Analyze) updateProjectPath(mdPth string) output.Output {
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
