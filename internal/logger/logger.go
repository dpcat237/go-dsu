package logger

import (
	"go.uber.org/zap"

	"github.com/dpcat237/go-dsu/internal/output"
)

const pkg = "logger"

// Logger wraps zap.logger
type Logger interface {
	Sugar() *zap.SugaredLogger
	WithOptions(opts ...zap.Option) *zap.Logger
	With(fields ...zap.Field) *zap.Logger
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type logger struct {
	*zap.Logger
}

// Init creates a new preconfigured zap.logger
func Init(mod output.Mode) (*logger, output.Output) {
	out := output.Create(pkg + ".Init")
	var lgr logger
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
