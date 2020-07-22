package executor

import "strings"

const (
	escapeGoDownload = "go: downloading"
	escapeGoFinding  = "go: finding"
)

//Response contains information returned from CLI command
type Response struct {
	StdOutput []byte
	StdError  []byte
	Success   bool
}

//HasError checks if Response has an error
func (rsp Response) HasError() bool {
	return !rsp.Success && len(rsp.StdError) > 0 && !rsp.hasFalsePositive()
}

//IsEmpty checks if Response's output is empty
func (rsp Response) IsEmpty() bool {
	return len(rsp.StdOutput) == 0
}

//StdErrorString returns STD error as a string
func (rsp Response) StdErrorString() string {
	return string(rsp.StdError)
}

//StdOutputString returns STD out as a string
func (rsp Response) StdOutputString() string {
	return string(rsp.StdOutput)
}

//hasFalsePositive check that STD error doesn't have not error messages
func (rsp Response) hasFalsePositive() bool {
	return strings.Contains(string(rsp.StdError), escapeGoDownload) || strings.Contains(string(rsp.StdError), escapeGoFinding)
}
