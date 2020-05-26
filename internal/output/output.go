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
	error    error
	method   string
	response string
}

type Mode uint16

func Create(mtd string) Output {
	return Output{
		method: mtd,
	}
}

func (out Output) HasError() bool {
	return out.error != nil
}

func (out Output) ToString(md Mode) string {
	if !out.HasError() {
		return out.response
	}

	if md == ModeDev {
		return fmt.Sprintf("[%s] %s", out.method, out.error)
	}
	return fmt.Sprintf("%s", out.error)
}

func (out Output) WithError(err error) Output {
	out.error = err
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
