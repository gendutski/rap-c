package authusecase

import (
	"context"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"
	"rap-c/app/repository/contract"
	usecasecontract "rap-c/app/usecase/contract"
	"rap-c/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

const (
	tokenStrID    string = "id"
	tokenStrName  string = "userName"
	tokenStrEmail string = "email"
)

func NewUsecase(cfg config.Config, userRepo contract.UserRepository) usecasecontract.AuthUsecase {
	return &usecase{cfg, userRepo}
}

type usecase struct {
	cfg      config.Config
	userRepo contract.UserRepository
}

func (uc *usecase) AttemptLogin(ctx context.Context, payload *entity.AttemptLoginPayload) (*entity.User, error) {
	// validate payload
	validate := helper.GenerateStructValidator()
	errMessages := payload.Validate(validate)
	if len(errMessages) > 0 {
		return nil, &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  errMessages,
			Internal: entity.NewInternalError(entity.ValidatorNotValid, errMessages...),
		}
	}

	// get user by email
	user, err := uc.userRepo.GetUserByField(ctx, "email", payload.Email, http.StatusBadRequest)
	if err != nil {
		return nil, err
	}

	// check user status
	if user.Disabled {
		return nil, &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  entity.UserUsecaseAttemptLoginDisableUserMessage,
			Internal: entity.NewInternalError(entity.UserUsecaseAttemptLoginDisableUser, entity.UserUsecaseAttemptLoginDisableUserMessage),
		}
	}

	// validate password
	if !helper.ValidateEncryptedPassword(user.Password, payload.Password) {
		return nil, &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  entity.UserusecaseAttemptLoginIncorrectPasswordMessage,
			Internal: entity.NewInternalError(entity.UserusecaseAttemptLogigIncorrectPassword, entity.UserusecaseAttemptLoginIncorrectPasswordMessage),
		}
	}

	return user, nil
}

func (uc *usecase) AttemptGuestLogin(ctx context.Context) (*entity.User, error) {
	if !uc.cfg.EnableGuestLogin {
		return nil, &echo.HTTPError{
			Code:     http.StatusForbidden,
			Message:  entity.UserUsecaseAttemptGuestLoginDisabledMessage,
			Internal: entity.NewInternalError(entity.UserUsecaseAttemptGuestLoginDisabled, entity.UserUsecaseAttemptGuestLoginDisabledMessage),
		}
	}
	users, err := uc.userRepo.GetUsersByRequest(ctx, &entity.GetUserListRequest{GuestOnly: true, Page: 1})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  entity.UserUsecaseAttemptGuestLoginNotFoundMessage,
			Internal: entity.NewInternalError(entity.UserUsecaseAttemptGuestLoginNotFound, entity.UserUsecaseAttemptGuestLoginNotFoundMessage),
		}
	}
	return users[0], nil
}

func (uc *usecase) GenerateJwtToken(ctx context.Context, user *entity.User, isLongSession bool) (string, error) {
	exp := time.Now().Add(time.Minute * time.Duration(uc.cfg.JwtExpirationInMinutes)).Unix()
	if isLongSession {
		// long session for remember login session
		exp = time.Now().Add(time.Hour * 24 * time.Duration(uc.cfg.JwtRememberInDays)).Unix()
	}
	claims := jwt.MapClaims{
		tokenStrID:    user.ID,
		tokenStrName:  user.Username,
		tokenStrEmail: user.Email,
		"exp":         exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(uc.cfg.JwtSecret))
	if err != nil {
		return "", &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserUsecaseGenerateJwtTokenError, err.Error()),
		}
	}
	return tokenStr, nil
}

func (uc *usecase) ValidateJwtToken(ctx context.Context, token *jwt.Token, guestAccepted bool) (*entity.User, error) {
	// get claims
	claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
	if !ok {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserUsecaseValidateJwtTokenError, "failed to cast claims as jwt.MapClaims"),
		}
	}
	// get user from claims
	return uc.getUserFromJwtClaims(ctx, claims, guestAccepted)
}

