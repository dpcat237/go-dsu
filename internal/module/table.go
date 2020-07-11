package module

import (
	"bytes"

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

type Table struct {
	printer         *tablewriter.Table
	vulnerabilities bool
	writer          *bytes.Buffer
}

func NewTable() Table {
	var tbl Table
	tbl.writer = &bytes.Buffer{}
	tbl.printer = tablewriter.NewWriter(tbl.writer)
	return tbl
}

func (tbl Table) cellColor(clTp tableColor) tablewriter.Colors {
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

func (tbl Table) rowColors(clsTb ...tableColor) []tablewriter.Colors {
	var cls []tablewriter.Colors
	for _, clTb := range clsTb {
		cls = append(cls, tbl.cellColor(clTb))
	}
	return cls
}
