package license

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/licenseclassifier"

	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/output"
)

const confidenceThreshold = float64(0.9)

//Handler handles functions related to license identification
type Handler interface {
	FindLicense(dir string) License
	InitializeClassifier() output.Output
}

type handler struct {
	cls *licenseclassifier.License
	lgr logger.Logger
}

// InitHandler initializes handler of licenses
func InitHandler(lgr logger.Logger) *handler {
	return &handler{
		lgr: lgr,
	}
}

// FindLicense looks for a license in given directory
func (hnd handler) FindLicense(dir string) License {
	var lic License

	pth, pthOut := hnd.licensePath(dir)
	if pthOut.HasError() || pth == "" {
		hnd.lgr.Debug(fmt.Sprintf("License not found in directory %s with error %s", dir, pthOut.String()))
		return lic
	}
	lic.Path = pth

	flHash, hsOut := hnd.fileHash(pth)
	if hsOut.HasError() {
		hnd.lgr.Debug(fmt.Sprintf("Error hashing license file from directory %s with error %s", dir, hsOut.String()))
		return lic
	}
	lic.Hash = flHash
	hnd.identifyType(&lic)

	return lic
}

//InitializeClassifier Initialize licenses classifier
func (hnd *handler) InitializeClassifier() output.Output {
	out := output.Create(pkg + ".InitializeClassifier")

	cls, err := licenseclassifier.New(confidenceThreshold)
	if err != nil {
		return out.WithError(err)
	}
	hnd.cls = cls
	return out
}

func (hnd handler) fileHash(flPath string) (string, output.Output) {
	out := output.Create(pkg + ".fileHash")

	hash := md5.New()
	flRd, err := os.Open(flPath)
	if err != nil {
		return "", out.WithError(err)
	}

	if _, err := io.Copy(hash, flRd); err != nil {
		return "", out.WithError(err)
	}
	return hex.EncodeToString(hash.Sum(nil)[:16]), out
}

func (hnd handler) isLicense(flName string) bool {
	flName = strings.ToLower(flName)
	if flName == licensesBase[0] {
		return true
	}

	for _, licBs := range licensesBase {
		if strings.HasPrefix(flName, licBs) {
			return true
		}
	}

	return false
}

func (hnd handler) identifyType(lic *License) {
	if lic.Path == "" {
		hnd.lgr.Debug("Empty path during license identification")
		return
	}
	content, err := ioutil.ReadFile(lic.Path)
	if err != nil {
		hnd.lgr.Debug(fmt.Sprintf("Error reading license file from path %s with error %s", lic.Path, err.Error()))
		return
	}
	matches := hnd.cls.MultipleMatch(string(content), true)
	if len(matches) == 0 {
		hnd.lgr.Debug("Unknown license during license identification")
		return
	}

	lic.Name = matches[0].Name
	lic.Type = Type(licenseclassifier.LicenseType(lic.Name))
}

func (hnd handler) licensePath(dir string) (string, output.Output) {
	out := output.Create(pkg + ".licenseHash")

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", out.WithError(err)
	}

	for _, fl := range files {
		if !hnd.isLicense(fl.Name()) {
			continue
		}
		return dir + "/" + fl.Name(), out
	}

	return "", out
}
