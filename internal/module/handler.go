package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/license"
	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	cmdChmodModule    = "(chmod 744 %s && chmod 655 %s/*)"
	cmdListModules    = "list -u -m -mod=mod -json all"
	cmdListSubModules = "(cd %s && go list -m -mod=mod -json all)"
	cmdModDownload    = "mod download -json"
)

type Handler struct {
	exc    *executor.Executor
	lgr    *logger.Logger
	licHnd *license.Handler
}

func InitHandler(exc *executor.Executor, lgr *logger.Logger, licHnd *license.Handler) *Handler {
	return &Handler{
		exc:    exc,
		lgr:    lgr,
		licHnd: licHnd,
	}
}

// AnalyzeUpdateDifferences analyze update differences of direct module and his direct modules recursively
func (hnd Handler) AnalyzeUpdateDifferences(md Module) (Differences, output.Output) {
	out := output.Create(pkg + ".AnalyzeUpdateDifferences")
	dffs := Differences{}

	if md.Update == nil {
		return dffs, out
	}

	upOut := hnd.updateDifferencesSubModule(md, *md.Update, &dffs)
	if upOut.HasError() {
		return dffs, out
	}
	return dffs, out
}

// ListAvailable list modules with available updates
func (hnd Handler) ListAvailable(direct bool) (Modules, output.Output) {
	var mds Modules
	out := output.Create(pkg + ".ListAvailable")

	excRsp, cmdOut := hnd.exc.ExecProject(cmdListModules)
	if cmdOut.HasError() {
		return mds, cmdOut
	}
	if excRsp.HasError() {
		return mds, out.WithErrorString(excRsp.StdErrorString())
	}
	if excRsp.IsEmpty() {
		return mds, out.WithErrorString("Not found any dependency")
	}

	return hnd.bytesToModules(excRsp.StdOutput, direct, true)
}

func (hnd Handler) bytesToModules(rspBts []byte, direct, withUpdate bool) (Modules, output.Output) {
	var mds Modules
	out := output.Create(pkg + ".bytesToModules")
	dec := json.NewDecoder(bytes.NewReader(rspBts))
	for {
		var md Module
		if err := dec.Decode(&md); err != nil {
			if err == io.EOF {
				break
			}
			return mds, out.WithError(err)
		}

		if md.Main || (direct && md.Indirect) || (withUpdate && md.Update == nil) {
			continue
		}
		mds = append(mds, md)
	}

	return mds, out
}

func (hnd Handler) downloadModule(md *Module) output.Output {
	out := output.Create(fmt.Sprintf("%s.%s '%s'", pkg, "downloadModule", md))

	// Download
	dwnRsp, dwnOut := hnd.exc.ExecProject(fmt.Sprintf("%s %s", cmdModDownload, md))
	if dwnOut.HasError() {
		return dwnOut
	}

	var mdDwn Module
	dec := json.NewDecoder(bytes.NewReader(dwnRsp.StdOutput))
	if err := dec.Decode(&mdDwn); err != nil {
		return out.WithError(err)
	}
	md.Dir = mdDwn.Dir

	// Double check permissions
	if _, prmOut := hnd.exc.ExecGlobal(fmt.Sprintf(cmdChmodModule, md.Dir, md.Dir)); prmOut.HasError() {
		return prmOut
	}

	return out
}

func (hnd Handler) listSubModules(pth string) (Modules, output.Output) {
	out := output.Create(pkg + ".listSubModules")
	var mds Modules

	excRsp, cmdOut := hnd.exc.ExecGlobal(fmt.Sprintf(cmdListSubModules, pth))
	if cmdOut.HasError() {
		return mds, cmdOut
	}
	if excRsp.HasError() {
		return mds, out.WithErrorString(excRsp.StdErrorString())
	}
	if excRsp.IsEmpty() {
		return mds, out.WithErrorString("Not found any dependency")
	}

	return hnd.bytesToModules(excRsp.StdOutput, true, false)
}

