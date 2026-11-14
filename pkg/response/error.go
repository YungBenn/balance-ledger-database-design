package response

type Error struct {
	Status int
	Err    error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func NewError(status int, err error) *Error {
	return &Error{
		Status: status,
		Err:    err,
	}
}
