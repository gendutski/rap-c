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
	user := databaseentity.User{
		Username:           payload.Username,
		FullName:           payload.FullName,
		Email:              payload.Email,
		Password:           encryptPassword,
		PasswordMustChange: true,
		IsGuest:            payload.IsGuest,
		CreatedBy:          author.Username,
		UpdatedBy:          author.Username,
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
	if req.Page < 1 {
		req.Page = 1
	}
	return uc.userRepo.GetTotalUsersByRequest(ctx, req)
}

func (uc *usecase) GetUserByUsername(ctx context.Context, username string) (*databaseentity.User, error) {
	user, err := uc.userRepo.GetUserByField(ctx, "username", username, http.StatusNotFound)
	if err != nil {
		return nil, err
	}
	return user, nil
}
