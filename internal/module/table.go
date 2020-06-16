package module

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

const (
	colorWhite = tableColor(iota)
	colorGreen
	colorBlue
	colorYellow
	colorRed
)

var tableHeader = []string{"Direct Module", "Version", "New Version", "Changes"}

type tableColor uint16

// ToTable generates a table for CLI with available updates
func (mds Modules) ToTable() string {
	var wrt bytes.Buffer
	tbl := tablewriter.NewWriter(&wrt)
	tbl.SetHeader(tableHeader)
	tbl.SetAutoMergeCells(true)
	tbl.SetRowLine(true)

	for _, md := range mds {
		md.addModuleRows(tbl)
	}
	tbl.Render()

	return wrt.String()
}

func (md Module) addModuleRows(tbl *tablewriter.Table) {
	dataBase := md.rowBase()
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
	}
	return tablewriter.Colors{tablewriter.Normal, cl}
}

func (md Module) differenceToString(dff Difference) string {
	var ln string
	switch dff.Type {
	case diffTypeModuleFetchError:
		ln = fmt.Sprintf("- Error fetching - %s", dff.Module)
	case diffTypeLicenseNotFound:
		ln = fmt.Sprintf("- License not found - %s", dff.Module)
	case diffTypeLicenseAdded:
		ln = fmt.Sprintf("- License %s would be added in update of %s", dff.ModuleUpdate.License.Name, dff.Module)
	case diffTypeLicenseMinorChanges:
		ln = fmt.Sprintf("- Minor changes in license %s from %s to %s", dff.ModuleUpdate.License.Name, dff.Module, dff.ModuleUpdate)
	case diffTypeLicenseNameChanged:
		ln = fmt.Sprintf("- License would change from %s in %s to %s in %s", dff.Module.License.Name, dff.Module, dff.ModuleUpdate.License.Name, dff.ModuleUpdate)
	case diffTypeLicenseLessStrictChanged:
		ln = fmt.Sprintf("- License would change to less strictive, from %s in %s to %s in %s", dff.Module.License.Name, dff.Module, dff.ModuleUpdate.License.Name, dff.ModuleUpdate)
	case diffTypeLicenseMoreStrictChanged:
		ln = fmt.Sprintf("- License would change to more strictive, from %s in %s to %s in %s", dff.Module.License.Name, dff.Module, dff.ModuleUpdate.License.Name, dff.ModuleUpdate)
	case diffTypeLicenseRemoved:
		ln = fmt.Sprintf("- License %s would be removed in %s", dff.Module.License.Name, dff.ModuleUpdate)
	case diffTypeNewSubmodule:
		if dff.Module.License.Name == "" {
			ln = fmt.Sprintf("- Would be added new submodule %s with unknown license", dff.Module)
		} else {
			ln = fmt.Sprintf("- Would be added new submodule %s with license %s", dff.Module, dff.Module.License.Name)
		}
	}
	return ln
}

func (md Module) levelToColor(lvl diffLevel) tableColor {
	cl := colorWhite
	switch lvl {
	case diffWeightLow:
		cl = colorBlue
	case diffWeightMedium:
		cl = colorYellow
	case diffWeightHigh:
		cl = colorRed
	}
	return cl
}

func (md Module) rowBase() []string {
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
