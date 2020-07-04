package module

import (
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/stretchr/testify/assert"
)

func TestNewTable(t *testing.T) {
	tests := []struct {
		name string
		want Table
	}{
		{
			name: "Test table created printer and writer initialized",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := NewTable()
			assert.NotNil(t, tbl.writer)
			assert.NotNil(t, tbl.printer)
		})
	}
}

func TestTable_cellColor(t *testing.T) {
	type args struct {
		color tableColor
	}
	tests := []struct {
		args args
		want int
	}{
		{
			args: args{color: colorWhite},
			want: tablewriter.FgWhiteColor,
		},
		{
			args: args{color: colorGreen},
			want: tablewriter.FgGreenColor,
		},
		{
			args: args{color: colorBlue},
			want: tablewriter.FgBlueColor,
		},
		{
			args: args{color: colorYellow},
			want: tablewriter.FgYellowColor,
		},
		{
			args: args{color: colorRed},
			want: tablewriter.FgHiRedColor,
		},
		{
			args: args{color: colorRedBg},
			want: tablewriter.BgRedColor,
		},
	}

	var tbl Table
	for _, tt := range tests {
		assert.Equal(t, tt.want, tbl.cellColor(tt.args.color)[1])
	}
}

func TestTable_rowColors(t *testing.T) {
	type args struct {
		colors []tableColor
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test 1 color",
			args: args{
				colors: []tableColor{colorWhite},
			},
			want: 1,
		},
		{
			name: "Test 3 colors",
			args: args{
				colors: []tableColor{colorWhite, colorYellow, colorRedBg},
			},
			want: 3,
		},
	}

	var tbl Table
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, len(tbl.rowColors(tt.args.colors...)))
		})
	}
}
