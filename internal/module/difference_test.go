package module_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/vulnerability"
)

func TestDifferences_AddModule(t *testing.T) {
	type args struct {
		md   module.Module
		dfLv module.DiffLevel
		dfTp module.DiffType
	}
	tests := []struct {
		name string
		dffs module.Differences
		args args
	}{
		{
			name: "Check AddModule filled 1",
			dffs: module.Differences{},
			args: args{
				md: module.Module{
					Dir:  "test/dir",
					Main: false,
				},
				dfLv: module.DiffWeightLow,
				dfTp: module.DiffTypeLicenseMinorChanges,
			},
		},
		{
			name: "Check AddModule filled 2",
			dffs: module.Differences{},
			args: args{
				md: module.Module{
					Dir:  "test/directory",
					Main: true,
				},
				dfLv: module.DiffWeightHigh,
				dfTp: module.DiffTypeLicenseRemoved,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dffs.AddModule(tt.args.md, tt.args.dfLv, tt.args.dfTp)
			assert.Equal(t, tt.args.md.Dir, tt.dffs[0].Module.Dir)
			assert.Equal(t, tt.args.md.Main, tt.dffs[0].Module.Main)
			assert.Equal(t, tt.args.dfLv, tt.dffs[0].Level)
			assert.Equal(t, tt.args.dfTp, tt.dffs[0].Type)
		})
	}
}

func TestDifferences_AddModules(t *testing.T) {
	type args struct {
		md   module.Module
		mdUp module.Module
		dfLv module.DiffLevel
		dfTp module.DiffType
	}
	tests := []struct {
		name string
		dffs module.Differences
		args args
	}{
		{
			name: "Check AddModules filled",
			dffs: module.Differences{},
			args: args{
				md: module.Module{
					Dir:     "test/dir",
					Main:    false,
					Path:    "github.com/uyf1/hlag",
					Version: "v1.0.3",
				},
				mdUp: module.Module{
					Dir:     "test/dir",
					Main:    false,
					Path:    "github.com/uyf1/hlag",
					Version: "v1.0.5",
				},
				dfLv: module.DiffWeightLow,
				dfTp: module.DiffTypeLicenseMinorChanges,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dffs.AddModules(tt.args.md, tt.args.mdUp, tt.args.dfLv, tt.args.dfTp)
			assert.Equal(t, tt.args.md.Dir, tt.dffs[0].Module.Dir)
			assert.Equal(t, tt.args.md.Main, tt.dffs[0].Module.Main)
			assert.Equal(t, tt.args.md.Path, tt.dffs[0].Module.Path)
			assert.Equal(t, tt.args.md.Version, tt.dffs[0].Module.Version)
			assert.Equal(t, tt.args.mdUp.Dir, tt.dffs[0].ModuleUpdate.Dir)
			assert.Equal(t, tt.args.mdUp.Main, tt.dffs[0].ModuleUpdate.Main)
			assert.Equal(t, tt.args.mdUp.Path, tt.dffs[0].ModuleUpdate.Path)
			assert.Equal(t, tt.args.mdUp.Version, tt.dffs[0].ModuleUpdate.Version)
			assert.Equal(t, tt.args.dfLv, tt.dffs[0].Level)
			assert.Equal(t, tt.args.dfTp, tt.dffs[0].Type)
		})
	}
}

func TestDifferences_AddVulnerability(t *testing.T) {
	type args struct {
		md  module.Module
		vln vulnerability.Vulnerability
	}
	tests := []struct {
		name string
		dffs module.Differences
		args args
	}{
		{
			name: "Check AddVulnerability filled",
			dffs: module.Differences{},
			args: args{
				md: module.Module{
					Dir:     "test/dir",
					Main:    false,
					Path:    "github.com/uyf1/hlag",
					Version: "v1.0.5",
				},
				vln: vulnerability.Vulnerability{
					ID:    "45634h5k67h6767",
					Title: "testttlt",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dffs.AddVulnerability(tt.args.md, tt.args.vln)
			assert.Equal(t, tt.args.md.Dir, tt.dffs[0].Module.Dir)
			assert.Equal(t, tt.args.md.Main, tt.dffs[0].Module.Main)
			assert.Equal(t, tt.args.md.Path, tt.dffs[0].Module.Path)
			assert.Equal(t, tt.args.md.Version, tt.dffs[0].Module.Version)
			assert.Equal(t, tt.args.vln.ID, tt.dffs[0].Vulnerability.ID)
			assert.Equal(t, tt.args.vln.Title, tt.dffs[0].Vulnerability.Title)
		})
	}
}
