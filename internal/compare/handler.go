package compare

import (
	"fmt"

	"github.com/dpcat237/go-dsu/internal/download"
	"github.com/dpcat237/go-dsu/internal/license"
	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
	"github.com/dpcat237/go-dsu/internal/vulnerability"
)

const pkg = "compare"

//Handler compare module and his update to find differences
type Handler struct {
	dwnHnd *download.Handler
	lgr    *logger.Logger
	licHnd *license.Handler
	mdHnd  *module.Handler
	vlnHnd *vulnerability.Handler
}

//Init initializes compare Handler
func Init(dwnHnd *download.Handler, lgr *logger.Logger, licHnd *license.Handler, mdHnd *module.Handler, vlnHnd *vulnerability.Handler) *Handler {
	return &Handler{
		dwnHnd: dwnHnd,
		lgr:    lgr,
		licHnd: licHnd,
		mdHnd:  mdHnd,
		vlnHnd: vlnHnd,
	}
}

// AnalyzeUpdateDifferences analyze update differences of direct module and his direct modules recursively
func (hnd Handler) AnalyzeUpdateDifferences(md module.Module) (module.Differences, output.Output) {
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

func (hnd Handler) addLicenseDifferences(md, mdUp module.Module, dffs *module.Differences) output.Output {
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

func (hnd Handler) addNewModuleDifferences(upMd module.Module, dffs *module.Differences) output.Output {
	out := output.Create(pkg + ".addNewModuleDifferences")

	if out := hnd.updateModuleDirectory(&upMd); out.HasError() {
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

func (hnd Handler) addVulnerabilityDifferences(md, mdUp module.Module, dffs *module.Differences) output.Output {
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

func (hnd Handler) updateDifferencesModule(md, mdUp module.Module, dffs *module.Differences) output.Output {
	out := output.Create(pkg + ".updateDifferencesModule")

	if licOut := hnd.addLicenseDifferences(md, mdUp, dffs); licOut.HasError() {
		return licOut
	}

	if vlnOut := hnd.addVulnerabilityDifferences(md, mdUp, dffs); vlnOut.HasError() {
		return vlnOut
	}

	return out
}

func (hnd Handler) updateDifferencesSubModule(md, mdUp module.Module, dffs *module.Differences) output.Output {
	out := output.Create(pkg + ".updateDifferencesSubModule")
	if out := hnd.updateModuleDirectory(&md); out.HasError() {
		dffs.AddModule(md, module.DiffWeightHigh, module.DiffTypeModuleFetchError)
		return out
	}
	if out := hnd.updateModuleDirectory(&mdUp); out.HasError() {
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

func (hnd Handler) updateDifferencesSubModules(subMds, subUpMds module.Modules, dffs *module.Differences) output.Output {
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

//updateModuleDirectory checks that module's directory is accessible and downloads if it isn't
func (hnd Handler) updateModuleDirectory(md *module.Module) output.Output {
	out := output.Create(fmt.Sprintf("%s.%s '%s'", pkg, "updateDir", md))
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
	return out
}
