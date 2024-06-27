package usermodule

import (
	"context"
	"errors"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"
	"rap-c/app/repository/contract"
	"rap-c/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	tokenStrID    string = "id"
	tokenStrName  string = "userName"
	tokenStrEmail string = "email"
)

type UserUsecase interface {
	// create user
	Create(ctx context.Context, payload *entity.CreateUserPayload, author *entity.User) (*entity.User, error)
	// attempt to login with email and password
	AttemptLogin(ctx context.Context, payload *entity.AttemptLoginPayload) (*entity.User, error)
	// convert user to jwt token
	GenerateJwtToken(ctx context.Context, user *entity.User) (string, error)
	// validate jwt token into user
	ValidateJwtToken(ctx context.Context, token *jwt.Token, guestAccepted bool) (*entity.User, error)
}

func NewUsecase(cfg config.Config, userRepo contract.UserRepository) UserUsecase {
	return &usecase{cfg, userRepo}
}

type usecase struct {
	cfg      config.Config
	userRepo contract.UserRepository
}

func (uc *usecase) Create(ctx context.Context, payload *entity.CreateUserPayload, author *entity.User) (*entity.User, error) {
	// validate payload
	validate := helper.GenerateStructValidator()
	errMessages := payload.Validate(validate)
	if len(errMessages) > 0 {
		return nil, &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  errMessages,
			Internal: entity.NewInternalError(entity.ValidateCreateUserFailed, errMessages...),
		}
	}
	// validate author
	if author == nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  entity.CreateUserErrorEmptyAuthor,
			Internal: entity.NewInternalError(entity.CreateUserError, entity.CreateUserErrorEmptyAuthor),
		}
	}

	// generate password
	pass, err := uc.generateUserPassword("")
	if err != nil {
		return nil, err
	}

	// set payload & result
	user := entity.User{
		Username:           payload.Username,
		FullName:           payload.FullName,
		Email:              payload.Email,
		Password:           pass,
		PasswordMustChange: true,
		IsGuest:            payload.IsGuest,
		CreatedBy:          author.Username,
		UpdatedBy:          author.Username,
	}

	// save
	err = uc.userRepo.Create(ctx, &user)
	if err != nil {
		if echoError, ok := err.(*echo.HTTPError); ok {
			return nil, echoError
		}
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.CreateUserError, err.Error()),
		}
	}
	return &user, nil
}

func (uc *usecase) AttemptLogin(ctx context.Context, payload *entity.AttemptLoginPayload) (*entity.User, error) {
	// validate payload
	validate := helper.GenerateStructValidator()
	errMessages := payload.Validate(validate)
	if len(errMessages) > 0 {
		return nil, &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  errMessages,
			Internal: entity.NewInternalError(entity.ValidateAttemptLoginFailed, errMessages...),
		}
	}

	// get user by email
	user, err := uc.userRepo.GetUserByField(ctx, "email", payload.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// user not found
			return nil, &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  entity.AttemptLoginIncorrectMessage,
				Internal: entity.NewInternalError(entity.ValidateAttemptLoginFailed, entity.AttemptLoginIncorrectMessage),
			}
		}
		// internal database error
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AttemptLoginError, err.Error()),
		}
	}

	// check user status
	if user.Disabled {
		return nil, &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  entity.AttemptLoginDisabledMessage,
			Internal: entity.NewInternalError(entity.AttemptLoginFailed, entity.AttemptLoginDisabledMessage),
		}
	}

	// validate password
	if !helper.ValidateEncryptedPassword(user.Password, payload.Password) {
		return nil, &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  entity.AttemptLoginIncorrectMessage,
			Internal: entity.NewInternalError(entity.AttemptLoginFailed, entity.AttemptLoginIncorrectMessage),
		}
	}

	return user, nil
}

func (uc *usecase) GenerateJwtToken(ctx context.Context, user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		tokenStrID:    user.ID,
		tokenStrName:  user.Username,
		tokenStrEmail: user.Email,
		"exp":         time.Now().Add(time.Minute * time.Duration(uc.cfg.JwtExpirationInMinutes)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(uc.cfg.JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (uc *usecase) ValidateJwtToken(ctx context.Context, token *jwt.Token, guestAccepted bool) (*entity.User, error) {
	// get claims
	claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
	if !ok {
		return nil, errors.New("failed to cast claims as jwt.MapClaims")
	}

	// get user
	user, err := uc.userRepo.GetUserByField(ctx, "id", claims[tokenStrID])
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// user not found
			return nil, &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  entity.ValidateTokenUserNotFoundMessage,
				Internal: entity.NewInternalError(entity.ValidateTokenUserNotFound, entity.ValidateTokenUserNotFoundMessage),
			}
		}
		// internal database error
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.ValidateTokenDBError, err.Error()),
		}
	}

	// check user
	if user.Username != claims[tokenStrName].(string) || user.Email != claims[tokenStrEmail].(string) || user.Disabled {
		return nil, &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  entity.ValidateTokenUserNotMatchMessage,
			Internal: entity.NewInternalError(entity.ValidateTokenUserNotMatch, entity.ValidateTokenUserNotMatchMessage),
		}
	}

	// if guest not accepted
	if !guestAccepted && user.IsGuest {
		return nil, &echo.HTTPError{
			Code:     http.StatusForbidden,
			Message:  entity.ValidateTokenGuestNotAcceptedMessage,
			Internal: entity.NewInternalError(entity.ValidateTokenGuestNotAccepted, entity.ValidateTokenGuestNotAcceptedMessage),
		}
	}

	return user, nil
}
