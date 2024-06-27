package entity

import "strings"

const (
	// bad request
	ValidateCreateUserFailed           int    = 40001
	CreateUserEmailDuplicate           int    = 40002
	CreateUserEmailDuplicateMessage    string = "duplicate email, email %s already exists"
	CreateUserUsernameDuplicate        int    = 40003
	CreateUserUsernameDuplicateMessage string = "duplicate username, username %s already exists"
	MysqlDuplicateKeyError             int    = 40099
	MysqlDuplicateKeyErrorMessage      string = "duplicate key error"

	// unauthorized

	// forbidden

	// internal error
	CreateUserError            int    = 50001
	CreateUserErrorEmptyAuthor string = "author is nil"
	GeneratePasswordError      int    = 50002
)

// for echo.HTTPError.Internal
type InternalError struct {
	Code    int
	Message string
}

func (e *InternalError) Error() string {
	return e.Message
}

func NewInternalError(code int, messages ...string) *InternalError {
	res := &InternalError{code, ""}
	if len(messages) > 0 {
		res.Message = strings.Join(messages, "; ")
	}
	return res
}
