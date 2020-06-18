package download

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dpcat237/go-dsu/internal/executor"
	"github.com/dpcat237/go-dsu/internal/module"
	"github.com/dpcat237/go-dsu/internal/output"
)

const (
	cmdChmodModule = "(chmod 744 %s && chmod 655 %s/*)"
	cmdModDownload = "mod download -json"

	pkg = "download"
)

//Handler handles functions related to download modules
type Handler struct {
	exc *executor.Executor
}

//InitHandler initializes downloads handler
func InitHandler(exc *executor.Executor) *Handler {
	return &Handler{
		exc: exc,
	}
}

//DownloadModule download module and returns local directory to module
func (hnd Handler) DownloadModule(mdPth string) (string, output.Output) {
	out := output.Create(pkg + ".DownloadModule")

	if mdPth == "" {
		return "", out
	}

	return hnd.modDownload(mdPth)
}

func (hnd Handler) modDownload(mdPth string) (string, output.Output) {
	out := output.Create(fmt.Sprintf("%s.%s '%s'", pkg, "downloadModule", mdPth))

	// Download
	dwnRsp, dwnOut := hnd.exc.ExecProject(fmt.Sprintf("%s %s", cmdModDownload, mdPth))
	if dwnOut.HasError() {
		return "", dwnOut
	}

	var mdDwn module.Module
	dec := json.NewDecoder(bytes.NewReader(dwnRsp.StdOutput))
	if err := dec.Decode(&mdDwn); err != nil {
		return "", out.WithError(err)
	}
	dir := mdDwn.Dir

	// Double check permissions
	if _, prmOut := hnd.exc.ExecGlobal(fmt.Sprintf(cmdChmodModule, dir, dir)); prmOut.HasError() {
		return "", prmOut
	}

	return dir, out
}
