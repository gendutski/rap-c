package authusecase_test

import (
	"context"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"
	repomocks "rap-c/app/repository/contract/mocks"
	authusecase "rap-c/app/usecase/auth-usecase"
	"rap-c/app/usecase/contract"
	"rap-c/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func initUsecase(ctrl *gomock.Controller, cfg config.Config) (contract.AuthUsecase, *repomocks.MockUserRepository) {
	userRepo := repomocks.NewMockUserRepository(ctrl)
	uc := authusecase.NewUsecase(cfg, userRepo)
	return uc, userRepo
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
		userRepo.EXPECT().GetUserByField(ctx, "email", validUser.Email, http.StatusBadRequest).Return(&validUser, nil).Times(1)

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
		userRepo.EXPECT().GetUserByField(ctx, "email", validUser.Email, http.StatusBadRequest).Return(&validUser, nil).Times(1)

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
		userRepo.EXPECT().GetUserByField(ctx, "email", validUser.Email, http.StatusBadRequest).Return(&validUser, nil).Times(1)

		res, err := uc.AttemptLogin(ctx, &entity.AttemptLoginPayload{
			Email:    "mvp.firman.darmawan@gmail.com",
			Password: "password123",
		})
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
		assert.Nil(t, res)
	})

	t.Run("failed, user not found", func(t *testing.T) {
		userRepo.EXPECT().GetUserByField(ctx, "email", "mvp.firman.darmawan@gmail.com", http.StatusBadRequest).Return(nil, &echo.HTTPError{
			Code: http.StatusBadRequest,
		}).Times(1)

		res, err := uc.AttemptLogin(ctx, &entity.AttemptLoginPayload{
			Email:    "mvp.firman.darmawan@gmail.com",
			Password: "password123",
		})
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
		assert.Nil(t, res)
	})
}

func Test_GenerateJwtToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := config.Config{
		JwtSecret:              "secret",
		JwtExpirationInMinutes: 2,
		JwtRememberInDays:      2,
	}
	uc, _ := initUsecase(ctrl, cfg)
	ctx := context.Background()
	t.Run("success, short time token session", func(t *testing.T) {
		exp := time.Now().Add(time.Minute * time.Duration(cfg.JwtExpirationInMinutes)).Unix()

		tokenStr, err := uc.GenerateJwtToken(ctx, &entity.User{
			ID:       1,
			Username: "gendutski",
			Email:    "mvp.firman.darmawan@gmail.com",
		}, false)

		assert.Nil(t, err)

		// parse token
		claims := jwt.MapClaims{}
		jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecret), nil
		})
		assert.Equal(t, float64(1), claims["id"])
		assert.Equal(t, "gendutski", claims["userName"])
		assert.Equal(t, "mvp.firman.darmawan@gmail.com", claims["email"])
		assert.Equal(t, float64(exp), claims["exp"])
	})

	t.Run("success, long time token session", func(t *testing.T) {
		exp := time.Now().Add(time.Hour * 24 * time.Duration(cfg.JwtRememberInDays))
		tokenStr, err := uc.GenerateJwtToken(ctx, &entity.User{
			ID:       1,
			Username: "gendutski",
			Email:    "mvp.firman.darmawan@gmail.com",
		}, true)

		assert.Nil(t, err)

		// parse token
		claims := jwt.MapClaims{}
		jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecret), nil
		})
		assert.Equal(t, float64(1), claims["id"])
		assert.Equal(t, "gendutski", claims["userName"])
		assert.Equal(t, "mvp.firman.darmawan@gmail.com", claims["email"])
		assert.Equal(t, exp.Format("2006-01-02"), time.Unix(int64(claims["exp"].(float64)), 0).Format("2006-01-02"))
	})
}

func Test_RenewPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, userRepo := initUsecase(ctrl, config.Config{})
	ctx := context.Background()

	oldPassword := "0LdF4s#ionP455W0rd"
	encrypted, _ := helper.EncryptPassword(oldPassword)

	t.Run("success", func(t *testing.T) {
		user := entity.User{
			ID:                 1,
			Username:           "gendutski",
			Password:           encrypted,
			PasswordMustChange: true,
		}

		pass := "Tr!al123#"
		userRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil).Times(1)

		err := uc.RenewPassword(ctx, &user, &entity.RenewPasswordPayload{
			Password:        pass,
			ConfirmPassword: pass,
		})
		assert.Nil(t, err)
		assert.True(t, helper.ValidateEncryptedPassword(user.Password, pass))
		assert.False(t, user.PasswordMustChange)
	})

	t.Run("failed, password not match", func(t *testing.T) {
		user := entity.User{
			ID:                 1,
			Username:           "gendutski",
			Password:           encrypted,
			PasswordMustChange: true,
		}

		pass := "Tr!al123#"
		err := uc.RenewPassword(ctx, &user, &entity.RenewPasswordPayload{
			Password:        pass,
			ConfirmPassword: "password",
		})
		assert.NotNil(t, err)
	})

	t.Run("failed, password not changed", func(t *testing.T) {
		user := entity.User{
			ID:                 1,
			Username:           "gendutski",
			Password:           encrypted,
			PasswordMustChange: true,
		}

		err := uc.RenewPassword(ctx, &user, &entity.RenewPasswordPayload{
			Password:        oldPassword,
			ConfirmPassword: oldPassword,
		})
		assert.NotNil(t, err)
	})
}
