package compare

import "github.com/dpcat237/go-dsu/internal/module"

type licenseChangeHandler interface {
	compareLicenses(md module.Module, mdUp module.Module, addDff addDifference)
}

type licenseComparer struct {
}

func initLicenseComparer() *licenseComparer {
	return &licenseComparer{}
}

type addDifference func(md, mdUp module.Module, dfLv module.DiffLevel, dfTp module.DiffType)
type compareLicensesFunc func(md module.Module, mdUp module.Module, addDff addDifference)

func (f compareLicensesFunc) compareLicenses(md module.Module, mdUp module.Module, addDff addDifference) {
	f(md, mdUp, addDff)
}

// License name changed maintaining restrictiveness type
func (hnd *licenseComparer) changedSameRestrictiveness(nextHnd licenseChangeHandler) licenseChangeHandler {
	return compareLicensesFunc(func(md module.Module, mdUp module.Module, addDff addDifference) {
		if md.License.Type == mdUp.License.Type && md.License.Name != mdUp.License.Name {
			addDff(md, mdUp, module.DiffWeightMedium, module.DiffTypeLicenseNameChanged)
			return
		}
		nextHnd.compareLicenses(md, mdUp, addDff)
	})
}

// License changed to more restrictive with critical restrictiveness
func (hnd *licenseComparer) criticalRestrictiveness(nextHnd licenseChangeHandler) licenseChangeHandler {
	return compareLicensesFunc(func(md module.Module, mdUp module.Module, addDff addDifference) {
		if mdUp.License.IsCritical() {
			addDff(md, mdUp, module.DiffWeightCritical, module.DiffTypeLicenseMoreStrictChanged)
			return
		}
		nextHnd.compareLicenses(md, mdUp, addDff)
	})
}

// License changed to less restrictive
func (hnd *licenseComparer) lessRestrictive(nextHnd licenseChangeHandler) licenseChangeHandler {
	return compareLicensesFunc(func(md module.Module, mdUp module.Module, addDff addDifference) {
		if !md.License.IsMoreRestrictive(mdUp.License.Type) {
			addDff(md, mdUp, module.DiffWeightLow, module.DiffTypeLicenseLessStrictChanged)
			return
		}
		nextHnd.compareLicenses(md, mdUp, addDff)
	})
}

// License added
func (hnd *licenseComparer) licenseAdded(nextHnd licenseChangeHandler) licenseChangeHandler {
	return compareLicensesFunc(func(md module.Module, mdUp module.Module, addDff addDifference) {
		if !md.License.Found() && mdUp.License.Found() {
			addDff(md, mdUp, module.DiffWeightHigh, module.DiffTypeLicenseAdded)
			return
		}
		nextHnd.compareLicenses(md, mdUp, addDff)
	})
}

// License not found
func (hnd *licenseComparer) licenseNotFound(nextHnd licenseChangeHandler) licenseChangeHandler {
	return compareLicensesFunc(func(md module.Module, mdUp module.Module, addDff addDifference) {
		if !md.License.Found() && !mdUp.License.Found() {
			addDff(md, mdUp, module.DiffWeightLow, module.DiffTypeLicenseNotFound)
			return
		}
		nextHnd.compareLicenses(md, mdUp, addDff)
	})
}

// License removed
func (hnd *licenseComparer) licenseRemoved(nextHnd licenseChangeHandler) licenseChangeHandler {
	return compareLicensesFunc(func(md module.Module, mdUp module.Module, addDff addDifference) {
		if md.License.Found() && !mdUp.License.Found() {
			addDff(md, mdUp, module.DiffWeightHigh, module.DiffTypeLicenseRemoved)
			return
		}
		nextHnd.compareLicenses(md, mdUp, addDff)
	})
}

// Minor changes in the same license
func (hnd *licenseComparer) minorChanges(nextHnd licenseChangeHandler) licenseChangeHandler {
	return compareLicensesFunc(func(md module.Module, mdUp module.Module, addDff addDifference) {
		if md.License.Name == mdUp.License.Name {
			addDff(md, mdUp, module.DiffWeightLow, module.DiffTypeLicenseMinorChanges)
			return
		}
		nextHnd.compareLicenses(md, mdUp, addDff)
	})
}

// License changed to more restrictive
func (hnd *licenseComparer) moreRestrictive() licenseChangeHandler {
	return compareLicensesFunc(func(md module.Module, mdUp module.Module, addDff addDifference) {
		addDff(md, mdUp, module.DiffWeightHigh, module.DiffTypeLicenseMoreStrictChanged)
	})
}

// Same license
func (hnd *licenseComparer) sameLicense(nextHnd licenseChangeHandler) licenseChangeHandler {
	return compareLicensesFunc(func(md module.Module, mdUp module.Module, addDff addDifference) {
		if md.License.Hash == mdUp.License.Hash {
			return
		}
		nextHnd.compareLicenses(md, mdUp, addDff)
	})
}
