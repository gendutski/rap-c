package usermodule

import (
	"context"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"
	"rap-c/app/repository/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

type UserModule interface {
	// create user
	Create(ctx context.Context, payload *entity.CreateUserPayload, author *entity.User) (*entity.User, error)
}

func NewUsecase(cfg config.Config, userRepo contract.UserRepository) UserModule {
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
