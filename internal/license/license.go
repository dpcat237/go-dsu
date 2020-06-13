package license

const pkg = "license"

// License types
// Thank you https://github.com/google/go-licenses
const (
	// Unknown license type.
	Unknown = Type("")
	// Restricted licenses require mandatory source distribution if we ship a
	// product that includes third-party code protected by such a license.
	Restricted = Type("restricted")
	// Reciprocal licenses allow usage of software made available under such
	// licenses freely in *unmodified* form. If the third-party source code is
	// modified in any way these modifications to the original third-party
	// source code must be made available.
	Reciprocal = Type("reciprocal")
	// Notice licenses contain few restrictions, allowing original or modified
	// third-party software to be shipped in any product without endangering or
	// encumbering our source code. All of the licenses in this category do,
	// however, have an "original Copyright notice" or "advertising clause",
	// wherein any external distributions must include the notice or clause
	// specified in the license.
	Notice = Type("notice")
	// Permissive licenses are even more lenient than a 'notice' license.
	// Not even a copyright notice is required for license compliance.
	Permissive = Type("permissive")
	// Unencumbered covers licenses that basically declare that the code is "free for any use".
	Unencumbered = Type("unencumbered")
	// Forbidden licenses are forbidden to be used.
	Forbidden = Type("FORBIDDEN")
)

var licensesBase = []string{
	"license",
	"copying",
	"copyright",
	"licence",
	"unlicense",
	"copyleft",
}

// License software license details.
type License struct {
	Hash string
	Name string
	Path string
	Type Type
}

// Type identifies a class of software license.
type Type string

// Found checks by hash if a license was found
func (lic License) Found() bool {
	return lic.Hash != ""
}

// IsMoreRestrictive defines if license if more restrictions comparing licenses types
func (lic License) IsMoreRestrictive(nTp Type) bool {
	if lic.Type == nTp {
		return false
	}

	switch {
	case lic.Type == Unencumbered && (nTp == Permissive || nTp == Notice || nTp == Reciprocal || nTp == Restricted || nTp == Forbidden):
		return true
	case lic.Type == Permissive && (nTp == Notice || nTp == Reciprocal || nTp == Restricted || nTp == Forbidden):
		return true
	case lic.Type == Notice && (nTp == Reciprocal || nTp == Restricted || nTp == Forbidden):
		return true
	case lic.Type == Reciprocal && (nTp == Restricted || nTp == Forbidden):
		return true
	case lic.Type == Restricted && nTp == Forbidden:
		return true
	}

	return false
}
