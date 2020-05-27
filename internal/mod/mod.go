package mod

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/olekukonko/tablewriter"
)

const (
	surveyPageSize = 10
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

// NewModule returns the path and version of the update taking in account any Replace settings
func (m Module) NewModule() string {
	mod := m
	if m.Replace != nil && m.Replace.Update != nil {
		mod = *m.Replace
	}
	return fmt.Sprintf("%s@%s", mod.Update.Path, mod.Update.Version)
}

// String returns the path and version of current module
func (m Module) String() string {
	return fmt.Sprintf("%s@%s", m.Path, m.Version)
}

// SelectCLI allows interactively select modules to update
func (mds *Modules) SelectCLI() error {
	prompt := &survey.MultiSelect{
		Message:  "Choose which modules to update",
		Options:  mds.surveyOptions(),
		PageSize: surveyPageSize,
	}

	var chs []int
	if err := survey.AskOne(prompt, &chs); err != nil {
		return err
	}

	var sltMds Modules
	for _, n := range chs {
		sltMds = append(sltMds, (*mds)[n])
	}
	*mds = sltMds

	return nil
}

// ToTable generates a table for CLI
func (mds Modules) ToTable() string {
	var wrt bytes.Buffer
	tbl := tablewriter.NewWriter(&wrt)
	tbl.SetHeader([]string{"Module", "Version", "New Version", "Direct"})

	for _, m := range mds {
		tbl.Append([]string{
			m.Path,
			m.Version,
			m.newVersion(),
			strconv.FormatBool(!m.Indirect),
		})
	}
	tbl.Render()

	return wrt.String()
}

func (mds Modules) surveyOptions() []string {
	var ops []string
	for _, m := range mds {
		ops = append(ops, fmt.Sprintf("%s %s -> %s", m.Path, m.Version, m.newVersion()))
	}
	return ops
}

// newVersion returns the version of the update taking in account any Replace settings
func (m Module) newVersion() string {
	mod := m
	if m.Replace != nil && m.Replace.Update != nil {
		mod = *m.Replace
	}
	return mod.Update.Version
}
