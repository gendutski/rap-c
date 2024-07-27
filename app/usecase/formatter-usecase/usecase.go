package formatterusecase

import (
	"context"
	"net/http"
	"rap-c/app/entity"
	databaseentity "rap-c/app/entity/database-entity"
	responseentity "rap-c/app/entity/response-entity"
	"rap-c/app/repository/contract"
	usecasecontract "rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

const (
	defaultUserUsername string = "SYSTEM"
	createdByString     string = "createdBy"
	updatedByString     string = "updatedBy"
)

func NewUsecase(cfg *config.Config, userRepo contract.UserRepository) usecasecontract.FormatterUsecase {
	return &usecase{cfg, userRepo}
}

type usecase struct {
	cfg      *config.Config
	userRepo contract.UserRepository
}

func (uc *usecase) FormatUser(ctx context.Context, user *databaseentity.User, mapUsers map[int]string) (*responseentity.UserResponse, error) {
	var err error
	if user == nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.FormatterUsecaseFormatUserError, "empty user"),
		}
	}
	if mapUsers == nil {
		mapUsers, err = uc.userRepo.MapUserUsername(ctx, user)
		if err != nil {
			return nil, err
		}
	}
	return uc.formatUser(user, mapUsers), nil
}

func (uc *usecase) FormatUsers(ctx context.Context, users []*databaseentity.User) ([]*responseentity.UserResponse, error) {
	if len(users) == 0 {
		return nil, nil
	}
	mapUsers, err := uc.userRepo.MapUserUsername(ctx, users)
	if err != nil {
		return nil, err
	}
	var result []*responseentity.UserResponse
	for _, usr := range users {
		result = append(result, uc.formatUser(usr, mapUsers))
	}
	return result, nil
}
