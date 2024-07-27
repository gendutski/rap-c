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
		return []*responseentity.UserResponse{}, nil
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

func (uc *usecase) FormatUnit(ctx context.Context, unit *databaseentity.Unit, mapUsers map[int]string) (*responseentity.UnitResponse, error) {
	var err error
	if unit == nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.FormatterUsecaseFormatUnitError, "empty user"),
		}
	}
	if mapUsers == nil {
		mapUsers, err = uc.userRepo.MapUserUsername(ctx, []int{unit.CreatedBy})
		if err != nil {
			return nil, err
		}
	}
	return uc.formatUnit(unit, mapUsers), nil
}

func (uc *usecase) FormatUnits(ctx context.Context, units []*databaseentity.Unit) ([]*responseentity.UnitResponse, error) {
	if len(units) == 0 {
		return []*responseentity.UnitResponse{}, nil
	}
	mapUsers, err := uc.userRepo.MapUserUsername(ctx, units)
	if err != nil {
		return nil, err
	}
	var result []*responseentity.UnitResponse
	for _, unit := range units {
		result = append(result, uc.formatUnit(unit, mapUsers))
	}
	return result, nil
}
