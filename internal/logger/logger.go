package logger

import (
	"go.uber.org/zap"

	"github.com/dpcat237/go-dsu/internal/output"
)

const pkg = "logger"

// Logger wraps zap.Logger
type Logger struct {
	*zap.Logger
}

// Init creates a new preconfigured zap.Logger
func Init(mod output.Mode) (*Logger, output.Output) {
	out := output.Create(pkg + ".Init")
	var lgr Logger
	var zapLg *zap.Logger
	var err error

	if mod.IsProduction() {
		zapLg, err = zap.NewProduction()
	} else {
		zapLg, err = zap.NewDevelopment()
	}
	if err != nil {
		return &lgr, out.WithError(err)
	}
	lgr.Logger = zapLg

	return &lgr, out
}
