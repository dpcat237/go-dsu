package module

import (
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

var tablePreviewHeader = []string{"Direct Module", "Version", "New Version", "Changes"}

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

	md.addUpdateDifferencesRows(tbl)
}

func (md Module) addUpdateDifferencesRows(tbl *tablewriter.Table) {
	dataBase := md.previewRowBase()
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

func (md Module) levelToColor(lvl DiffLevel) tableColor {
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

func (md Module) previewRowBase() []string {
	return []string{
		md.Path,
		md.Version,
		md.newVersion(),
	}
}
