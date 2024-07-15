package authusecase

import (
	"context"
	"net/http"
	"rap-c/app/entity"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
	"rap-c/app/helper"
	"rap-c/app/repository/contract"
	usecasecontract "rap-c/app/usecase/contract"
	"rap-c/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	tokenStrID    string = "id"
	tokenStrName  string = "userName"
	tokenStrEmail string = "email"
)

func NewUsecase(cfg *config.Config, authRepo contract.AuthRepository) usecasecontract.AuthUsecase {
	return &usecase{cfg, authRepo}
}

type usecase struct {
	cfg      *config.Config
	authRepo contract.AuthRepository
}

func (uc *usecase) AttemptLogin(ctx context.Context, payload *payloadentity.AttemptLoginPayload) (*databaseentity.User, error) {
	// validate payload
	err := entity.InitValidator().Validate(payload)
	if err != nil {
		return nil, err
	}

	// do login
	user, err := uc.authRepo.DoUserLogin(ctx, payload)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *usecase) AttemptGuestLogin(ctx context.Context) (*databaseentity.User, error) {
	if !uc.cfg.EnableGuestLogin() {
		return nil, &echo.HTTPError{
			Code:     http.StatusForbidden,
			Message:  entity.AttemptGuestLoginForbiddenMessage,
			Internal: entity.NewInternalError(entity.AttemptGuestLoginForbidden, entity.AttemptGuestLoginForbiddenMessage),
		}
	}
	user, err := uc.authRepo.DoUserLogin(ctx, &payloadentity.AttemptLoginPayload{
		Email:    config.GuestEmail,
		Password: config.GuestPassword,
	})
	if err != nil {
		return nil, err
	}
	if !user.IsGuest {
		return nil, &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  entity.NonGuestAttemptGuestLoginMessage,
			Internal: entity.NewInternalError(entity.NonGuestAttemptGuestLogin, entity.NonGuestAttemptGuestLoginMessage),
		}
	}

	return user, nil
}

func (uc *usecase) GenerateJwtToken(ctx context.Context, user *databaseentity.User, isLongSession bool) (string, error) {
	exp := time.Now().Add(time.Minute * time.Duration(uc.cfg.JwtExpirationInMinutes())).Unix()
	if isLongSession {
		// long session for remember login session
		exp = time.Now().Add(time.Hour * 24 * time.Duration(uc.cfg.JwtRememberInDays())).Unix()
	}
	claims := jwt.MapClaims{
		tokenStrID:    user.ID,
		tokenStrName:  user.Username,
		tokenStrEmail: user.Email,
		"exp":         exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(uc.cfg.JwtSecret()))
	if err != nil {
		return "", &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AuthUsecaseGenerateJwtTokenError, err.Error()),
		}
	}
	return tokenStr, nil
}

func (uc *usecase) ValidateJwtToken(ctx context.Context, token *jwt.Token, guestAccepted bool) (*databaseentity.User, error) {
	// get claims
	claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
	if !ok {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AuthUsecaseValidateJwtTokenError, "failed to cast claims as jwt.MapClaims"),
		}
	}

	// get user
	email, ok := claims[tokenStrEmail].(string)
	if !ok {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AuthUsecaseValidateJwtTokenError, "jwt claims email not valid"),
		}
	}
	user, err := uc.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// check user
	var validUserID bool
	switch t := claims[tokenStrID].(type) {
	case int:
		validUserID = t == user.ID
	case int64:
		validUserID = t == int64(user.ID)
	case float32:
		validUserID = t == float32(user.ID)
	case float64:
		validUserID = t == float64(user.ID)
	}
	if user.Username != claims[tokenStrName].(string) || !validUserID || user.Disabled {
		return nil, &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  entity.ValidateTokenFailedMessage,
			Internal: entity.NewInternalError(entity.ValidateTokenFailed, entity.ValidateTokenFailedMessage),
		}
	}

	// if guest not accepted
	if !guestAccepted && user.IsGuest {
		return nil, &echo.HTTPError{
			Code:     http.StatusForbidden,
			Message:  entity.GuestTokenForbiddenMessage,
			Internal: entity.NewInternalError(entity.GuestTokenForbidden, entity.GuestTokenForbiddenMessage),
		}
	}
	return user, nil
}

func (uc *usecase) RenewPassword(ctx context.Context, user *databaseentity.User, payload *payloadentity.RenewPasswordPayload) error {
	// validate payload
	err := entity.InitValidator().Validate(payload)
	if err != nil {
		return err
	}

	// check whether the new password is the same as the old password
	if helper.ValidateEncryptedPassword(user.Password, payload.Password) {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  entity.RenewPasswordWithUnchangedPasswordMessage,
			Internal: entity.NewInternalError(entity.RenewPasswordWithUnchangedPassword, entity.RenewPasswordWithUnchangedPasswordMessage),
		}
	}

	// update
	return uc.authRepo.DoRenewPassword(ctx, user, payload)
}

func (uc *usecase) RequestResetPassword(ctx context.Context, payload *payloadentity.RequestResetPayload) (*databaseentity.User, *databaseentity.PasswordResetToken, error) {
	// validate payload
	err := entity.InitValidator().Validate(payload)
	if err != nil {
		return nil, nil, err
	}

	// check email
	user, err := uc.authRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return nil, nil, err
	}

	// generate reset password token
	token, err := uc.authRepo.GenerateUserResetPassword(ctx, payload)
	if err != nil {
		return nil, nil, err
	}

	return user, token, nil
}

func (uc *usecase) ValidateResetToken(ctx context.Context, payload *payloadentity.ValidateResetTokenPayload) error {
	// validate payload
	err := entity.InitValidator().Validate(payload)
	if err != nil {
		return err
	}
	_, err = uc.authRepo.ValidateResetToken(ctx, payload)
	if err != nil {
		return err
	}
	return nil
}

func (uc *usecase) SubmitResetPassword(ctx context.Context, payload *payloadentity.ResetPasswordPayload) (*databaseentity.User, error) {
	// validate payload
	err := entity.InitValidator().Validate(payload)
	if err != nil {
		return nil, err
	}

	// get reset pasword data
	reset, err := uc.authRepo.ValidateResetToken(ctx, &payloadentity.ValidateResetTokenPayload{
		Email: payload.Email,
		Token: payload.Token,
	})
	if err != nil {
		return nil, err
	}

	// get user
	user, err := uc.authRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	// encrypt password
	encryptPass, err := helper.EncryptPassword(payload.Password)
	if err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.HelperEncryptPasswordError, err.Error()),
		}
	}

	// update user
	user.Password = encryptPass
	user.UpdatedBy = user.Username

	// update reset token
	reset.Token = ""

	// save
	err = uc.authRepo.DoResetPassword(ctx, user, reset)
	if err != nil {
		return nil, err
	}
	return user, nil
}
