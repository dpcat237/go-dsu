package module

import (
	"github.com/dpcat237/go-dsu/internal/vulnerability"
)

const (
	diffWeightLow = diffLevel(iota)
	diffWeightMedium
	diffWeightHigh
	diffWeightCritical
)

const (
	diffTypeModuleFetchError = diffType(iota)
	diffTypeLicenseNotFound
	diffTypeLicenseAdded
	diffTypeLicenseMinorChanges
	diffTypeLicenseNameChanged
	diffTypeLicenseLessStrictChanged
	diffTypeLicenseMoreStrictChanged
	diffTypeLicenseRemoved
	diffTypeNewSubmodule
	diffTypeNewVulnerability
)

type diffLevel uint16
type diffType uint16

// Difference contains differences between module versions
type Difference struct {
	Level         diffLevel
	Module        Module
	ModuleUpdate  Module
	Type          diffType
	Vulnerability vulnerability.Vulnerability
}

// Differences contains multiple differences
type Differences []Difference

// AddModule adds difference of module
func (dffs *Differences) AddModule(md Module, dfLv diffLevel, dfTp diffType) {
	dif := Difference{
		Module: md,
		Level:  dfLv,
		Type:   dfTp,
	}
	*dffs = append(*dffs, dif)
}

// AddModules adds difference details with module and available update
func (dffs *Differences) AddModules(md, mdUp Module, dfLv diffLevel, dfTp diffType) {
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
		Type:          diffTypeNewVulnerability,
		Vulnerability: vln,
	}
	*dffs = append(*dffs, dif)
}

func (dffs Differences) highestLevel() diffLevel {
	var lvl diffLevel
	for _, dff := range dffs {
		if dff.Level > lvl {
			lvl = dff.Level
		}
	}
	return lvl
}

func (dffs Differences) vulnerabilityLevel(vln vulnerability.Vulnerability) diffLevel {
	switch vln.Severity() {
	case vulnerability.SeverityLow:
		return diffWeightLow
	case vulnerability.SeverityMedium:
		return diffWeightMedium
	case vulnerability.SeverityHigh:
		return diffWeightHigh
	case vulnerability.SeverityCritical:
		return diffWeightCritical
	}
	return diffWeightCritical
}
