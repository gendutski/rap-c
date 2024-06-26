package usermodule_test

import (
	"context"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"
	usermodule "rap-c/app/module/user-module"
	repomocks "rap-c/app/repository/contract/mocks"
	"rap-c/config"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func initUsecase(ctrl *gomock.Controller, cfg config.Config) (usermodule.UserUsecase, *repomocks.MockUserRepository) {
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

func Test_AttemptLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, userRepo := initUsecase(ctrl, config.Config{})
	ctx := context.Background()
	dt, _ := time.Parse("2006-01-02 15:04:05", "2024-06-25 01:14:36")

	password := "password"
	hashPass, _ := helper.EncryptPassword(password)

	t.Run("success", func(t *testing.T) {
		validUser := entity.User{
			ID:                 1,
			Username:           "gendutski",
			FullName:           "Firman Darmawan",
			Email:              "mvp.firman.darmawan@gmail.com",
			Password:           hashPass,
			PasswordMustChange: false,
			CreatedAt:          dt,
			UpdatedAt:          dt,
		}
		userRepo.EXPECT().GetUserByField(ctx, "email", validUser.Email).Return(&validUser, nil).Times(1)

		res, err := uc.AttemptLogin(ctx, &entity.AttemptLoginPayload{
			Email:    "mvp.firman.darmawan@gmail.com",
			Password: password,
		})
		assert.Nil(t, err)
		assert.Equal(t, &validUser, res)
	})

	t.Run("failed, disable user", func(t *testing.T) {
		validUser := entity.User{
			ID:                 1,
			Username:           "gendutski",
			FullName:           "Firman Darmawan",
			Email:              "mvp.firman.darmawan@gmail.com",
			Password:           hashPass,
			PasswordMustChange: false,
			Disabled:           true,
			CreatedAt:          dt,
			UpdatedAt:          dt,
		}
		userRepo.EXPECT().GetUserByField(ctx, "email", validUser.Email).Return(&validUser, nil).Times(1)

		res, err := uc.AttemptLogin(ctx, &entity.AttemptLoginPayload{
			Email:    "mvp.firman.darmawan@gmail.com",
			Password: password,
		})
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
		assert.Nil(t, res)
	})

	t.Run("failed, incorrect password", func(t *testing.T) {
		validUser := entity.User{
			ID:                 1,
			Username:           "gendutski",
			FullName:           "Firman Darmawan",
			Email:              "mvp.firman.darmawan@gmail.com",
			Password:           hashPass,
			PasswordMustChange: false,
			CreatedAt:          dt,
			UpdatedAt:          dt,
		}
		userRepo.EXPECT().GetUserByField(ctx, "email", validUser.Email).Return(&validUser, nil).Times(1)

		res, err := uc.AttemptLogin(ctx, &entity.AttemptLoginPayload{
			Email:    "mvp.firman.darmawan@gmail.com",
			Password: "password123",
		})
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
		assert.Nil(t, res)
	})

	t.Run("failed, user not found", func(t *testing.T) {
		userRepo.EXPECT().GetUserByField(ctx, "email", "mvp.firman.darmawan@gmail.com").Return(nil, gorm.ErrRecordNotFound).Times(1)

		res, err := uc.AttemptLogin(ctx, &entity.AttemptLoginPayload{
			Email:    "mvp.firman.darmawan@gmail.com",
			Password: "password123",
		})
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
		assert.Nil(t, res)
	})
}
