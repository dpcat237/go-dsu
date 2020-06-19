package download

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/logger"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	baseTmpFolder    = "go-dsu"
	cmdChmodModule   = "(chmod 744 %s && chmod 655 %s/*)"
	cmdModDownload   = "go mod download -json"
	modDownloadError = "not a known dependency"
	cmdPermissions   = "(cd %s && go list -m -mod=mod -json)"

	pkg = "download"
)

type details struct {
	tempDir string
	vcs     string
	version string
}

//Handler handles functions related to download modules
type Handler struct {
	exc *executor.Executor
	lgr *logger.Logger
}

//InitHandler initializes downloads handler
func InitHandler(exc *executor.Executor, lgr *logger.Logger) *Handler {
	return &Handler{
		exc: exc,
		lgr: lgr,
	}
}

//CleanTemporaryData removes temporary folder
func (hnd Handler) CleanTemporaryData() {
	bsPth := fmt.Sprintf("%s/%s", os.TempDir(), baseTmpFolder)
	if _, err := os.Stat(bsPth); err != nil {
		return
	}

	if err := os.RemoveAll(bsPth); err != nil {
		hnd.lgr.Debug(err.Error())
	}
}

//DownloadModule download module and returns local directory to module
func (hnd Handler) DownloadModule(mdPth string) (string, output.Output) {
	out := output.Create(pkg + ".DownloadModule")

	if mdPth == "" {
		return "", out
	}

	dir, dirOut := hnd.modDownload(mdPth)
	if dirOut.HasError() {
		if dirOut.ErrorContainsString(modDownloadError) {
			return hnd.gitDownload(mdPth)
		}
		return "", dirOut
	}
	return dir, out
}

// FolderAccessible verifies that provided folder is accessible and allow commands execution
func (hnd Handler) FolderAccessible(pth string) bool {
	if pth == "" {
		return false
	}

	if _, err := os.Stat(pth); os.IsNotExist(err) {
		return false
	}

	if _, out := hnd.exc.ExecGlobal(fmt.Sprintf(cmdPermissions, pth)); out.HasError() {
		return false
	}
	return true
}

func (hnd Handler) gitDownload(mdPth string) (string, output.Output) {
	out := output.Create(pkg + ".gitDownload")

	dt, err := hnd.transformModulePath(mdPth)
	if err != nil {
		return "", out.WithError(err)
	}

	rep, err := git.PlainClone(dt.tempDir, false, &git.CloneOptions{
		URL: dt.vcs,
	})
	if err != nil {
		return "", out.WithError(err)
	}

	if dt.version == "" {
		return dt.tempDir, out
	}

	w, err := rep.Worktree()
	if err != nil {
		return "", out.WithError(err)
	}

	refV, err := rep.Tag(dt.version)
	if err != nil {
		return "", out.WithError(err)
	}

	chkOpt := git.CheckoutOptions{
		Hash: refV.Hash(),
	}
	if err := w.Checkout(&chkOpt); err != nil {
		return "", out.WithError(err)
	}
	return dt.tempDir, out
}

func (hnd Handler) modDownload(mdPth string) (string, output.Output) {
	out := output.Create(fmt.Sprintf("%s.%s '%s'", pkg, "downloadModule", mdPth))

	// Download
	dwnRsp, dwnOut := hnd.exc.ExecGlobal(fmt.Sprintf("%s %s", cmdModDownload, mdPth))
	if dwnOut.HasError() {
		return "", dwnOut
	}

	var mdDwn module.Module
	dec := json.NewDecoder(bytes.NewReader(dwnRsp.StdOutput))
	if err := dec.Decode(&mdDwn); err != nil {
		return "", out.WithError(err)
	}

	// Double check permissions
	if _, prmOut := hnd.exc.ExecGlobal(fmt.Sprintf(cmdChmodModule, mdDwn.Dir, mdDwn.Dir)); prmOut.HasError() {
		return "", prmOut
	}

	return mdDwn.Dir, out
}

func (hnd Handler) transformModulePath(mdPth string) (details, error) {
	var dt details
	if strings.Contains(mdPth, "@") {
		prts := strings.Split(mdPth, "@")
		mdPth = prts[0]
		dt.version = prts[1]
	}

	uri, err := url.Parse("https://" + mdPth)
	if err != nil {
		return dt, err
	}

	dt.vcs = fmt.Sprintf("git@%s:%s.git", uri.Host, uri.Path)
	dt.tempDir = fmt.Sprintf("%s/%s/%s", os.TempDir(), baseTmpFolder, mdPth)

	return dt, nil
}
