package userusecase

import (
	"context"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"
	"rap-c/app/repository/contract"
	usecasecontract "rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

func NewUsecase(cfg config.Config, userRepo contract.UserRepository) usecasecontract.UserUsecase {
	return &usecase{cfg, userRepo}
}

type usecase struct {
	cfg      config.Config
	userRepo contract.UserRepository
}

func (uc *usecase) Create(ctx context.Context, payload *entity.CreateUserPayload, author *entity.User) (*entity.User, string, error) {
	// validate payload
	validate := helper.GenerateStructValidator()
	errMessages := payload.Validate(validate)
	if len(errMessages) > 0 {
		return nil, "", &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  errMessages,
			Internal: entity.NewInternalError(entity.ValidatorNotValid, errMessages...),
		}
	}

	// generate password
	pass, encryptedPass, err := uc.generateUserPassword("")
	if err != nil {
		return nil, "", err
	}

	// set payload & result
	user := entity.User{
		Username:           payload.Username,
		FullName:           payload.FullName,
		Email:              payload.Email,
		Password:           encryptedPass,
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
	return &user, pass, nil
}

func (uc *usecase) GetUserList(ctx context.Context, req *entity.GetUserListRequest) ([]*entity.User, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	return uc.userRepo.GetUsersByRequest(ctx, req)
}

func (uc *usecase) GetTotalUserList(ctx context.Context, req *entity.GetUserListRequest) (int64, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	return uc.userRepo.GetTotalUsersByRequest(ctx, req)
}

func (uc *usecase) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, err := uc.userRepo.GetUserByField(ctx, "username", username, http.StatusNotFound)
	if err != nil {
		return nil, err
	}
	return user, nil
}
