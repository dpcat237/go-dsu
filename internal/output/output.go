package output

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Output collects error and/or response, and print them by specified mode.
type Output struct {
	error    error
	method   string
	response string
	status   uint16
}

//Create creates Output object
func Create(mtd string) Output {
	return Output{
		method: mtd,
	}
}

//ErrorContainsString check if Output error contains string
func (out Output) ErrorContainsString(str string) bool {
	return strings.Contains(out.error.Error(), str)
}

//GetError returns an error from Output
func (out Output) GetError() error {
	return out.error
}

//GetStatus returns an status code from Output
func (out Output) GetStatus() int {
	return int(out.status)
}

//HasError check if Output has an error
func (out Output) HasError() bool {
	return out.error != nil
}

//IsToManyRequests check if error was caused by to many requests
func (out Output) IsToManyRequests() bool {
	return out.status == http.StatusTooManyRequests
}

// String returns Output as string wrapping method and error
func (out Output) String() string {
	return fmt.Sprintf("[%s] %s", out.method, out.error)
}

//ToString returns Output as string by specified mode
func (out Output) ToString(md Mode) string {
	if !out.HasError() {
		return out.response
	}

	if md == ModeDev {
		return out.String()
	}
	return fmt.Sprintf("%s", out.error)
}

//WithError adds an error to Output and returns same Output
func (out Output) WithError(err error) Output {
	out.error = err
	return out
}

//WithErrorString adds an error from string to Output and returns same Output
func (out Output) WithErrorString(msg string) Output {
	out.error = errors.New(msg)
	return out
}

//WithResponse adds response message to Output and returns same Output
func (out Output) WithResponse(rsp string) Output {
	out.response = rsp
	return out
}

//WithStatus adds response status code to Output and returns same Output
func (out Output) WithStatus(sts int) Output {
	out.status = uint16(sts)
	return out
}
