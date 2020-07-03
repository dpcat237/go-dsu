package module

import "github.com/dpcat237/go-dsu/internal/vulnerability"

var tableAnalyzeHeader = []string{"Direct Module", "Submodule", "Version", "License", "Vulnerabilities"}

//GenerateAnalyzeTable generates a table for CLI with analyze of current dependencies
func (tbl *Table) GenerateAnalyzeTable(mds Modules) string {
	tbl.printer.SetHeader(tableAnalyzeHeader)
	tbl.printer.SetAutoMergeCells(true)
	tbl.printer.SetRowLine(true)

	for _, md := range mds {
		tbl.addModuleAnalyzeRows(md)
	}
	tbl.printer.Render()

	return tbl.writer.String()
}

func (tbl *Table) addSubmoduleAnalyzeRow(md Module, baseClm []string, baseSbt vulnerability.Severity) {
	bsRow := append(baseClm, md.Path, md.Version, md.License.Name)
	bsTlt := tbl.severityToColor(baseSbt)
	mdTlt := tbl.severityToColor(md.Vulnerabilities.HighestSeverity())
	tbl.printer.Rich(append(bsRow, ""), tbl.rowColors(bsTlt, mdTlt, colorWhite, colorWhite, colorWhite))

	if len(md.Vulnerabilities) == 0 {
		return
	}

	for _, vln := range md.Vulnerabilities {
		tbl.printer.Rich(append(bsRow, ""), tbl.rowColors(bsTlt, mdTlt, colorWhite, colorWhite, tbl.severityToColor(vln.Severity())))
	}
}

func (tbl *Table) addModuleAnalyzeRows(md Module) {
	frsRow := []string{
		md.Path,
		"",
		md.Version,
		md.License.Name,
		"",
	}
	if len(md.Dependencies) == 0 {
		tbl.printer.Append(frsRow)
		return
	}

	baseClm := []string{md.Path}
	md.DependenciesMap = make(map[string]Module)
	md.mapDependencies(md.Dependencies)
	mdSvt := tbl.dependenciesMapHighestSeverity(md)
	tbl.printer.Rich(frsRow, tbl.rowColors(tbl.severityToColor(mdSvt), colorWhite, colorWhite, colorWhite, colorWhite))
	for _, subMd := range md.Dependencies {
		tbl.addSubmoduleAnalyzeRow(subMd, baseClm, mdSvt)
	}
}

func (tbl *Table) dependencyHighestSeverity(md Module, svt vulnerability.Severity) vulnerability.Severity {
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

func (tbl *Table) dependenciesMapHighestSeverity(md Module) vulnerability.Severity {
	var svt vulnerability.Severity
	for _, subMd := range md.DependenciesMap {
		subMdSvt := tbl.dependencyHighestSeverity(subMd, svt)
		if subMdSvt > svt {
			svt = subMdSvt
		}
		if svt == vulnerability.SeverityCritical {
			return svt
		}
	}
	return svt
}

func (tbl Table) severityToColor(svt vulnerability.Severity) tableColor {
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
