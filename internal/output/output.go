package output

import (
	"errors"
	"fmt"
)

const (
	ModeProd Mode = iota + 1
	ModeDev
)

// Output collects error and/or response, and print them by specified mode.
type Output struct {
	error      error
	method     string
	cmdSuccess bool
	response   string
}

type Mode uint16

func Create(mtd string) Output {
	return Output{
		method: mtd,
	}
}

func (out Output) GetError() error {
	return out.error
}

func (out Output) HasError() bool {
	return out.error != nil
}

func (out Output) IsCmdSuccessful() bool {
	return out.cmdSuccess
}

func (out Output) SetCmdSuccessful(sc bool) {
	out.cmdSuccess = sc
}

func (out Output) String() string {
	return fmt.Sprintf("[%s] %s", out.method, out.error)
}

func (out Output) ToString(md Mode) string {
	if !out.HasError() {
		return out.response
	}

	if md == ModeDev {
		return out.String()
	}
	return fmt.Sprintf("%s", out.error)
}

func (out Output) WithError(err error) Output {
	out.error = err
	return out
}

func (out Output) WithErrorPrefix(msg string) Output {
	out.error = fmt.Errorf("%s: \n%s", msg, out.error)
	return out
}

func (out Output) WithErrorString(msg string) Output {
	out.error = errors.New(msg)
	return out
}

func (out Output) WithResponse(rsp string) Output {
	out.response = rsp
	return out
}

func (md Mode) IsProduction() bool {
	return md == ModeProd
}
