package module

import (
	"github.com/olekukonko/tablewriter"
)

const (
	colorWhite = tableColor(iota)
	colorGreen
	colorBlue
	colorYellow
	colorRed
	colorRedBg
)

type tableColor uint16

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

func (md Module) rowColors(clsTb ...tableColor) []tablewriter.Colors {
	var cls []tablewriter.Colors
	for _, clTb := range clsTb {
		cls = append(cls, md.cellColor(clTb))
	}
	return cls
}
