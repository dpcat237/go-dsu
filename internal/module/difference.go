package module

const (
	diff_weight_none = diffLevel(iota)
	diff_weight_low
	diff_weight_medium
	diff_weight_high
)

const (
	diff_type_module_fetch_error = diffType(iota)
	diff_type_license_not_found
	diff_type_license_added
	diff_type_license_minor_changes
	diff_type_license_name_changed
	diff_type_license_less_strict_changed
	diff_type_license_more_strict_changed
	diff_type_license_removed
	diff_type_new_submodule
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
