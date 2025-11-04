package common

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{Message: msg}
}
