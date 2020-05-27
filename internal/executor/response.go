package executor

type Response struct {
	StdOutput []byte
	StdError  []byte
}

func (rsp Response) HasError() bool {
	return len(rsp.StdError) > 0
}

func (rsp Response) IsEmpty() bool {
	return len(rsp.StdOutput) == 0
}

func (rsp Response) StdErrorString() string {
	return string(rsp.StdError)
}

func (rsp Response) StdOutputString() string {
	return string(rsp.StdOutput)
}
