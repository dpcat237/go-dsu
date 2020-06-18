package previewer

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

const (
	pkg = "previewer"
)

//Preview handles changes preview processes
type Preview struct {
	dwnHnd *download.Handler
	exc    *executor.Executor
	lgr    *logger.Logger
	licHnd *license.Handler
	mdHnd  *module.Handler
	vlnHnd *vulnerability.Handler
}

//Init initializes Preview
func Init(dwnHnd *download.Handler, exc *executor.Executor, lgr *logger.Logger, licHnd *license.Handler, mdHnd *module.Handler, vlnHnd *vulnerability.Handler) *Preview {
	return &Preview{
		dwnHnd: dwnHnd,
		exc:    exc,
		lgr:    lgr,
		licHnd: licHnd,
		mdHnd:  mdHnd,
		vlnHnd: vlnHnd,
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
			dfs, dfsOut := hnd.analyzeUpdateDifferences(*md)
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

func (hnd Preview) addLicenseDifferences(md, mdUp module.Module, dffs *module.Differences) output.Output {
	out := output.Create(pkg + ".addLicenseDifferences")

	// Checks if in updated module are some changes in license
	md.License = hnd.licHnd.FindLicense(md.Dir)
	mdUp.License = hnd.licHnd.FindLicense(mdUp.Dir)

	if !md.License.Found() && !mdUp.License.Found() {
		dffs.AddModules(md, mdUp, module.DiffWeightLow, module.DiffTypeLicenseNotFound)
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, module.DiffTypeLicenseNotFound))
		return out
	}

	if md.License.Hash == mdUp.License.Hash {
		return out
	}

	if md.License.Found() && !mdUp.License.Found() {
		dffs.AddModules(md, mdUp, module.DiffWeightHigh, module.DiffTypeLicenseRemoved)
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, module.DiffTypeLicenseRemoved))
		return out
	}

	if !md.License.Found() && mdUp.License.Found() {
		dffs.AddModules(md, mdUp, module.DiffWeightHigh, module.DiffTypeLicenseAdded)
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, module.DiffTypeLicenseAdded))
		return out
	}

	// Identify license name and type
	hnd.licHnd.IdentifyType(&md.License)
	hnd.licHnd.IdentifyType(&mdUp.License)

	// Minor changes in the same license
	if md.License.Name == mdUp.License.Name {
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, module.DiffTypeLicenseMinorChanges))
		dffs.AddModules(md, mdUp, module.DiffWeightLow, module.DiffTypeLicenseMinorChanges)
		return out
	}

	// License name changed maintaining restrictiveness type
	if md.License.Type == mdUp.License.Type && md.License.Name != mdUp.License.Name {
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, module.DiffTypeLicenseNameChanged))
		dffs.AddModules(md, mdUp, module.DiffWeightMedium, module.DiffTypeLicenseNameChanged)
		return out
	}

	// License changed to less restrictive
	if !md.License.IsMoreRestrictive(mdUp.License.Type) {
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, module.DiffTypeLicenseLessStrictChanged))
		dffs.AddModules(md, mdUp, module.DiffWeightLow, module.DiffTypeLicenseLessStrictChanged)
		return out
	}

	// License changed to more restrictive with critical restrictiveness
	if mdUp.License.IsCritical() {
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, module.DiffTypeLicenseMoreStrictChanged))
		dffs.AddModules(md, mdUp, module.DiffWeightCritical, module.DiffTypeLicenseMoreStrictChanged)
		return out
	}

	// License changed to more restrictive
	hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, module.DiffTypeLicenseMoreStrictChanged))
	dffs.AddModules(md, mdUp, module.DiffWeightHigh, module.DiffTypeLicenseMoreStrictChanged)
	return out
}

func (hnd Preview) addNewModuleDifferences(upMd module.Module, dffs *module.Differences) output.Output {
	out := output.Create(pkg + ".addNewModuleDifferences")

	if out := hnd.updateDir(&upMd); out.HasError() {
		dffs.AddModule(upMd, module.DiffWeightHigh, module.DiffTypeModuleFetchError)
		return out
	}

	// Add vulnerabilities
	mdUpVlns, mdUpVlnsOut := hnd.vlnHnd.ModuleVulnerabilities(upMd.String())
	if mdUpVlnsOut.HasError() {
		return mdUpVlnsOut
	}
	for _, mdUpVln := range mdUpVlns {
		hnd.lgr.Debug(fmt.Sprintf("New module %s with vulnerability %s, type %d", upMd, mdUpVln.ID, module.DiffTypeNewVulnerability))
		dffs.AddVulnerability(upMd, mdUpVln)
	}

	// Add license
	upMd.License = hnd.licHnd.FindLicense(upMd.Dir)
	hnd.licHnd.IdentifyType(&upMd.License)
	if upMd.License.IsCritical() {
		dffs.AddModule(upMd, module.DiffWeightCritical, module.DiffTypeNewSubmodule)
		return out
	}
	dffs.AddModule(upMd, module.DiffWeightHigh, module.DiffTypeNewSubmodule)

	return out
}

