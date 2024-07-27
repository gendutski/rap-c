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
			CreatedByDB:        1,
			UpdatedByDB:        1,
		})).Return(nil).Times(1)

		res, pass, err := uc.Create(ctx, &payloadentity.CreateUserPayload{
			Username: "gendutski",
			FullName: "Firman Darmawan",
			Email:    "mvp.firman.darmawan@gmail.com",
		}, &databaseentity.User{ID: 1})
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
			CreatedByDB:        1,
			UpdatedByDB:        1,
		})).Return(nil).Times(1)

		res, pass, err := uc.Create(ctx, &payloadentity.CreateUserPayload{
			Username: "gendutski",
			FullName: "Firman Darmawan",
			Email:    "mvp.firman.darmawan@gmail.com",
			IsGuest:  true,
		}, &databaseentity.User{ID: 1})
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

func Test_GetUserByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, userRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		validUser := &databaseentity.User{
			Username: "gendutski",
		}

		userRepo.EXPECT().GetUserByField(ctx, "username", "gendutski", 404).Return(validUser, nil).Times(1)

		res, err := uc.GetUserByUsername(ctx, &payloadentity.GetUserDetailRequest{
			Username: "gendutski",
		})
		assert.Nil(t, err)
		assert.Equal(t, validUser, res)
	})

	t.Run("validator fails", func(t *testing.T) {
		res, err := uc.GetUserByUsername(ctx, &payloadentity.GetUserDetailRequest{})
		assert.Nil(t, res)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, entity.ValidatorBadRequest, herr.Internal.(*entity.InternalError).Code)
	})
}

func Test_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, userRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		author := &databaseentity.User{
			ID:       1,
			Username: "gendutski",
			FullName: "Firman Darmawan",
			Email:    "gendutski@gmail.com",
			Password: "password",
			Token:    "token",
		}

		userRepo.EXPECT().Update(ctx, CreateMatcher(&databaseentity.User{
			Username:           "gendutski-1",
			FullName:           "Lord Firman Darmawan",
			Email:              "mvp.firman.darmawan@gmail.com",
			PasswordMustChange: false,
			UpdatedByDB:        1,
		})).Return(nil).Times(1)

		err := uc.Update(ctx, &payloadentity.UpdateUserPayload{
			Username:        "gendutski-1",
			FullName:        "Lord Firman Darmawan",
			Email:           "mvp.firman.darmawan@gmail.com",
			Password:        "new awesome password",
			ConfirmPassword: "new awesome password",
		}, author)
		assert.Nil(t, err)
		assert.True(t, helper.ValidateEncryptedPassword(author.Password, "new awesome password"))
		assert.False(t, author.PasswordMustChange)
		assert.Equal(t,
			[]string{"gendutski-1", "Lord Firman Darmawan", "mvp.firman.darmawan@gmail.com"},
			[]string{author.Username, author.FullName, author.Email},
		)
	})

	t.Run("no change (empty payload)", func(t *testing.T) {
		author := &databaseentity.User{
			Username: "gendutski",
			FullName: "Firman Darmawan",
			Email:    "gendutski@gmail.com",
			Password: "password",
			Token:    "token",
		}

		err := uc.Update(ctx, &payloadentity.UpdateUserPayload{}, author)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusConflict, herr.Code)
		assert.Equal(t, entity.UpdateUserNoChange, herr.Internal.(*entity.InternalError).Code)
	})

	t.Run("validator failed", func(t *testing.T) {
		author := &databaseentity.User{
			Username: "gendutski",
			FullName: "Firman Darmawan",
			Email:    "gendutski@gmail.com",
			Password: "password",
			Token:    "token",
		}

		err := uc.Update(ctx, &payloadentity.UpdateUserPayload{
			Password:        "short",
			ConfirmPassword: "not match",
		}, author)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, entity.ValidatorBadRequest, herr.Internal.(*entity.InternalError).Code)
		assert.Equal(t, map[string][]*entity.ValidatorMessage{
			"password":        {{Tag: "min", Param: "8"}},
			"confirmPassword": {{Tag: "eqfield", Param: "Password"}},
		}, herr.Message)
	})

}

func Test_UpdateActiveStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, userRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	author := &databaseentity.User{
		ID:       1,
		Username: "gendutski",
	}

	t.Run("deactivating active user", func(t *testing.T) {
		currentUser := &databaseentity.User{
			Username: "other-user",
			Password: "password",
			Token:    "token",
		}
		userRepo.EXPECT().GetUserByField(ctx, "username", "other-user", 404).Return(currentUser, nil).Times(1)
		userRepo.EXPECT().Update(ctx, CreateMatcher(&databaseentity.User{
			Username:    "other-user",
			Password:    "password",
			Disabled:    true,
			UpdatedByDB: author.ID,
		})).Return(nil).Times(1)

		res, err := uc.UpdateActiveStatus(ctx, &payloadentity.ActiveStatusPayload{
			Username: "other-user",
			Disabled: true,
		}, author)
		assert.Nil(t, err)
		assert.True(t, res.Disabled)
	})

	t.Run("activating inactive user", func(t *testing.T) {
		currentUser := &databaseentity.User{
			Username: "other-user",
			Password: "password",
			Disabled: true,
			Token:    "token",
		}
		userRepo.EXPECT().GetUserByField(ctx, "username", "other-user", 404).Return(currentUser, nil).Times(1)
		userRepo.EXPECT().Update(ctx, CreateMatcher(&databaseentity.User{
			Username:    "other-user",
			Password:    "password",
			Disabled:    false,
			UpdatedByDB: author.ID,
		})).Return(nil).Times(1)

		res, err := uc.UpdateActiveStatus(ctx, &payloadentity.ActiveStatusPayload{
			Username: "other-user",
			Disabled: false,
		}, author)
		assert.Nil(t, err)
		assert.False(t, res.Disabled)
	})

	t.Run("activating active user", func(t *testing.T) {
		currentUser := &databaseentity.User{
			Username: "other-user",
			Password: "password",
			Disabled: false,
			Token:    "token",
		}
		userRepo.EXPECT().GetUserByField(ctx, "username", "other-user", 404).Return(currentUser, nil).Times(1)

		_, err := uc.UpdateActiveStatus(ctx, &payloadentity.ActiveStatusPayload{
			Username: "other-user",
			Disabled: false,
		}, author)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, entity.ActivatingActiveUser, herr.Internal.(*entity.InternalError).Code)
	})

	t.Run("deactivating inactive user", func(t *testing.T) {
		currentUser := &databaseentity.User{
			Username: "other-user",
			Password: "password",
			Disabled: true,
			Token:    "token",
		}
		userRepo.EXPECT().GetUserByField(ctx, "username", "other-user", 404).Return(currentUser, nil).Times(1)

		_, err := uc.UpdateActiveStatus(ctx, &payloadentity.ActiveStatusPayload{
			Username: "other-user",
			Disabled: true,
		}, author)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, entity.DeactivatingInActiveUser, herr.Internal.(*entity.InternalError).Code)
	})
}
