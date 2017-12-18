package subcmd

type Error interface {
	error
	IsPrintDefaults() bool
}

type subcmdError struct {
	err error
}

func (err *subcmdError) Error() string {
	return err.err.Error()
}

func (err *subcmdError) IsPrintDefaults() bool {
	return true
}

func NewError(err error) *subcmdError {
	return &subcmdError{err: err}
}

var _ Error = &subcmdError{}
