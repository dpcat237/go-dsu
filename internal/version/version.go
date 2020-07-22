package version

import (
	"runtime"
	"strings"

	ver "github.com/hashicorp/go-version"
)

func IsModSupported() bool {
	min, err := ver.NewVersion("1.14")
	if err != nil {
		return false
	}
	cur, err := ver.NewVersion(current())
	if err != nil {
		return false
	}

	return cur.GreaterThanOrEqual(min)
}

func current() string {
	return strings.ReplaceAll(runtime.Version(), "go", "")
}
