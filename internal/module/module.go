package module

import (
	"fmt"
	"strings"
	"time"

	"github.com/dpcat237/go-dsu/internal/license"
)

const (
	pkg = "mod"

	surveyPageSize = 10
)

//Details contains additional information about module
type Details struct {
	License           license.License
	UpdateDifferences Differences
}

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
	Details
}

//Modules contains collections of Modules
type Modules []Module

// ModuleError represents the error when a module cannot be loaded
type ModuleError struct {
	Err string
}

// NewModule returns the path and version of the update taking in account any Replace settings
func (md Module) NewModule() string {
	mod := md
	if md.Replace != nil && md.Replace.Update != nil {
		mod = *md.Replace
	}
	return fmt.Sprintf("%s@%s", mod.Update.Path, mod.Update.Version)
}

// HasUpdate checks if module has an update
func (md Module) HasUpdate() bool {
	return md.Update != nil
}

// String returns the path and version of current module
func (md Module) String() string {
	return fmt.Sprintf("%s@%s", md.Path, md.Version)
}

//PathCleaned cleans module path
func (md Module) PathCleaned() string {
	if strings.Contains(md.Path, ".v") {
		parts := strings.Split(md.Path, ".v")
		return parts[0]
	}
	return md.Path
}

// newVersion returns the version of the update taking in account any Replace settings
func (md Module) newVersion() string {
	mod := md
	if md.Replace != nil && md.Replace.Update != nil {
		mod = *md.Replace
	}
	return mod.Update.Version
}
