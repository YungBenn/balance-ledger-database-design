package response

import "github.com/gofiber/fiber/v2"

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

func ErrorHandler(err error) error {
	if customErr, ok := err.(*Error); ok {
		return fiber.NewError(customErr.Status, customErr.Err.Error())
	}

	return nil
}
