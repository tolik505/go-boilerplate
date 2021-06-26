package apperr

// Known helps to handle expected logical errors.
// E.g. a client of the API doesn't provide required data, the Known error
// is being propagated to request's handler and returned as a clear message
// with a specific error code.
type Known struct {
	code    string
	message string
}

// NewKnown creates Known error
func NewKnown(code, msg string) Known {
	return Known{
		code:    code,
		message: msg,
	}
}

// Error returns error message string
func (e Known) Error() string {
	return e.message
}

// Code returns error code
func (e Known) Code() string {
	return e.code
}
