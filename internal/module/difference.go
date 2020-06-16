package module

const (
	diffWeightNone = diffLevel(iota)
	diffWeightLow
	diffWeightMedium
	diffWeightHigh
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
)

type diffLevel uint16
type diffType uint16

// Difference contains differences between module versions
type Difference struct {
	Module       Module
	ModuleUpdate Module
	Level        diffLevel
	Type         diffType
}

// Differences contains multiple differences
type Differences []Difference

// AddDifference adds difference details with module and available update
func (dffs *Differences) AddDifference(md, mdUp Module, dfLv diffLevel, dfTp diffType) {
	df := Difference{
		Module:       md,
		ModuleUpdate: mdUp,
		Level:        dfLv,
		Type:         dfTp,
	}
	*dffs = append(*dffs, df)
}

// AddModule adds difference of module
func (dffs *Differences) AddModule(md Module, dfLv diffLevel, dfTp diffType) {
	dif := Difference{
		Module: md,
		Level:  dfLv,
		Type:   dfTp,
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
