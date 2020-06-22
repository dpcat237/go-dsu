package module

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"

	"github.com/dpcat237/go-dsu/internal/vulnerability"
)

const (
	colorWhite = tableColor(iota)
	colorGreen
	colorBlue
	colorYellow
	colorRed
	colorRedBg
)

var tableAnalyzeHeader = []string{"Direct Module", "Submodules", "Version", "License", "Vulnerabilities"}
var tablePreviewHeader = []string{"Direct Module", "Version", "New Version", "Changes"}

type tableColor uint16

//ToAnalyzeTable generates a table for CLI with analyze of current dependencies
func (mds Modules) ToAnalyzeTable() string {
	var wrt bytes.Buffer
	tbl := tablewriter.NewWriter(&wrt)
	tbl.SetHeader(tableAnalyzeHeader)
	tbl.SetAutoMergeCells(true)
	tbl.SetRowLine(true)

	for _, md := range mds {
		md.addModuleAnalyzeRows(tbl)
	}
	tbl.Render()

	return wrt.String()
}

// ToPreviewTable generates a table for CLI with available updates
func (mds Modules) ToPreviewTable() string {
	var wrt bytes.Buffer
	tbl := tablewriter.NewWriter(&wrt)
	tbl.SetHeader(tablePreviewHeader)
	tbl.SetAutoMergeCells(true)
	tbl.SetRowLine(true)

	for _, md := range mds {
		md.addModulePreviewRows(tbl)
	}
	tbl.Render()

	return wrt.String()
}

func (md Module) addSubmoduleAnalyzeRow(tbl *tablewriter.Table, baseClm []string, baseSbt vulnerability.Severity) {
	bsRow := append(baseClm, md.Path, md.Version, md.License.Name)
	bsTlt := md.severityToColor(baseSbt)
	mdTlt := md.severityToColor(md.Vulnerabilities.HighestSeverity())
	tbl.Rich(append(bsRow, ""), md.rowColors(bsTlt, mdTlt, colorWhite, colorWhite, colorWhite))

	if len(md.Vulnerabilities) == 0 {
		return
	}

	for _, vln := range md.Vulnerabilities {
		tbl.Rich(append(bsRow, ""), md.rowColors(bsTlt, mdTlt, colorWhite, colorWhite, md.severityToColor(vln.Severity())))
	}
}

func (md Module) addModuleAnalyzeRows(tbl *tablewriter.Table) {
	frsRow := []string{
		md.Path,
		"",
		md.Version,
		md.License.Name,
		"",
	}
	if len(md.Dependencies) == 0 {
		tbl.Append(frsRow)
		return
	}

	baseClm := []string{md.Path}
	md.DependenciesMap = make(map[string]Module)
	md.mapDependencies(md.Dependencies)
	mdSvt := md.dependenciesMapHighestSeverity()
	tbl.Rich(frsRow, md.rowColors(md.severityToColor(mdSvt), colorWhite, colorWhite, colorWhite, colorWhite))
	for _, subMd := range md.Dependencies {
		subMd.addSubmoduleAnalyzeRow(tbl, baseClm, mdSvt)
	}
}

func (md Module) addModulePreviewRows(tbl *tablewriter.Table) {
	dataBase := md.previewRowBase()
	if len(md.UpdateDifferences) == 0 {
		dataBase = append(dataBase, "")
		tbl.Rich(dataBase, md.rowColors(colorGreen, colorWhite, colorWhite, colorWhite))
		return
	}

	if len(md.UpdateDifferences) == 1 {
		dff := md.UpdateDifferences[0]
		dataBase = append(dataBase, md.differenceToString(dff))
		cls := md.rowColors(md.levelToColor(dff.Level), colorWhite, colorWhite, md.levelToColor(dff.Level))
		tbl.Rich(dataBase, cls)
		return
	}

	var data []string
	fst := false
	hgLvl := md.UpdateDifferences.highestLevel()
	for _, dff := range md.UpdateDifferences {
		data = dataBase
		data = append(data, md.differenceToString(dff))
		if fst {
			cls := md.rowColors(md.levelToColor(hgLvl), colorWhite, colorWhite, md.levelToColor(dff.Level))
			tbl.Rich(data, cls)
			continue
		}

		cls := md.rowColors(md.levelToColor(md.UpdateDifferences.highestLevel()), colorWhite, colorWhite, md.levelToColor(dff.Level))
		tbl.Rich(data, cls)
		fst = true
	}
}

func (md Module) cellColor(clTp tableColor) tablewriter.Colors {
	cl := tablewriter.FgWhiteColor
	switch clTp {
	case colorWhite:
		cl = tablewriter.FgWhiteColor
	case colorGreen:
		cl = tablewriter.FgGreenColor
	case colorBlue:
		cl = tablewriter.FgBlueColor
	case colorYellow:
		cl = tablewriter.FgYellowColor
	case colorRed:
		cl = tablewriter.FgHiRedColor
	case colorRedBg:
		cl = tablewriter.BgRedColor
	}
	return tablewriter.Colors{tablewriter.Normal, cl}
}

