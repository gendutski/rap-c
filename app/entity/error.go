package entity

import "strings"

const (
	// bad request
	UserRepoCreateEmailDuplicate                    int    = 40001
	UserRepoCreateEmailDuplicateMessage             string = "duplicate email, email %s already exists"
	UserRepoCreateUsernameDuplicate                 int    = 40002
	UserRepoCreateUsernameDuplicateMessage          string = "duplicate username, username %s already exists"
	UserUsecaseAttemptLoginDisableUser              int    = 40003
	UserUsecaseAttemptLoginDisableUserMessage       string = "user is disabled"
	UserusecaseAttemptLogigIncorrectPassword        int    = 40004
	UserusecaseAttemptLoginIncorrectPasswordMessage string = "incorrect password"
	UserUsecaseRenewPasswordUnchanged               int    = 40005
	UserUsecaseRenewPasswordUnchangedMessage        string = "you cannot use old password for new password"
	ValidatorNotValid                               int    = 40098
	MysqlDuplicateKeyError                          int    = 40099
	MysqlDuplicateKeyErrorMessage                   string = "duplicate key error"

	// unauthorized
	UserUsecaseGetUserFromJwtClaimsUnauthorized        int    = 40101
	UserUsecaseGetUserFromJwtClaimsUnauthorizedMessage string = "user does not match with token"
	UserUsecaseValidateSessionJwtTokenUnauthorized     int    = 40102

	// forbidden
	UserUsecaseAttemptGuestLoginDisabled              int    = 40301
	UserUsecaseAttemptGuestLoginDisabledMessage       string = "guest login disabled"
	UserUsecaseGetUserFromJwtClaimsForbidGuest        int    = 40302
	UserUsecaseGetUserFromJwtClaimsForbidGuestMessage string = "this page cannot accessed by guest"
	MiddlewarePasswordNotChangedSamePassword          int    = 40302
	MiddlewarePasswordNotChangedSamePasswordMessage   string = "password must be changed, cannot use same password"

	// not found
	UserRepoGetUserByFieldNotFound              int    = 40401
	UserRepoGetUserByFieldNotFoundMessage       string = "user with the %s `%v` not found"
	UserRepoValidateResetTokenNotFound          int    = 40402
	UserRepoValidateResetTokenNotFoundMessage   string = "email `%s` and token `%s` not valid"
	UserUsecaseAttemptGuestLoginNotFound        int    = 40403
	UserUsecaseAttemptGuestLoginNotFoundMessage string = "guest user not found"

	// internal error
	UserRepoCreateError                      int    = 50001
	UserRepoUpdateError                      int    = 50002
	UserRepoGetUserByFieldError              int    = 50003
	UserRepoGetTotalUsersByRequestError      int    = 50004
	UserRepoGetUsersByRequestError           int    = 50005
	UserRepoGenerateUserResetPasswordError   int    = 50006
	UserRepoValidateResetTokenError          int    = 50007
	UserUsecaseCreateError                   int    = 50008
	UserUsecaseGenerateStrongPasswordError   int    = 50009
	UserUsecaseEncryptPasswordError          int    = 50010
	UserUsecaseAttemptLoginError             int    = 50011
	UserUsecaseGenerateJwtTokenError         int    = 50012
	UserUsecaseValidateJwtTokenError         int    = 50013
	BaseHandlerGetAuthorError                int    = 50014
	MailUsecaseGeneratingEmailHTMLError      int    = 50015
	MailUsecaseGeneratingEmailPlainTextError int    = 50016
	SessionError                             int    = 50099
	SessionErrorMessage                      string = "session error"
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
