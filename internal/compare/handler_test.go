package compare

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/dpcat237/go-dsu/internal/license"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
	"github.com/dpcat237/go-dsu/internal/vulnerability"
)

type mockDownloadHandler struct {
}

func (mockDownloadHandler) CleanTemporaryData() {
}

func (mockDownloadHandler) DownloadModule(mdPth string) (string, output.Output) {
	return "", output.Output{}
}

func (mockDownloadHandler) FolderAccessible(pth string) bool {
	return false
}

func (mockDownloadHandler) GitDownload(mdPth string) (string, output.Output) {
	return "", output.Output{}
}

type mockLicenseHandler struct {
}

func (m mockLicenseHandler) FindLicense(dir string) license.License {
	if dir == "" {
		return license.License{}
	}

	var lic license.License
	lic.Path = dir
	if lic.Path[0] == 100 { // d - different
		lic.Hash = newHash()
	} else if lic.Path[0] == 101 { // e - equal
		lic.Hash = "owuem4353m4cewhdf"
	}

	if len(lic.Path) < 2 {
		return lic
	}

	if lic.Path[1] == 100 { // d - different
		lic.Name = newHash()
	} else if lic.Path[1] == 101 { // e - equal
		lic.Name = "test-license-name"
	}

	if len(lic.Path) < 3 {
		return lic
	}

	switch lic.Path[2] {
	case 97: // a
		lic.Type = license.Unknown
	case 98: // b
		lic.Type = license.Restricted
	case 99: // c
		lic.Type = license.Reciprocal
	case 100: // d
		lic.Type = license.Notice
	case 101: // e
		lic.Type = license.Permissive
	case 102: // f
		lic.Type = license.Unencumbered
	case 103: // g
		lic.Type = license.Forbidden
	}

	return lic
}

func (mockLicenseHandler) InitializeClassifier() output.Output {
	return output.Output{}
}

type mockLogger struct {
}

func (mockLogger) Sugar() *zap.SugaredLogger {
	return &zap.SugaredLogger{}
}

func (mockLogger) WithOptions(opts ...zap.Option) *zap.Logger {
	return &zap.Logger{}
}

func (mockLogger) With(fields ...zap.Field) *zap.Logger {
	return &zap.Logger{}
}

func (mockLogger) Debug(msg string, fields ...zap.Field) {
}

func (mockLogger) Info(msg string, fields ...zap.Field) {
}

func (mockLogger) Warn(msg string, fields ...zap.Field) {
}

func (mockLogger) Error(msg string, fields ...zap.Field) {
}

func (mockLogger) Fatal(msg string, fields ...zap.Field) {
}

type mockModuleHandler struct {
}

func (mockModuleHandler) ListAvailable(direct, withUpdate bool) (module.Modules, output.Output) {
	return nil, output.Output{}
}

func (mockModuleHandler) ListSubModules(pth string) (module.Modules, output.Output) {
	return nil, output.Output{}
}

type mockVulnerabilityHandler struct {
}

func (mockVulnerabilityHandler) ModuleVulnerabilities(pth string) (vulnerability.Vulnerabilities, output.Output) {
	return nil, output.Output{}
}

func (mockVulnerabilityHandler) IsSet() bool {
	return false
}

