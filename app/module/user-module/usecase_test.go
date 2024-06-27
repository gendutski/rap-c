package usermodule_test

import (
	"context"
	"rap-c/app/entity"
	usermodule "rap-c/app/module/user-module"
	repomocks "rap-c/app/repository/contract/mocks"
	"rap-c/config"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func initUsecase(ctrl *gomock.Controller, cfg config.Config) (usermodule.UserModule, *repomocks.MockUserRepository) {
	userRepo := repomocks.NewMockUserRepository(ctrl)
	uc := usermodule.NewUsecase(cfg, userRepo)
	return uc, userRepo
}

func Test_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, userRepo := initUsecase(ctrl, config.Config{})
	ctx := context.Background()

	t.Run("success non guest", func(t *testing.T) {
		userRepo.EXPECT().Create(ctx, CreateMatcher(&entity.User{
			Username:           "gendutski",
			FullName:           "Firman Darmawan",
			Email:              "mvp.firman.darmawan@gmail.com",
			PasswordMustChange: true,
			CreatedBy:          "SYSTEM",
			UpdatedBy:          "SYSTEM",
		})).Return(nil).Times(1)

		res, err := uc.Create(ctx, &entity.CreateUserPayload{
			Username: "gendutski",
			FullName: "Firman Darmawan",
			Email:    "mvp.firman.darmawan@gmail.com",
		}, &entity.User{Username: "SYSTEM"})
		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("success guest", func(t *testing.T) {
		userRepo.EXPECT().Create(ctx, CreateMatcher(&entity.User{
			Username:           "gendutski",
			FullName:           "Firman Darmawan",
			Email:              "mvp.firman.darmawan@gmail.com",
			PasswordMustChange: true,
			IsGuest:            true,
			CreatedBy:          "SYSTEM",
			UpdatedBy:          "SYSTEM",
		})).Return(nil).Times(1)

		res, err := uc.Create(ctx, &entity.CreateUserPayload{
			Username: "gendutski",
			FullName: "Firman Darmawan",
			Email:    "mvp.firman.darmawan@gmail.com",
			IsGuest:  true,
		}, &entity.User{Username: "SYSTEM"})
		assert.Nil(t, err)
		assert.NotNil(t, res)
	})

	t.Run("not valid payload", func(t *testing.T) {
		res, err := uc.Create(ctx, &entity.CreateUserPayload{
			FullName: "Firman Darmawan",
			Email:    "mvp.firman.darmawan@gmail.com",
		}, &entity.User{Username: "gendutski"})
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

	t.Run("empty author", func(t *testing.T) {
		res, err := uc.Create(ctx, &entity.CreateUserPayload{
			Username: "gendutski",
			FullName: "Firman Darmawan",
			Email:    "mvp.firman.darmawan@gmail.com",
		}, nil)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}
