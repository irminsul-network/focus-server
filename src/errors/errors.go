package errors

import "fmt"

type CustomError struct {
	message string
	cause   error
}

type UserError struct {
	CustomError
}

type ServerError struct {
	CustomError
}

func (e CustomError) Error() string {
	return fmt.Sprintf("Server Error: %s", e.cause)
}

func (e UserError) Error() string {
	return fmt.Sprintf("User Error: %s", e.cause)
}

func NewServerError(message string, cause error) ServerError {
	return ServerError{
		CustomError{message, cause},
	}
}
func NewUserError(message string, cause error) ServerError {
	return ServerError{
		CustomError{message, cause},
	}
}
