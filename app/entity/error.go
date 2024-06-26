package entity

import "strings"

const (
	// bad request
	ValidateCreateUserFailed           int    = 40001
	CreateUserEmailDuplicate           int    = 40002
	CreateUserEmailDuplicateMessage    string = "duplicate email, email %s already exists"
	CreateUserUsernameDuplicate        int    = 40003
	CreateUserUsernameDuplicateMessage string = "duplicate username, username %s already exists"
	ValidateAttemptLoginFailed         int    = 40004
	AttemptLoginFailed                 int    = 40005
	AttemptLoginIncorrectMessage       string = "incorrect email or password"
	AttemptLoginDisabledMessage        string = "user is disabled"
	MysqlDuplicateKeyError             int    = 40099
	MysqlDuplicateKeyErrorMessage      string = "duplicate key error"

	// unauthorized
	ValidateTokenUserNotFound        int    = 40101
	ValidateTokenUserNotFoundMessage string = "user not found"
	ValidateTokenUserNotMatch        int    = 40102
	ValidateTokenUserNotMatchMessage string = "user does not match with token"

	// forbidden
	ValidateTokenGuestNotAccepted        int    = 40301
	ValidateTokenGuestNotAcceptedMessage string = "this page cannot accessed by guest"

	// internal error
	CreateUserError            int    = 50001
	CreateUserErrorEmptyAuthor string = "author is nil"
	GeneratePasswordError      int    = 50002
	ValidateTokenDBError       int    = 50003
	AttemptLoginError          int    = 50003
	GetAuthorFromJwtError      int    = 50004
	SessionError               int    = 50099
	SessionErrorMessage        string = "session error"
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
