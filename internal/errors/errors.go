package serviceErrors

type ServiceError struct {
	Code    int
	Message string
	err     error
}

func (e *ServiceError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}
func NewServiceErrorWithMessage(code int, message string) error {
	return &ServiceError{
		Code:    code,
		Message: message,
	}
}

func NewServiceErrorWithMessageAndError(code int, message string, err error) error {
	return &ServiceError{
		Code:    code,
		Message: message,
		err:     err,
	}
}
