package module

import "fmt"

var tablePreviewHeader = []string{"Direct Module", "Version", "New Version", "Changes"}

// GeneratePreviewTable generates a table for CLI with available updates
func (tbl Table) GeneratePreviewTable(mds Modules) string {
	tbl.printer.SetHeader(tablePreviewHeader)
	tbl.printer.SetAutoMergeCells(true)
	tbl.printer.SetRowLine(true)

	for _, md := range mds {
		tbl.addModulePreviewRows(md)
	}
	tbl.printer.Render()

	return tbl.writer.String()
}

func (tbl *Table) addModulePreviewRows(md Module) {
	dataBase := tbl.previewRowBase(md)
	if len(md.UpdateDifferences) == 0 {
		dataBase = append(dataBase, "")
		tbl.printer.Rich(dataBase, tbl.rowColors(colorGreen, colorWhite, colorWhite, colorWhite))
		return
	}

	if len(md.UpdateDifferences) == 1 {
		dff := md.UpdateDifferences[0]
		dataBase = append(dataBase, tbl.differenceToString(dff))
		cls := tbl.rowColors(tbl.levelToColor(dff.Level), colorWhite, colorWhite, tbl.levelToColor(dff.Level))
		tbl.printer.Rich(dataBase, cls)
		return
	}

	tbl.addUpdateDifferencesRows(md)
}

func (tbl *Table) addUpdateDifferencesRows(md Module) {
	dataBase := tbl.previewRowBase(md)
	var data []string
	fst := false
	hgLvl := md.UpdateDifferences.highestLevel()
	for _, dff := range md.UpdateDifferences {
		data = dataBase
		data = append(data, tbl.differenceToString(dff))
		if fst {
			cls := tbl.rowColors(tbl.levelToColor(hgLvl), colorWhite, colorWhite, tbl.levelToColor(dff.Level))
			tbl.printer.Rich(data, cls)
			continue
		}

		cls := tbl.rowColors(tbl.levelToColor(md.UpdateDifferences.highestLevel()), colorWhite, colorWhite, tbl.levelToColor(dff.Level))
		tbl.printer.Rich(data, cls)
		fst = true
	}
}

func (tbl Table) differenceToString(dff Difference) string {
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

func (tbl Table) levelToColor(lvl DiffLevel) tableColor {
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

func (tbl Table) previewRowBase(md Module) []string {
	return []string{
		md.Path,
		md.Version,
		md.newVersion(),
	}
}
