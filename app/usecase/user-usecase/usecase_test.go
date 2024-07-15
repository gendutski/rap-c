package userusecase_test

import (
	"context"
	"net/http"
	"rap-c/app/entity"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
	"rap-c/app/helper"
	repomocks "rap-c/app/repository/contract/mocks"
	"rap-c/app/usecase/contract"
	userusecase "rap-c/app/usecase/user-usecase"
	"rap-c/config"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func initUsecase(ctrl *gomock.Controller, cfg *config.Config) (contract.UserUsecase, *repomocks.MockUserRepository) {
	userRepo := repomocks.NewMockUserRepository(ctrl)
	uc := userusecase.NewUsecase(cfg, userRepo)
	return uc, userRepo
}

func Test_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, userRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	t.Run("success non guest", func(t *testing.T) {
		userRepo.EXPECT().Create(ctx, CreateMatcher(&databaseentity.User{
			Username:           "gendutski",
			FullName:           "Firman Darmawan",
			Email:              "mvp.firman.darmawan@gmail.com",
			PasswordMustChange: true,
			CreatedBy:          "SYSTEM",
			UpdatedBy:          "SYSTEM",
		})).Return(nil).Times(1)

		res, pass, err := uc.Create(ctx, &payloadentity.CreateUserPayload{
			Username: "gendutski",
			FullName: "Firman Darmawan",
			Email:    "mvp.firman.darmawan@gmail.com",
		}, &databaseentity.User{Username: "SYSTEM"})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.True(t, helper.ValidateEncryptedPassword(res.Password, pass))
	})

	t.Run("success guest", func(t *testing.T) {
		userRepo.EXPECT().Create(ctx, CreateMatcher(&databaseentity.User{
			Username:           "gendutski",
			FullName:           "Firman Darmawan",
			Email:              "mvp.firman.darmawan@gmail.com",
			PasswordMustChange: true,
			IsGuest:            true,
			CreatedBy:          "SYSTEM",
			UpdatedBy:          "SYSTEM",
		})).Return(nil).Times(1)

		res, pass, err := uc.Create(ctx, &payloadentity.CreateUserPayload{
			Username: "gendutski",
			FullName: "Firman Darmawan",
			Email:    "mvp.firman.darmawan@gmail.com",
			IsGuest:  true,
		}, &databaseentity.User{Username: "SYSTEM"})
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.True(t, helper.ValidateEncryptedPassword(res.Password, pass))
	})

	t.Run("not valid payload", func(t *testing.T) {
		_, _, err := uc.Create(ctx, &payloadentity.CreateUserPayload{
			FullName: "Firman Darmawan",
			Email:    "gendutski.gmail.com",
		}, &databaseentity.User{Username: "gendutski"})
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, entity.ValidatorBadRequest, herr.Internal.(*entity.InternalError).Code)
		assert.Equal(t, map[string][]*entity.ValidatorMessage{
			"email":    {{Tag: "email"}},
			"username": {{Tag: "required"}},
		}, herr.Message)
	})
}
