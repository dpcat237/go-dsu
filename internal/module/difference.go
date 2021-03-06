package module

import (
	"github.com/dpcat237/go-dsu/internal/vulnerability"
)

const (
	//DiffWeightLow difference with low severity
	DiffWeightLow = DiffLevel(iota)
	//DiffWeightMedium difference with medium severity
	DiffWeightMedium
	//DiffWeightHigh difference with high severity
	DiffWeightHigh
	//DiffWeightCritical difference with critical severity
	DiffWeightCritical
)

const (
	//DiffTypeModuleFetchError error during module fetch
	DiffTypeModuleFetchError = DiffType(iota)
	//DiffTypeLicenseNotFound license not found
	DiffTypeLicenseNotFound
	//DiffTypeLicenseAdded license added
	DiffTypeLicenseAdded
	//DiffTypeLicenseMinorChanges minor changes in license
	DiffTypeLicenseMinorChanges
	//DiffTypeLicenseNameChanged changed license name
	DiffTypeLicenseNameChanged
	//DiffTypeLicenseLessStrictChanged changed license to less strict
	DiffTypeLicenseLessStrictChanged
	//DiffTypeLicenseMoreStrictChanged  changed license to more strict
	DiffTypeLicenseMoreStrictChanged
	//DiffTypeLicenseRemoved license removed
	DiffTypeLicenseRemoved
	//DiffTypeNewSubmodule new submodule
	DiffTypeNewSubmodule
	//DiffTypeNewVulnerability new vulnerability
	DiffTypeNewVulnerability
)

type DiffLevel uint16
type DiffType uint16

// Difference contains differences between module versions
type Difference struct {
	Level         DiffLevel
	Module        Module
	ModuleUpdate  Module
	Type          DiffType
	Vulnerability vulnerability.Vulnerability
}

// Differences contains multiple differences
type Differences []Difference

// AddModule adds difference of module
func (dffs *Differences) AddModule(md Module, dfLv DiffLevel, dfTp DiffType) {
	dif := Difference{
		Module: md,
		Level:  dfLv,
		Type:   dfTp,
	}
	*dffs = append(*dffs, dif)
}

// AddModules adds difference details with module and available update
func (dffs *Differences) AddModules(md, mdUp Module, dfLv DiffLevel, dfTp DiffType) {
	df := Difference{
		Module:       md,
		ModuleUpdate: mdUp,
		Level:        dfLv,
		Type:         dfTp,
	}
	*dffs = append(*dffs, df)
}

// AddVulnerability adds difference of vulnerability
func (dffs *Differences) AddVulnerability(md Module, vln vulnerability.Vulnerability) {
	dif := Difference{
		Level:         dffs.vulnerabilityLevel(vln),
		Module:        md,
		Type:          DiffTypeNewVulnerability,
		Vulnerability: vln,
	}
	*dffs = append(*dffs, dif)
}

func (dffs Differences) highestLevel() DiffLevel {
	var lvl DiffLevel
	for _, dff := range dffs {
		if dff.Level > lvl {
			lvl = dff.Level
		}
	}
	return lvl
}

func (dffs Differences) vulnerabilityLevel(vln vulnerability.Vulnerability) DiffLevel {
	switch vln.Severity() {
	case vulnerability.SeverityLow:
		return DiffWeightLow
	case vulnerability.SeverityMedium:
		return DiffWeightMedium
	case vulnerability.SeverityHigh:
		return DiffWeightHigh
	case vulnerability.SeverityCritical:
		return DiffWeightCritical
	}
	return DiffWeightCritical
}
