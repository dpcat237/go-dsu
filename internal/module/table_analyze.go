package module

import (
	"fmt"
	"strconv"
)

var (
	tableAnalyzeHeader                = []string{"Module", "Direct", "Version", "License"}
	tableAnalyzeHeaderVulnerabilities = []string{"Module", "Direct", "Version", "License", "Vulnerabilities"}
)

//GenerateAnalyzeTable generates a table for CLI with analyze of current dependencies
func (tbl *Table) GenerateAnalyzeTable(mds Modules, vln bool) string {
	tbl.vulnerabilities = vln
	if tbl.vulnerabilities {
		tbl.printer.SetHeader(tableAnalyzeHeaderVulnerabilities)
	} else {
		tbl.printer.SetHeader(tableAnalyzeHeader)
	}

	tbl.printer.SetAutoMergeCells(true)
	tbl.printer.SetRowLine(true)

	for _, md := range mds {
		tbl.addModuleAnalyzeRows(md)
	}
	tbl.printer.Render()

	return tbl.writer.String()
}

func (tbl *Table) addModuleAnalyzeRows(md Module) {
	frsRow := []string{
		md.Path,
		strconv.FormatBool(!md.Indirect),
		md.Version,
		md.License.Name,
	}

	if !tbl.vulnerabilities {
		tbl.printer.Append(frsRow)
		return
	}

	if len(md.Vulnerabilities) == 0 {
		tbl.printer.Append(append(frsRow, ""))
		return
	}

	for _, vln := range md.Vulnerabilities {
		vlnStr := fmt.Sprintf("- %s; %s", vln.Title, vln.Reference)
		tbl.printer.Append(append(frsRow, vlnStr))
	}
}
