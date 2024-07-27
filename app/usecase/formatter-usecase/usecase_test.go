package formatterusecase_test

import (
	"context"
	"net/http"
	"rap-c/app/entity"
	databaseentity "rap-c/app/entity/database-entity"
	responseentity "rap-c/app/entity/response-entity"
	"rap-c/app/repository/contract/mocks"
	"rap-c/app/usecase/contract"
	formatterusecase "rap-c/app/usecase/formatter-usecase"
	"rap-c/config"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func initUsecase(ctrl *gomock.Controller, cfg *config.Config) (contract.FormatterUsecase, *mocks.MockUserRepository) {
	userRepo := mocks.NewMockUserRepository(ctrl)
	usecase := formatterusecase.NewUsecase(cfg, userRepo)
	return usecase, userRepo
}

func Test_FormatUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, userRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	t.Run("success, use map in db", func(t *testing.T) {
		user := &databaseentity.User{Username: "user-2", CreatedBy: 1, UpdatedBy: 1}
		userRepo.EXPECT().MapUserUsername(ctx, user).Return(map[int]string{
			1: "user-1",
		}, nil).Times(1)

		resp, err := uc.FormatUser(ctx, user, nil)
		assert.Nil(t, err)
		assert.Equal(t, &responseentity.UserResponse{
			Username:  "user-2",
			CreatedBy: "user-1",
			UpdatedBy: "user-1",
		}, resp)
	})

	t.Run("success, use map in param", func(t *testing.T) {
		user := &databaseentity.User{Username: "user-2", CreatedBy: 1, UpdatedBy: 1}

		resp, err := uc.FormatUser(ctx, user, map[int]string{1: "gendutski"})
		assert.Nil(t, err)
		assert.Equal(t, &responseentity.UserResponse{
			Username:  "user-2",
			CreatedBy: "gendutski",
			UpdatedBy: "gendutski",
		}, resp)
	})

	t.Run("empty users", func(t *testing.T) {
		resp, err := uc.FormatUser(ctx, nil, map[int]string{1: "gendutski"})
		assert.Nil(t, resp)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, herr.Code)
		assert.Equal(t, entity.FormatterUsecaseFormatUserError, herr.Internal.(*entity.InternalError).Code)
	})
}

func Test_FormatUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, userRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		users := []*databaseentity.User{
			{Username: "user-1"},
			{Username: "user-2", CreatedBy: 1, UpdatedBy: 1},
		}
		userRepo.EXPECT().MapUserUsername(ctx, users).Return(map[int]string{
			1: "user-1",
		}, nil).Times(1)

		resp, err := uc.FormatUsers(ctx, users)
		assert.Nil(t, err)
		assert.Equal(t, []*responseentity.UserResponse{
			{Username: "user-1", CreatedBy: "SYSTEM", UpdatedBy: "SYSTEM"},
			{Username: "user-2", CreatedBy: "user-1", UpdatedBy: "user-1"},
		}, resp)
	})

	t.Run("empty users", func(t *testing.T) {
		resp, err := uc.FormatUsers(ctx, []*databaseentity.User{})
		assert.Nil(t, err)
		assert.Empty(t, resp)
	})
}

func Test_FormatUnit(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, userRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	t.Run("success, use map in db", func(t *testing.T) {
		userRepo.EXPECT().MapUserUsername(ctx, []int{1}).Return(map[int]string{
			1: "user-1",
		}, nil).Times(1)

		resp, err := uc.FormatUnit(ctx, &databaseentity.Unit{
			Name:      "Kilogram",
			CreatedBy: 1,
		}, nil)
		assert.Nil(t, err)
		assert.Equal(t, &responseentity.UnitResponse{
			Name:      "Kilogram",
			CreatedBy: "user-1",
		}, resp)
	})

	t.Run("success, use map in param", func(t *testing.T) {
		resp, err := uc.FormatUnit(ctx, &databaseentity.Unit{
			Name:      "Kilogram",
			CreatedBy: 1,
		}, map[int]string{1: "gendutski"})
		assert.Nil(t, err)
		assert.Equal(t, &responseentity.UnitResponse{
			Name:      "Kilogram",
			CreatedBy: "gendutski",
		}, resp)
	})

	t.Run("empty unit", func(t *testing.T) {
		resp, err := uc.FormatUnit(ctx, nil, map[int]string{1: "gendutski"})
		assert.Nil(t, resp)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, herr.Code)
		assert.Equal(t, entity.FormatterUsecaseFormatUnitError, herr.Internal.(*entity.InternalError).Code)
	})
}

func Test_FormatUnits(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, userRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		payload := []*databaseentity.Unit{
			{Name: "Kilogram", CreatedBy: 1},
			{Name: "Gram", CreatedBy: 2},
		}
		userRepo.EXPECT().MapUserUsername(ctx, payload).Return(map[int]string{
			1: "user-1",
		}, nil).Times(1)

		resp, err := uc.FormatUnits(ctx, payload)
		assert.Nil(t, err)
		assert.Equal(t, []*responseentity.UnitResponse{
			{Name: "Kilogram", CreatedBy: "user-1"},
			{Name: "Gram", CreatedBy: "DELETED USER"},
		}, resp)
	})

	t.Run("empty units", func(t *testing.T) {
		resp, err := uc.FormatUnits(ctx, []*databaseentity.Unit{})
		assert.Nil(t, err)
		assert.Empty(t, resp)
	})
}