func (uc *usecase) ValidateSessionJwtToken(ctx context.Context, r *http.Request, w http.ResponseWriter, store sessions.Store, guestAccepted bool) (*entity.User, string, error) {
	// get token from session
	sess := entity.InitSession(r, w, store, entity.SessionID, uc.cfg.LogMode, uc.cfg.EnableWarnFileLog)
	tokenStr, ok := sess.Get(entity.TokenSessionName).(string)
	if !ok {
		message := "token not found in session"
		return nil, "", &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  message,
			Internal: entity.NewInternalError(entity.UserUsecaseValidateSessionJwtTokenUnauthorized, message),
		}
	}

	// parse token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(uc.cfg.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, "", &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  "parse token failed",
			Internal: entity.NewInternalError(entity.UserUsecaseValidateSessionJwtTokenUnauthorized, err.Error()),
		}
	}

	// get user from claims
	user, err := uc.getUserFromJwtClaims(ctx, claims, guestAccepted)
	if err != nil {
		return nil, "", err
	}
	return user, tokenStr, nil
}

func (uc *usecase) RenewPassword(ctx context.Context, user *entity.User, payload *entity.RenewPasswordPayload) error {
	// validate payload
	validate := helper.GenerateStructValidator()
	errMessages := payload.Validate(validate)
	if len(errMessages) > 0 {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  errMessages,
			Internal: entity.NewInternalError(entity.ValidatorNotValid, errMessages...),
		}
	}

	// check whether the new password is the same as the old password
	if helper.ValidateEncryptedPassword(user.Password, payload.Password) {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  entity.UserUsecaseRenewPasswordUnchangedMessage,
			Internal: entity.NewInternalError(entity.UserUsecaseRenewPasswordUnchanged, entity.UserUsecaseRenewPasswordUnchangedMessage),
		}
	}

	// set new password
	_, encryptedPass, err := uc.generateUserPassword(payload.Password)
	if err != nil {
		return err
	}
	user.Password = encryptedPass
	user.PasswordMustChange = false
	user.UpdatedBy = user.Username

	// save
	return uc.userRepo.Update(ctx, user)
}

func (uc *usecase) RequestResetPassword(ctx context.Context, email string) (*entity.User, *entity.PasswordResetToken, error) {
	if email == "" {
		return nil, nil, echo.NewHTTPError(http.StatusBadRequest)
	}

	// check email
	user, err := uc.userRepo.GetUserByField(ctx, "email", email, http.StatusBadRequest)
	if err != nil {
		return nil, nil, err
	}

	// generate reset password token
	token, err := uc.userRepo.GenerateUserResetPassword(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}

func (uc *usecase) ValidateResetPassword(ctx context.Context, email string, token string) error {
	if email == "" || token == "" {
		var errMessages []string
		if email == "" {
			errMessages = append(errMessages, "email is required")
		}
		if token == "" {
			errMessages = append(errMessages, "token is required")
		}

		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  errMessages,
			Internal: entity.NewInternalError(entity.ValidatorNotValid, errMessages...),
		}
	}
	_, err := uc.userRepo.ValidateResetToken(ctx, email, token)
	if err != nil {
		return err
	}
	return nil
}

func (uc *usecase) SubmitResetPassword(ctx context.Context, payload *entity.ResetPasswordPayload) (*entity.User, error) {
	// validate payload
	validate := helper.GenerateStructValidator()
	errMessages := payload.Validate(validate)
	if len(errMessages) > 0 {
		return nil, &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  errMessages,
			Internal: entity.NewInternalError(entity.ValidatorNotValid, errMessages...),
		}
	}

	// get reset pasword data
	reset, err := uc.userRepo.ValidateResetToken(ctx, payload.Email, payload.Token)
	if err != nil {
		return nil, err
	}

	// get user
	user, err := uc.userRepo.GetUserByField(ctx, "email", payload.Email, http.StatusNotFound)
	if err != nil {
		return nil, err
	}

	// encrypt password
	_, encryptPass, err := uc.generateUserPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	// update user
	user.Password = encryptPass
	user.UpdatedBy = user.Username

	// update reset token
	reset.Token = ""

	// save
	err = uc.userRepo.ResetPassword(ctx, user, reset)
	if err != nil {
		return nil, err
	}
	return user, nil
}
