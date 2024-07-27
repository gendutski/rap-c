package userusecase

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

	"github.com/labstack/echo/v4"
)

func NewUsecase(cfg *config.Config, userRepo contract.UserRepository) usecasecontract.UserUsecase {
	return &usecase{cfg, userRepo}
}

type usecase struct {
	cfg      *config.Config
	userRepo contract.UserRepository
}

func (uc *usecase) Create(ctx context.Context, payload *payloadentity.CreateUserPayload, author *databaseentity.User) (*databaseentity.User, string, error) {
	// validate payload
	err := entity.InitValidator().Validate(payload)
	if err != nil {
		return nil, "", err
	}

	// generate strong password
	password, err := helper.GenerateStrongPassword()
	if err != nil {
		return nil, "", &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.HelperGenerateStrongPasswordError, err.Error()),
		}
	}
	// encrypt password
	encryptPassword, err := helper.EncryptPassword(password)
	if err != nil {
		return nil, "", &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.HelperEncryptPasswordError, err.Error()),
		}
	}

	// set payload & result
	token, err := helper.GenerateToken(64)
	if err != nil {
		return nil, "", &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.HelperGenerateTokenError, err.Error()),
		}
	}
	user := databaseentity.User{
		Username:           payload.Username,
		FullName:           payload.FullName,
		Email:              payload.Email,
		Password:           encryptPassword,
		PasswordMustChange: true,
		IsGuest:            payload.IsGuest,
		Token:              token,
		CreatedBy:          author.ID,
		UpdatedBy:          author.ID,
	}

	// save
	err = uc.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, "", err
	}
	return &user, password, nil
}

func (uc *usecase) GetUserList(ctx context.Context, req *payloadentity.GetUserListRequest) ([]*databaseentity.User, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	return uc.userRepo.GetUsersByRequest(ctx, req)
}

func (uc *usecase) GetTotalUserList(ctx context.Context, req *payloadentity.GetUserListRequest) (int64, error) {
	return uc.userRepo.GetTotalUsersByRequest(ctx, req)
}

func (uc *usecase) GetUserByUsername(ctx context.Context, req *payloadentity.GetUserDetailRequest) (*databaseentity.User, error) {
	// validate request
	err := entity.InitValidator().Validate(req)
	if err != nil {
		return nil, err
	}

	// get user by username
	user, err := uc.userRepo.GetUserByField(ctx, "username", req.Username, http.StatusNotFound)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *usecase) Update(ctx context.Context, payload *payloadentity.UpdateUserPayload, author *databaseentity.User) error {
	// validate payload
	err := entity.InitValidator().Validate(payload)
	if err != nil {
		return err
	}

	var isModified bool
	if payload.Username != "" {
		author.Username = payload.Username
		isModified = true
	}
	if payload.FullName != "" {
		author.FullName = payload.FullName
		isModified = true
	}
	if payload.Email != "" {
		author.Email = payload.Email
		isModified = true
	}
	if payload.Password != "" {
		// encrypt password
		encryptPassword, err := helper.EncryptPassword(payload.Password)
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  http.StatusText(http.StatusInternalServerError),
				Internal: entity.NewInternalError(entity.HelperEncryptPasswordError, err.Error()),
			}
		}
		author.Password = encryptPassword
		author.PasswordMustChange = false
		isModified = true
	}
	if isModified {
		author.UpdatedBy = author.ID
		return uc.userRepo.Update(ctx, author)
	}
	return &echo.HTTPError{
		Code:     http.StatusConflict,
		Message:  entity.UpdateUserNoChangeMessage,
		Internal: entity.NewInternalError(entity.UpdateUserNoChange, entity.UpdateUserNoChangeMessage),
	}
}

func (uc *usecase) UpdateActiveStatus(ctx context.Context, payload *payloadentity.ActiveStatusPayload, author *databaseentity.User) (*databaseentity.User, error) {
	// validate payload
	err := entity.InitValidator().Validate(payload)
	if err != nil {
		return nil, err
	}

	// get target user
	user, err := uc.userRepo.GetUserByField(ctx, "username", payload.Username, http.StatusNotFound)
	if err != nil {
		return nil, err
	}

	// validate status
	if user.Disabled == payload.Disabled {
		if user.Disabled {
			return nil, &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  entity.DeactivatingInActiveUserMessage,
				Internal: entity.NewInternalError(entity.DeactivatingInActiveUser, entity.DeactivatingInActiveUserMessage),
			}
		} else {
			return nil, &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  entity.ActivatingActiveUserMessage,
				Internal: entity.NewInternalError(entity.ActivatingActiveUser, entity.ActivatingActiveUserMessage),
			}
		}
	}

	// update user
	user.Disabled = payload.Disabled
	user.UpdatedBy = author.ID
	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
