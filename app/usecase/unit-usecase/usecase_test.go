package unitusecase_test

import (
	"context"
	"net/http"
	"rap-c/app/entity"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
	"rap-c/app/repository/contract/mocks"
	"rap-c/app/usecase/contract"
	unitusecase "rap-c/app/usecase/unit-usecase"
	"rap-c/config"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func initUsecase(ctrl *gomock.Controller, cfg *config.Config) (contract.UnitUsecase, *mocks.MockUnitRepository) {
	unitRepo := mocks.NewMockUnitRepository(ctrl)
	usecase := unitusecase.NewUsecase(cfg, unitRepo)
	return usecase, unitRepo
}

func Test_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	usecase, unitRepo := initUsecase(ctrl, nil)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		validUnit := &databaseentity.Unit{
			Name:      "Kilogram",
			CreatedBy: 1,
		}
		unitRepo.EXPECT().Create(ctx, validUnit).Return(nil).Times(1)

		resp, err := usecase.Create(ctx, &payloadentity.CreateUnitPayload{
			Name: "Kilogram",
		}, &databaseentity.User{ID: 1})
		assert.Nil(t, err)
		assert.Equal(t, validUnit, resp)
	})

	t.Run("empty payload", func(t *testing.T) {
		_, err := usecase.Create(ctx, &payloadentity.CreateUnitPayload{}, &databaseentity.User{ID: 1})
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, map[string][]*entity.ValidatorMessage{
			"name": {{Tag: "required"}},
		}, herr.Message)
	})

	t.Run("unit name exceded max", func(t *testing.T) {
		_, err := usecase.Create(ctx, &payloadentity.CreateUnitPayload{
			Name: "very long unit name that exceeded max 30 characters",
		}, &databaseentity.User{ID: 1})
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, map[string][]*entity.ValidatorMessage{
			"name": {{Tag: "max", Param: "30"}},
		}, herr.Message)
	})
}

func Test_DeleteCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	usecase, unitRepo := initUsecase(ctrl, nil)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		validUnit := &databaseentity.Unit{
			Name:      "Kilogram",
			CreatedBy: 1,
		}
		unitRepo.EXPECT().GetUnitByName(ctx, "Kilogram").Return(validUnit, nil).Times(1)
		unitRepo.EXPECT().Delete(ctx, validUnit).Return(nil).Times(1)

		err := usecase.Delete(ctx, &payloadentity.DeleteUnitPayload{
			Name: "Kilogram",
		})
		assert.Nil(t, err)
	})

	t.Run("unit not exists", func(t *testing.T) {
		unitRepo.EXPECT().GetUnitByName(ctx, "Kilogram").Return(nil, gorm.ErrRecordNotFound).Times(1)

		err := usecase.Delete(ctx, &payloadentity.DeleteUnitPayload{
			Name: "Kilogram",
		})
		assert.NotNil(t, err)
	})

	t.Run("empty payload", func(t *testing.T) {
		err := usecase.Delete(ctx, &payloadentity.DeleteUnitPayload{})
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, map[string][]*entity.ValidatorMessage{
			"name": {{Tag: "required"}},
		}, herr.Message)
	})
}