func (hnd Handler) updateDifferencesModule(md, mdUp Module, dffs *Differences) output.Output {
	out := output.Create(pkg + ".updateDifferencesModule")

	// Checks if in updated module are some changes in license
	md.License = hnd.licHnd.FindLicense(md.Dir)
	mdUp.License = hnd.licHnd.FindLicense(mdUp.Dir)

	if !md.License.Found() && !mdUp.License.Found() {
		dffs.AddDifference(md, mdUp, diffWeightLow, diffTypeLicenseNotFound)
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, diffTypeLicenseNotFound))
		return out
	}

	if md.License.Hash == mdUp.License.Hash {
		return out
	}

	if md.License.Found() && !mdUp.License.Found() {
		dffs.AddDifference(md, mdUp, diffWeightHigh, diffTypeLicenseRemoved)
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, diffTypeLicenseRemoved))
		return out
	}

	if !md.License.Found() && mdUp.License.Found() {
		dffs.AddDifference(md, mdUp, diffWeightHigh, diffTypeLicenseAdded)
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, diffTypeLicenseAdded))
		return out
	}

	// Identify license name and type
	hnd.licHnd.IdentifyType(&md.License)
	hnd.licHnd.IdentifyType(&mdUp.License)

	// Minor changes in the same license
	if md.License.Name == mdUp.License.Name {
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, diffTypeLicenseMinorChanges))
		dffs.AddDifference(md, mdUp, diffWeightLow, diffTypeLicenseMinorChanges)
		return out
	}

	// License name changed maintaining restrictiveness type
	if md.License.Type == mdUp.License.Type && md.License.Name != mdUp.License.Name {
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, diffTypeLicenseNameChanged))
		dffs.AddDifference(md, mdUp, diffWeightMedium, diffTypeLicenseNameChanged)
		return out
	}

	// License changed to less restrictive
	if !md.License.IsMoreRestrictive(mdUp.License.Type) {
		hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, diffTypeLicenseLessStrictChanged))
		dffs.AddDifference(md, mdUp, diffWeightMedium, diffTypeLicenseLessStrictChanged)
		return out
	}

	// License changed to more restrictive
	hnd.lgr.Debug(fmt.Sprintf("Module %s -> %s differences %d", md, mdUp, diffTypeLicenseMoreStrictChanged))
	dffs.AddDifference(md, mdUp, diffWeightMedium, diffTypeLicenseMoreStrictChanged)
	return out
}

func (hnd Handler) updateDifferencesSubModule(md, mdUp Module, dffs *Differences) output.Output {
	out := output.Create(pkg + ".updateDifferencesSubModule")
	if out := hnd.updateDir(&md); out.HasError() {
		dffs.AddModule(md, diffWeightHigh, diffTypeModuleFetchError)
		return out
	}
	if out := hnd.updateDir(&mdUp); out.HasError() {
		dffs.AddModule(mdUp, diffWeightHigh, diffTypeModuleFetchError)
		return out
	}

	upOut := hnd.updateDifferencesModule(md, mdUp, dffs)
	if upOut.HasError() {
		return out
	}

	subMds, mdsOut := hnd.listSubModules(md.Dir)
	if mdsOut.HasError() {
		return out
	}
	if len(subMds) == 0 {
		return out
	}
	subUpMds, mdsOut := hnd.listSubModules(mdUp.Dir)
	if mdsOut.HasError() {
		return out
	}

	if subMdsOut := hnd.updateDifferencesSubModules(subMds, subUpMds, dffs); subMdsOut.HasError() {
		return subMdsOut
	}
	return out
}

func (hnd Handler) updateDifferencesSubModules(subMds, subUpMds Modules, dffs *Differences) output.Output {
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
			if out := hnd.updateDir(&upMd); out.HasError() {
				dffs.AddModule(upMd, diffWeightHigh, diffTypeModuleFetchError)
				return out
			}

			upMd.License = hnd.licHnd.FindLicense(upMd.Dir)
			hnd.licHnd.IdentifyType(&upMd.License)
			dffs.AddModule(upMd, diffWeightHigh, diffTypeNewSubmodule)
		}
	}
	return out
}

// updateDir checks that module's directory is accessible and downloads it if it isn't
func (hnd Handler) updateDir(md *Module) output.Output {
	out := output.Create(fmt.Sprintf("%s.%s '%s'", pkg, "checkModuleExistence", md))
	if md.Dir == "" || !hnd.exc.FolderAccessible(md.Dir) {
		if upDirOut := hnd.downloadModule(md); upDirOut.HasError() || md.Dir == "" {
			return upDirOut
		}
	}
	return out
}