func (md Module) dependencyHighestSeverity(svt vulnerability.Severity) vulnerability.Severity {
	if svt == vulnerability.SeverityCritical {
		return svt
	}

	for _, vln := range md.Vulnerabilities {
		if vln.Severity() > svt {
			svt = vln.Severity()
		}
		if svt == vulnerability.SeverityCritical {
			return svt
		}
	}

	return svt
}

func (md Module) dependenciesMapHighestSeverity() vulnerability.Severity {
	var svt vulnerability.Severity
	for _, subMd := range md.DependenciesMap {
		subMdSvt := subMd.dependencyHighestSeverity(svt)
		if subMdSvt > svt {
			svt = subMdSvt
		}
		if svt == vulnerability.SeverityCritical {
			return svt
		}
	}
	return svt
}

func (md Module) differenceToString(dff Difference) string {
	var ln string
	switch dff.Type {
	case DiffTypeModuleFetchError:
		ln = fmt.Sprintf("- Error fetching - %s", dff.Module)
	case DiffTypeLicenseNotFound:
		ln = fmt.Sprintf("- License not found - %s", dff.Module)
	case DiffTypeLicenseAdded:
		ln = fmt.Sprintf("- License %s would be added in update of %s", dff.ModuleUpdate.License.Name, dff.Module)
	case DiffTypeLicenseMinorChanges:
		ln = fmt.Sprintf("- Minor changes in license %s from %s to %s", dff.ModuleUpdate.License.Name, dff.Module, dff.ModuleUpdate)
	case DiffTypeLicenseNameChanged:
		ln = fmt.Sprintf("- License would change from %s in %s to %s in %s", dff.Module.License.Name, dff.Module, dff.ModuleUpdate.License.Name, dff.ModuleUpdate)
	case DiffTypeLicenseLessStrictChanged:
		ln = fmt.Sprintf("- License would change to less strictive, from %s in %s to %s in %s", dff.Module.License.Name, dff.Module, dff.ModuleUpdate.License.Name, dff.ModuleUpdate)
	case DiffTypeLicenseMoreStrictChanged:
		ln = fmt.Sprintf("- License would change to more strictive, from %s in %s to %s in %s", dff.Module.License.Name, dff.Module, dff.ModuleUpdate.License.Name, dff.ModuleUpdate)
	case DiffTypeLicenseRemoved:
		ln = fmt.Sprintf("- License %s would be removed in %s", dff.Module.License.Name, dff.ModuleUpdate)
	case DiffTypeNewSubmodule:
		if dff.Module.License.Name == "" {
			ln = fmt.Sprintf("- Would be added new indirect module %s with unknown license", dff.Module)
		} else {
			ln = fmt.Sprintf("- Would be added new indirect module %s with license %s", dff.Module, dff.Module.License.Name)
		}
	case DiffTypeNewVulnerability:
		ln = fmt.Sprintf("- Update of module %s has vulnerability %s, more info %s", dff.Module.String(), dff.Vulnerability.Title, dff.Vulnerability.Reference)
	}
	return ln
}

func (md Module) levelToColor(lvl diffLevel) tableColor {
	cl := colorWhite
	switch lvl {
	case DiffWeightLow:
		cl = colorBlue
	case DiffWeightMedium:
		cl = colorYellow
	case DiffWeightHigh:
		cl = colorRed
	case DiffWeightCritical:
		cl = colorRedBg
	}
	return cl
}

func (md *Module) mapDependency(subMd Module) {
	if _, ok := md.DependenciesMap[subMd.String()]; !ok {
		md.DependenciesMap[subMd.String()] = subMd
	}
}

func (md *Module) mapDependencies(subMds []Module) {
	if len(subMds) == 0 {
		return
	}
	for _, subMd := range subMds {
		md.mapDependency(subMd)
		md.mapDependencies(subMd.Dependencies)
	}
}

func (md Module) previewRowBase() []string {
	return []string{
		md.Path,
		md.Version,
		md.newVersion(),
	}
}

func (md Module) rowColors(clsTb ...tableColor) []tablewriter.Colors {
	var cls []tablewriter.Colors
	for _, clTb := range clsTb {
		cls = append(cls, md.cellColor(clTb))
	}
	return cls
}

func (md Module) severityToColor(svt vulnerability.Severity) tableColor {
	cl := colorWhite
	switch svt {
	case vulnerability.SeverityNone:
		cl = colorWhite
	case vulnerability.SeverityLow:
		cl = colorBlue
	case vulnerability.SeverityMedium:
		cl = colorYellow
	case vulnerability.SeverityHigh:
		cl = colorRed
	case vulnerability.SeverityCritical:
		cl = colorRedBg
	}
	return cl
}