func TestHandler_addLicenseDifferences(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	type args struct {
		md   module.Module
		mdUp module.Module
		dffs *module.Differences
	}
	type want struct {
		out   output.Output
		dff   module.Difference
		isDff bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "License not found in both",
			args: args{
				md: module.Module{
					Dir:   "",
					GoMod: "testm1a",
				},
				mdUp: module.Module{
					Dir:   "",
					GoMod: "testm1b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out: output.Output{},
				dff: module.Difference{
					Level: module.DiffWeightLow,
					Module: module.Module{
						GoMod: "testm1a",
					},
					ModuleUpdate: module.Module{
						GoMod: "testm1b",
					},
					Type: module.DiffTypeLicenseNotFound,
				},
				isDff: true,
			},
		},
		{
			name: "Same license",
			args: args{
				md: module.Module{
					Dir:   "e",
					GoMod: "testm2a",
				},
				mdUp: module.Module{
					Dir:   "e",
					GoMod: "testm2b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out:   output.Output{},
				dff:   module.Difference{},
				isDff: false,
			},
		},
		{
			name: "License removed",
			args: args{
				md: module.Module{
					Dir:   "e",
					GoMod: "testm3a",
				},
				mdUp: module.Module{
					Dir:   "",
					GoMod: "testm3b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out: output.Output{},
				dff: module.Difference{
					Level: module.DiffWeightHigh,
					Module: module.Module{
						GoMod: "testm3a",
					},
					ModuleUpdate: module.Module{
						GoMod: "testm3b",
					},
					Type: module.DiffTypeLicenseRemoved,
				},
				isDff: true,
			},
		},
		{
			name: "License added",
			args: args{
				md: module.Module{
					Dir:   "",
					GoMod: "testm4a",
				},
				mdUp: module.Module{
					Dir:   "e",
					GoMod: "testm4b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out: output.Output{},
				dff: module.Difference{
					Level: module.DiffWeightHigh,
					Module: module.Module{
						GoMod: "testm4a",
					},
					ModuleUpdate: module.Module{
						GoMod: "testm4b",
					},
					Type: module.DiffTypeLicenseAdded,
				},
				isDff: true,
			},
		},
	}

	hnd := Handler{
		dwnHnd: mockDownloadHandler{},
		lgr:    mockLogger{},
		licHnd: mockLicenseHandler{},
		mdHnd:  mockModuleHandler{},
		vlnHnd: mockVulnerabilityHandler{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dffs module.Differences
			out := hnd.addLicenseDifferences(tt.args.md, tt.args.mdUp, &dffs)

			assert.Equal(t, tt.want.out.GetError(), out.GetError())
			assert.Equal(t, tt.want.isDff, len(dffs) > 0)

			if !tt.want.isDff {
				return
			}
			dff := dffs[0]

			assert.Equal(t, tt.want.dff.Level, dff.Level)
			assert.Equal(t, tt.want.dff.Type, dff.Type)
			assert.Equal(t, tt.want.dff.Module.GoMod, dff.Module.GoMod)
			assert.Equal(t, tt.want.dff.Module.GoMod, dff.Module.GoMod)
			assert.Equal(t, tt.want.dff.ModuleUpdate.GoMod, dff.ModuleUpdate.GoMod)
		})
	}
}

func TestHandler_addLicenseTypeDifferences(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	type args struct {
		md   module.Module
		mdUp module.Module
		dffs *module.Differences
	}
	type want struct {
		out   output.Output
		dff   module.Difference
		isDff bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Minor changes in license",
			args: args{
				md: module.Module{
					Dir:   "de",
					GoMod: "testm1a",
				},
				mdUp: module.Module{
					Dir:   "de",
					GoMod: "testm1b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out: output.Output{},
				dff: module.Difference{
					Level: module.DiffWeightLow,
					Module: module.Module{
						GoMod: "testm1a",
					},
					ModuleUpdate: module.Module{
						GoMod: "testm1b",
					},
					Type: module.DiffTypeLicenseMinorChanges,
				},
				isDff: true,
			},
		},
		{
			name: "License name changed",
			args: args{
				md: module.Module{
					Dir:   "ddc",
					GoMod: "testm2a",
				},
				mdUp: module.Module{
					Dir:   "ddc",
					GoMod: "testm2b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out: output.Output{},
				dff: module.Difference{
					Level: module.DiffWeightMedium,
					Module: module.Module{
						GoMod: "testm2a",
					},
					ModuleUpdate: module.Module{
						GoMod: "testm2b",
					},
					Type: module.DiffTypeLicenseNameChanged,
				},
				isDff: true,
			},
		},
		{
			name: "License changed to less strict 1",
			args: args{
				md: module.Module{
					Dir:   "ddg",
					GoMod: "testm2a",
				},
				mdUp: module.Module{
					Dir:   "ddb",
					GoMod: "testm2b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out: output.Output{},
				dff: module.Difference{
					Level: module.DiffWeightLow,
					Module: module.Module{
						GoMod: "testm2a",
					},
					ModuleUpdate: module.Module{
						GoMod: "testm2b",
					},
					Type: module.DiffTypeLicenseLessStrictChanged,
				},
				isDff: true,
			},
		},
		{
			name: "License changed to less strict 2",
			args: args{
				md: module.Module{
					Dir:   "ddb",
					GoMod: "testm3a",
				},
				mdUp: module.Module{
					Dir:   "ddc",
					GoMod: "testm3b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out: output.Output{},
				dff: module.Difference{
					Level: module.DiffWeightLow,
					Module: module.Module{
						GoMod: "testm3a",
					},
					ModuleUpdate: module.Module{
						GoMod: "testm3b",
					},
					Type: module.DiffTypeLicenseLessStrictChanged,
				},
				isDff: true,
			},
		},
		{
			name: "License changed to less strict 3",
			args: args{
				md: module.Module{
					Dir:   "ddc",
					GoMod: "testm4a",
				},
				mdUp: module.Module{
					Dir:   "ddd",
					GoMod: "testm4b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out: output.Output{},
				dff: module.Difference{
					Level: module.DiffWeightLow,
					Module: module.Module{
						GoMod: "testm4a",
					},
					ModuleUpdate: module.Module{
						GoMod: "testm4b",
					},
					Type: module.DiffTypeLicenseLessStrictChanged,
				},
				isDff: true,
			},
		},
		{
			name: "License changed to less strict 4",
			args: args{
				md: module.Module{
					Dir:   "ddc",
					GoMod: "testm5a",
				},
				mdUp: module.Module{
					Dir:   "dde",
					GoMod: "testm5b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out: output.Output{},
				dff: module.Difference{
					Level: module.DiffWeightLow,
					Module: module.Module{
						GoMod: "testm5a",
					},
					ModuleUpdate: module.Module{
						GoMod: "testm5b",
					},
					Type: module.DiffTypeLicenseLessStrictChanged,
				},
				isDff: true,
			},
		},
		{
			name: "License changed to less strict 5",
			args: args{
				md: module.Module{
					Dir:   "ddb",
					GoMod: "testm6a",
				},
				mdUp: module.Module{
					Dir:   "ddf",
					GoMod: "testm6b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out: output.Output{},
				dff: module.Difference{
					Level: module.DiffWeightLow,
					Module: module.Module{
						GoMod: "testm6a",
					},
					ModuleUpdate: module.Module{
						GoMod: "testm6b",
					},
					Type: module.DiffTypeLicenseLessStrictChanged,
				},
				isDff: true,
			},
		},
		{
			name: "License changed to more restrictive with critical restrictiveness",
			args: args{
				md: module.Module{
					Dir:   "dde",
					GoMod: "testm7a",
				},
				mdUp: module.Module{
					Dir:   "ddb",
					GoMod: "testm7b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out: output.Output{},
				dff: module.Difference{
					Level: module.DiffWeightCritical,
					Module: module.Module{
						GoMod: "testm7a",
					},
					ModuleUpdate: module.Module{
						GoMod: "testm7b",
					},
					Type: module.DiffTypeLicenseMoreStrictChanged,
				},
				isDff: true,
			},
		},
		{
			name: "License changed to more restrictive",
			args: args{
				md: module.Module{
					Dir:   "ddd",
					GoMod: "testm7a",
				},
				mdUp: module.Module{
					Dir:   "ddc",
					GoMod: "testm7b",
				},
				dffs: &module.Differences{},
			},
			want: want{
				out: output.Output{},
				dff: module.Difference{
					Level: module.DiffWeightHigh,
					Module: module.Module{
						GoMod: "testm7a",
					},
					ModuleUpdate: module.Module{
						GoMod: "testm7b",
					},
					Type: module.DiffTypeLicenseMoreStrictChanged,
				},
				isDff: true,
			},
		},
	}

	hnd := Handler{
		dwnHnd: mockDownloadHandler{},
		lgr:    mockLogger{},
		licHnd: mockLicenseHandler{},
		mdHnd:  mockModuleHandler{},
		vlnHnd: mockVulnerabilityHandler{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dffs module.Differences
			out := hnd.addLicenseDifferences(tt.args.md, tt.args.mdUp, &dffs)

			assert.Equal(t, tt.want.out.GetError(), out.GetError())
			assert.Equal(t, tt.want.isDff, len(dffs) > 0)

			if !tt.want.isDff {
				return
			}
			dff := dffs[0]

			assert.Equal(t, tt.want.dff.Level, dff.Level)
			assert.Equal(t, tt.want.dff.Type, dff.Type)
			assert.Equal(t, tt.want.dff.Module.GoMod, dff.Module.GoMod)
			assert.Equal(t, tt.want.dff.Module.GoMod, dff.Module.GoMod)
			assert.Equal(t, tt.want.dff.ModuleUpdate.GoMod, dff.ModuleUpdate.GoMod)
		})
	}
}

func newHash(n ...int) string {
	noRandomCharacters := 32
	if len(n) > 0 {
		noRandomCharacters = n[0]
	}

	randString := randomString(noRandomCharacters)

	hash := md5.New()
	_, err := hash.Write([]byte(randString))
	if err != nil {
		panic(err)
	}
	bs := hash.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString generates a random string of n length
func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}
