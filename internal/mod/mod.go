package mod

import (
	"bytes"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Module holds information about a specific module
type Module struct {
	Main      bool         `json:",omitempty"` // is this the main module?
	Indirect  bool         `json:",omitempty"` // module is only indirectly needed by main module
	Dir       string       `json:",omitempty"` // directory holding local copy of files, if any
	GoMod     string       `json:",omitempty"` // path to go.mod file describing module, if any
	GoVersion string       `json:",omitempty"` // go version used in module
	Path      string       `json:",omitempty"` // module path
	Version   string       `json:",omitempty"` // module version
	Versions  []string     `json:",omitempty"` // available module versions
	Error     *ModuleError `json:",omitempty"` // error loading module
	Replace   *Module      `json:",omitempty"` // replaced by this module
	Time      *time.Time   `json:",omitempty"` // time version was created
	Update    *Module      `json:",omitempty"` // available update
}

type Modules []Module

// ModuleError represents the error when a module cannot be loaded
type ModuleError struct {
	Err string
}

// ToTable generates a table for CLI
func (mds Modules) ToTable() string {
	var wrt bytes.Buffer
	tbl := tablewriter.NewWriter(&wrt)
	tbl.SetHeader([]string{"Module", "Version", "New Version", "Direct"})

	for _, m := range mds {
		tbl.Append([]string{
			m.Path,
			m.currentVersion(),
			m.newVersion(),
			strconv.FormatBool(!m.Indirect),
		})
	}
	tbl.Render()

	return wrt.String()
}

// currentVersion returns the current version of the module taking in account any Replace settings
func (m Module) currentVersion() string {
	if m.Replace != nil {
		return m.Replace.Version
	}
	return m.Version
}

// newVersion returns the version of the update taking in account any Replace settings
func (m Module) newVersion() string {
	mod := m
	if m.Replace != nil && m.Replace.Update != nil {
		mod = *m.Replace
	}
	return mod.Update.Version
}
