package entity

import "strings"

const (
	// bad request
	CreateUserEmailDuplicate                  int    = 400001
	CreateUserEmailDuplicateMessage           string = "duplicate email, email '%s` is already in use"
	CreateUserUsernameDuplicate               int    = 400002
	CreateUserUsernameDuplicateMessage        string = "duplicate username, username '%s` is already in use"
	RenewPasswordWithUnchangedPassword        int    = 400003
	RenewPasswordWithUnchangedPasswordMessage string = "cannot use same password"
	DeactivatingInActiveUser                  int    = 400004
	DeactivatingInActiveUserMessage           string = "try deactivating inactive users"
	ActivatingActiveUser                      int    = 400005
	ActivatingActiveUserMessage               string = "try activating active users"
	CreateUnitNameDuplicate                   int    = 400006
	CreateUnitNameDuplicateMessage            string = "duplicate unit name, `%s` is already in use"
	ValidatorBadRequest                       int    = 400999
	ValidatorBadRequestMessage                string = "bad request, validator failed"

	// unauthorize
	AttemptLoginFailed                 int    = 401001
	AttemptLoginFailedMessage          string = "wrong email or password"
	AttemptLoginUserDeactivated        int    = 401002
	AttemptLoginUserDeactivatedMessage string = "user is deactivated"
	ValidateTokenFailed                int    = 401003
	ValidateTokenFailedMessage         string = "invalid token"
	NonGuestAttemptGuestLogin          int    = 401004
	NonGuestAttemptGuestLoginMessage   string = "cannot login as guest"
	SessionTokenNotFound               int    = 401005
	SessionTokenNotFoundMessage        string = "token session not found"

	// forbidden
	AttemptGuestLoginForbidden         int    = 403001
	AttemptGuestLoginForbiddenMessage  string = "guest login is disabled"
	GuestTokenForbidden                int    = 403002
	GuestTokenForbiddenMessage         string = "guest token is forbidden"
	MustChangePasswordForbidden        int    = 403003
	MustChangePasswordForbiddenMessage string = "the password must be changed"
	DeleteUsedUnitForbidden            int    = 403004
	DeleteUsedUnitForbiddenMessage     string = "cannot delete used units"

	// not found
	ResetPasswordRequestNotFound        int    = 404001
	ResetPasswordRequestNotFoundMessage string = "request reset password not found"
	SearchSingleUserNotFOund            int    = 404002
	SearchSingleUserNotFOundMessage     string = "user with `%s` = `%s` not found"
	UnitNameNotFound                    int    = 404003
	UnitNameNotFoundMessage             string = "unit measurement `%s` not found"

	// conflict
	UpdateUserNoChange        int    = 409001
	UpdateUserNoChangeMessage string = "no change found"

	// internal service error
	// auth repository
	AuthRepoGetUserLoginError              int = 5000101
	AuthRepoGenerateUserResetPasswordError int = 5000102
	AuthRepoValidateResetTokenError        int = 5000103
	AuthRepoDoResetPasswordError           int = 5000104
	AuthRepoDoRenewPasswordError           int = 5000105
	// user repository
	UserRepoCreateError                 int = 5000201
	UserRepoUpdateError                 int = 5000202
	UserRepoGetUserByFieldError         int = 5000203
	UserRepoGetTotalUsersByRequestError int = 5000204
	UserRepoGetUsersByRequestError      int = 5000205
	UserRepoMapUserUsernameError        int = 5000206
	// unit repository
	UnitRepoCreateError                 int = 5000301
	UnitRepoGetUnitByNameError          int = 5000302
	UnitRepoGetTotalUnitsByRequestError int = 5000303
	UnitRepoGetUnitsByRequestError      int = 5000304
	UnitRepoDeleteError                 int = 5000305
	// auth usecase
	AuthUsecaseGenerateJwtTokenError int = 5003001
	AuthUsecaseValidateJwtTokenError int = 5003002
	// mail usecase
	MailUsecaseGenerateHTMLError      int = 5003101
	MailUsecaseGeneratePlainTextError int = 5003102
	// formatter usecase
	FormatterUsecaseFormatUserError int = 5003201
	FormatterUsecaseFormatUnitError int = 5003202
	// session usecase
	SessionUsecaseTokenInvalidType  int = 5003201
	SessionUsecaseErrorInvalidType  int = 5003202
	SessionUsecaseSaveJwtTokenError int = 5003203
	SessionUsecaseSetErrorError     int = 5003204
	SessionUsecaseSetInfoError      int = 5003205
	SessionUsecaseGetInfoError      int = 5003206
	SessionUsecaseLogoutError       int = 5003207
	SessionUsecaseGetErrorError     int = 5003208
	SessionUsecaseSetPrevRouteError int = 5003209
	// base handler
	BaseHandlerGetAuthorError int = 5006001
	BaseHandlerGetTokenError  int = 5006002
	// all handler
	AllHandlerBindError int = 5009101
	// helper
	HelperGenerateTokenError          int = 5009901
	HelperEncryptPasswordError        int = 5009902
	HelperGenerateStrongPasswordError int = 5009903
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
