package executor

import "strings"

const (
	escapeGoDownload = "go: downloading"
)

//Response contains information returned from CLI command
type Response struct {
	StdOutput []byte
	StdError  []byte
}

//HasError checks if Response has an error
func (rsp Response) HasError() bool {
	return len(rsp.StdError) > 0 && !strings.Contains(string(rsp.StdError), escapeGoDownload)
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