func (hnd Preview) addVulnerabilityDifferences(md, mdUp module.Module, dffs *module.Differences) output.Output {
	out := output.Create(pkg + ".addVulnerabilityDifferences")

	// Checks if in updated module are some changes in vulnerabilities
	mdUpVlns, mdUpVlnsOut := hnd.vlnHnd.ModuleVulnerabilities(mdUp.String())
	if mdUpVlnsOut.HasError() {
		return mdUpVlnsOut
	}
	if len(mdUpVlns) == 0 {
		return out
	}
	mdVlns, mdVlnsOut := hnd.vlnHnd.ModuleVulnerabilities(md.String())
	if mdVlnsOut.HasError() {
		return mdVlnsOut
	}

	for _, mdUpVln := range mdUpVlns {
		if !mdVlns.HasVulnerability(mdUpVln.ID) {
			// Indirect module has new vulnerability
			hnd.lgr.Debug(fmt.Sprintf("New vulnerability %s in module %s, type %d", mdUpVln.ID, mdUp, module.DiffTypeNewVulnerability))
			dffs.AddVulnerability(mdUp, mdUpVln)
		}
	}
	return out
}

// AnalyzeUpdateDifferences analyze update differences of direct module and his direct modules recursively
func (hnd Preview) analyzeUpdateDifferences(md module.Module) (module.Differences, output.Output) {
	out := output.Create(pkg + ".AnalyzeUpdateDifferences")
	dffs := module.Differences{}

	if md.Update == nil {
		return dffs, out
	}

	upOut := hnd.updateDifferencesSubModule(md, *md.Update, &dffs)
	if upOut.HasError() {
		return dffs, out
	}
	return dffs, out
}

func (hnd Preview) updateDifferencesModule(md, mdUp module.Module, dffs *module.Differences) output.Output {
	out := output.Create(pkg + ".updateDifferencesModule")

	if licOut := hnd.addLicenseDifferences(md, mdUp, dffs); licOut.HasError() {
		return licOut
	}

	if vlnOut := hnd.addVulnerabilityDifferences(md, mdUp, dffs); vlnOut.HasError() {
		return vlnOut
	}

	return out
}

func (hnd Preview) updateDifferencesSubModule(md, mdUp module.Module, dffs *module.Differences) output.Output {
	out := output.Create(pkg + ".updateDifferencesSubModule")
	if out := hnd.updateDir(&md); out.HasError() {
		dffs.AddModule(md, module.DiffWeightHigh, module.DiffTypeModuleFetchError)
		return out
	}
	if out := hnd.updateDir(&mdUp); out.HasError() {
		dffs.AddModule(mdUp, module.DiffWeightHigh, module.DiffTypeModuleFetchError)
		return out
	}

	upOut := hnd.updateDifferencesModule(md, mdUp, dffs)
	if upOut.HasError() {
		return out
	}

	subMds, mdsOut := hnd.mdHnd.ListSubModules(md.Dir)
	if mdsOut.HasError() {
		return out
	}
	if len(subMds) == 0 {
		return out
	}
	subUpMds, mdsOut := hnd.mdHnd.ListSubModules(mdUp.Dir)
	if mdsOut.HasError() {
		return out
	}

	if subMdsOut := hnd.updateDifferencesSubModules(subMds, subUpMds, dffs); subMdsOut.HasError() {
		return subMdsOut
	}
	return out
}

func (hnd Preview) updateDifferencesSubModules(subMds, subUpMds module.Modules, dffs *module.Differences) output.Output {
	out := output.Create(fmt.Sprintf("%s.%s '%s'", pkg, "updateDifferencesSubModules", subMds))
	for _, upMd := range subUpMds {
		if upMd.Indirect {
			continue
		}

		found := false
		for _, md := range subMds {
			if md.PathCleaned() != upMd.PathCleaned() {
				continue
			}

			found = true
			if md.Version == upMd.Version {
				break
			}
			if cmpOut := hnd.updateDifferencesSubModule(md, upMd, dffs); cmpOut.HasError() {
				return cmpOut
			}
		}

		if !found {
			if out := hnd.addNewModuleDifferences(upMd, dffs); out.HasError() {
				return out
			}
		}
	}
	return out
}

//updateDir checks that module's directory is accessible and downloads if it isn't
func (hnd Preview) updateDir(md *module.Module) output.Output {
	out := output.Create(fmt.Sprintf("%s.%s '%s'", pkg, "updateDir", md))
	if md.Dir == "" || !hnd.exc.FolderAccessible(md.Dir) {
		dir, dirOut := hnd.dwnHnd.DownloadModule(md.String())
		if dirOut.HasError() {
			return dirOut
		}
		if dir == "" {
			return dirOut.WithErrorString(fmt.Sprintf("%s.%s Empty dir after download module %s", pkg, "updateDir", md))
		}
		md.Dir = dir
	}
	return out
}

func (hnd Preview) updateProjectPath(mdPth string) output.Output {
	out := output.Create(pkg + ".updateProjectPath")

	dir, dirOut := hnd.dwnHnd.DownloadModule(mdPth)
	if dirOut.HasError() {
		return dirOut
	}
	hnd.exc.UpdateProjectPath(dir)

	return out
}
