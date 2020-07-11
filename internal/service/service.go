package service

import (
	"github.com/dpcat237/go-dsu/internal/analyzer"
	"github.com/dpcat237/go-dsu/internal/compare"
	"github.com/dpcat237/go-dsu/internal/download"
	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/license"
	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
	"github.com/dpcat237/go-dsu/internal/previewer"
	"github.com/dpcat237/go-dsu/internal/updater"
	"github.com/dpcat237/go-dsu/internal/vulnerability"
)

//InitAnalyzer initializes required dependencies for analyzer
func InitAnalyzer(mod output.Mode, ossTkn string) (*analyzer.Analyze, output.Output) {
	var out output.Output
	lgr, lgrOut := logger.Init(mod)
	if lgrOut.HasError() {
		return nil, lgrOut
	}

	exc, excOut := executor.Init(lgr)
	if excOut.HasError() {
		return nil, excOut
	}

	licHnd := license.InitHandler(lgr)
	dwnHnd := download.InitHandler(exc, lgr)
	vlnHnd, vlnOut := vulnerability.InitHandler(lgr, ossTkn)
	if vlnOut.HasError() {
		return nil, vlnOut
	}
	mdHnd := module.InitHandler(exc)

	return analyzer.Init(dwnHnd, exc, lgr, licHnd, mdHnd, vlnHnd), out
}

//InitPreviewer initializes required dependencies for previewer
func InitPreviewer(mod output.Mode, ossTkn string) (*previewer.Preview, output.Output) {
	var out output.Output
	lgr, lgrOut := logger.Init(mod)
	if lgrOut.HasError() {
		return nil, lgrOut
	}

	exc, excOut := executor.Init(lgr)
	if excOut.HasError() {
		return nil, excOut
	}

	dwnHnd := download.InitHandler(exc, lgr)
	licHnd := license.InitHandler(lgr)
	mdHnd := module.InitHandler(exc)
	vlnHnd, vlnOut := vulnerability.InitHandler(lgr, ossTkn)
	if vlnOut.HasError() {
		return nil, vlnOut
	}
	cmpHnd := compare.Init(dwnHnd, lgr, licHnd, mdHnd, vlnHnd)

	return previewer.Init(cmpHnd, dwnHnd, exc, lgr, mdHnd), out
}

//InitUpdater initializes required dependencies for updater
func InitUpdater(mod output.Mode) (*updater.Updater, output.Output) {
	var out output.Output
	lgr, lgrOut := logger.Init(mod)
	if lgrOut.HasError() {
		return nil, lgrOut
	}

	exc, excOut := executor.Init(lgr)
	if excOut.HasError() {
		return nil, excOut
	}

	dwnHnd := download.InitHandler(exc, lgr)
	licHnd := license.InitHandler(lgr)
	mdHnd := module.InitHandler(exc)
	vlnHnd, vlnOut := vulnerability.InitHandler(lgr, "")
	if vlnOut.HasError() {
		return nil, vlnOut
	}
	cmpHnd := compare.Init(dwnHnd, lgr, licHnd, mdHnd, vlnHnd)

	return updater.Init(cmpHnd, dwnHnd, exc, lgr, mdHnd), out
}
