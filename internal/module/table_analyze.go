package module

import (
	"bytes"

	"github.com/olekukonko/tablewriter"

	"github.com/dpcat237/go-dsu/internal/vulnerability"
)

var tableAnalyzeHeader = []string{"Direct Module", "Submodule", "Version", "License", "Vulnerabilities"}

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
