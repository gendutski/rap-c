package authusecase_test

import (
	"context"
	"errors"
	"net/http"
	"rap-c/app/entity"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
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

func initUsecase(ctrl *gomock.Controller, cfg *config.Config) (contract.AuthUsecase, *repomocks.MockAuthRepository) {
	authRepo := repomocks.NewMockAuthRepository(ctrl)
	uc := authusecase.NewUsecase(cfg, authRepo)
	return uc, authRepo
}

func Test_AttemptLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, authRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()
	password := "password"

	t.Run("success", func(t *testing.T) {
		validUser := databaseentity.User{
			ID:       1,
			Email:    "gendutski@gmail.com",
			Password: password,
		}
		payload := &payloadentity.AttemptLoginPayload{
			Email:    "gendutski@gmail.com",
			Password: password,
		}
		authRepo.EXPECT().DoUserLogin(ctx, payload).Return(&validUser, nil).Times(1)

		res, err := uc.AttemptLogin(ctx, payload)
		assert.Nil(t, err)
		assert.Equal(t, &validUser, res)
	})

	t.Run("empty email", func(t *testing.T) {
		payload := &payloadentity.AttemptLoginPayload{
			Password: password,
		}

		res, err := uc.AttemptLogin(ctx, payload)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, map[string][]*entity.ValidatorMessage{
			"email": {{Tag: "required"}},
		}, herr.Message)
	})

	t.Run("not valid email & empty password", func(t *testing.T) {
		payload := &payloadentity.AttemptLoginPayload{
			Email: "foo",
		}

		res, err := uc.AttemptLogin(ctx, payload)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, map[string][]*entity.ValidatorMessage{
			"email":    {{Tag: "email"}},
			"password": {{Tag: "required"}},
		}, herr.Message)
	})
}

func Test_AttemptGuestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		uc, authRepo := initUsecase(ctrl, config.InitTestConfig(map[string]string{"ENABLE_GUEST_LOGIN": "true"}))

		validUser := &databaseentity.User{
			ID:      2,
			Email:   config.GuestEmail,
			IsGuest: true,
		}
		authRepo.EXPECT().DoUserLogin(ctx, &payloadentity.AttemptLoginPayload{
			Email:    config.GuestEmail,
			Password: config.GuestPassword,
		}).Return(validUser, nil).Times(1)

		res, err := uc.AttemptGuestLogin(ctx)
		assert.Nil(t, err)
		assert.Equal(t, validUser, res)
	})

	t.Run("not guest users", func(t *testing.T) {
		uc, authRepo := initUsecase(ctrl, config.InitTestConfig(map[string]string{"ENABLE_GUEST_LOGIN": "true"}))

		validUser := &databaseentity.User{
			ID:    2,
			Email: config.GuestEmail,
		}
		authRepo.EXPECT().DoUserLogin(ctx, &payloadentity.AttemptLoginPayload{
			Email:    config.GuestEmail,
			Password: config.GuestPassword,
		}).Return(validUser, nil).Times(1)

		res, err := uc.AttemptGuestLogin(ctx)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, herr.Code)
		assert.Equal(t, entity.NonGuestAttemptGuestLogin, herr.Internal.(*entity.InternalError).Code)
	})

	t.Run("disable guest", func(t *testing.T) {
		uc, _ := initUsecase(ctrl, config.InitTestConfig(map[string]string{"ENABLE_GUEST_LOGIN": "false"}))

		res, err := uc.AttemptGuestLogin(ctx)
		assert.Nil(t, res)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, herr.Code)
		assert.Equal(t, entity.AttemptGuestLoginForbidden, herr.Internal.(*entity.InternalError).Code)
	})

}

func Test_GenerateJwtToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := config.InitTestConfig(map[string]string{
		"JWT_SECRET":                "secret",
		"JWT_EXPIRATION_IN_MINUTES": "2",
		"JWT_REMEMBER_IN_DAYS":      "2",
	})
	uc, _ := initUsecase(ctrl, cfg)
	ctx := context.Background()
	t.Run("success, short time token session", func(t *testing.T) {
		exp := time.Now().Add(time.Minute * time.Duration(cfg.JwtExpirationInMinutes())).Unix()

		tokenStr, err := uc.GenerateJwtToken(ctx, &databaseentity.User{
			ID:       1,
			Username: "gendutski",
			Email:    "mvp.firman.darmawan@gmail.com",
		}, false)

		assert.Nil(t, err)

		// parse token
		claims := jwt.MapClaims{}
		jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecret()), nil
		})
		assert.Equal(t, float64(1), claims["id"])
		assert.Equal(t, "gendutski", claims["username"])
		assert.Equal(t, "mvp.firman.darmawan@gmail.com", claims["email"])
		assert.Equal(t, float64(exp), claims["exp"])
	})

	t.Run("success, long time token session", func(t *testing.T) {
		exp := time.Now().Add(time.Hour * 24 * time.Duration(cfg.JwtRememberInDays()))
		tokenStr, err := uc.GenerateJwtToken(ctx, &databaseentity.User{
			ID:       1,
			Username: "gendutski",
			Email:    "mvp.firman.darmawan@gmail.com",
		}, true)

		assert.Nil(t, err)

		// parse token
		claims := jwt.MapClaims{}
		jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtSecret()), nil
		})
		assert.Equal(t, float64(1), claims["id"])
		assert.Equal(t, "gendutski", claims["username"])
		assert.Equal(t, "mvp.firman.darmawan@gmail.com", claims["email"])
		assert.Equal(t, exp.Format("2006-01-02"), time.Unix(int64(claims["exp"].(float64)), 0).Format("2006-01-02"))
	})
}

func Test_ValidateJwtToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, authRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	validUser := &databaseentity.User{
		ID:       1,
		Username: "gendutski",
		Email:    "gendutski@gmail.com",
	}
	validGuest := &databaseentity.User{
		ID:       2,
		Username: "guest",
		Email:    "guest@gmail.com",
		IsGuest:  true,
	}

	t.Run("success guest", func(t *testing.T) {
		authRepo.EXPECT().GetUserByEmail(ctx, "guest@gmail.com").Return(validGuest, nil).Times(1)

		res, err := uc.ValidateJwtToken(ctx, jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":       2,
			"username": "guest",
			"email":    "guest@gmail.com",
			"exp":      time.Now().Add(time.Hour).Unix(),
		}), true)
		assert.Nil(t, err)
		assert.Equal(t, validGuest, res)
	})

	t.Run("success non guest", func(t *testing.T) {
		authRepo.EXPECT().GetUserByEmail(ctx, "gendutski@gmail.com").Return(validUser, nil).Times(1)

		res, err := uc.ValidateJwtToken(ctx, jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":       1,
			"username": "gendutski",
			"email":    "gendutski@gmail.com",
			"exp":      time.Now().Add(time.Hour).Unix(),
		}), true)
		assert.Nil(t, err)
		assert.Equal(t, validUser, res)
	})

	t.Run("token not match", func(t *testing.T) {
		authRepo.EXPECT().GetUserByEmail(ctx, "guest@gmail.com").Return(validGuest, nil).Times(1)

		_, err := uc.ValidateJwtToken(ctx, jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":       1,
			"username": "gendutski",
			"email":    "guest@gmail.com",
			"exp":      time.Now().Add(time.Hour).Unix(),
		}), true)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, herr.Code)
		assert.Equal(t, entity.ValidateTokenFailed, herr.Internal.(*entity.InternalError).Code)
	})

	t.Run("guest forbid", func(t *testing.T) {
		authRepo.EXPECT().GetUserByEmail(ctx, "guest@gmail.com").Return(validGuest, nil).Times(1)

		_, err := uc.ValidateJwtToken(ctx, jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":       2,
			"username": "guest",
			"email":    "guest@gmail.com",
			"exp":      time.Now().Add(time.Hour).Unix(),
		}), false)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusForbidden, herr.Code)
		assert.Equal(t, entity.GuestTokenForbidden, herr.Internal.(*entity.InternalError).Code)
	})
}

func Test_RenewPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, authRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	oldPassword := "0LdF4s#ionP455W0rd"
	encrypted, _ := helper.EncryptPassword(oldPassword)

	user := &databaseentity.User{
		ID:                 1,
		Username:           "gendutski",
		Password:           encrypted,
		PasswordMustChange: true,
	}

	t.Run("success", func(t *testing.T) {
		pass := "Tr!al123#"
		payload := &payloadentity.RenewPasswordPayload{
			Password:        pass,
			ConfirmPassword: pass,
		}
		authRepo.EXPECT().DoRenewPassword(ctx, user, payload).Return(nil).Times(1)

		err := uc.RenewPassword(ctx, user, payload)
		assert.Nil(t, err)
	})

	t.Run("payload empty", func(t *testing.T) {
		err := uc.RenewPassword(ctx, user, &payloadentity.RenewPasswordPayload{})
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, map[string][]*entity.ValidatorMessage{
			"password":        {{Tag: "required"}},
			"confirmPassword": {{Tag: "required"}},
		}, herr.Message)
		assert.Equal(t, entity.ValidatorBadRequest, herr.Internal.(*entity.InternalError).Code)
	})

	t.Run("payload password not match", func(t *testing.T) {
		err := uc.RenewPassword(ctx, user, &payloadentity.RenewPasswordPayload{
			Password:        "trial123",
			ConfirmPassword: "Trial123",
		})
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, map[string][]*entity.ValidatorMessage{
			"confirmPassword": {{Tag: "eqfield", Param: "Password"}},
		}, herr.Message)
		assert.Equal(t, entity.ValidatorBadRequest, herr.Internal.(*entity.InternalError).Code)
	})

	t.Run("failed, password not changed", func(t *testing.T) {
		err := uc.RenewPassword(ctx, user, &payloadentity.RenewPasswordPayload{
			Password:        oldPassword,
			ConfirmPassword: oldPassword,
		})
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, entity.RenewPasswordWithUnchangedPassword, herr.Internal.(*entity.InternalError).Code)
	})
}

func Test_RequestResetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, authRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	validUser := &databaseentity.User{
		Email:              "gendutski@gmail.com",
		PasswordMustChange: true,
	}
	validToken := &databaseentity.PasswordResetToken{
		Token: "token",
	}

	t.Run("success", func(t *testing.T) {
		payload := &payloadentity.RequestResetPayload{Email: "gendutski@gmail.com"}
		authRepo.EXPECT().GetUserByEmail(ctx, "gendutski@gmail.com").Return(validUser, nil).Times(1)
		authRepo.EXPECT().GenerateUserResetPassword(ctx, payload).Return(validToken, nil).Times(1)

		user, token, err := uc.RequestResetPassword(ctx, payload)
		assert.Nil(t, err)
		assert.Equal(t, validToken, token)
		assert.Equal(t, validUser, user)
	})

	t.Run("token failed", func(t *testing.T) {
		payload := &payloadentity.RequestResetPayload{Email: "gendutski@gmail.com"}
		authRepo.EXPECT().GetUserByEmail(ctx, "gendutski@gmail.com").Return(validUser, nil).Times(1)
		authRepo.EXPECT().GenerateUserResetPassword(ctx, payload).Return(nil, errors.New("accident happen")).Times(1)

		user, token, err := uc.RequestResetPassword(ctx, payload)
		assert.NotNil(t, err)
		assert.Nil(t, token)
		assert.Nil(t, user)
	})

	t.Run("user failed", func(t *testing.T) {
		payload := &payloadentity.RequestResetPayload{Email: "gendutski@gmail.com"}
		authRepo.EXPECT().GetUserByEmail(ctx, "gendutski@gmail.com").Return(nil, errors.New("user not found")).Times(1)

		user, token, err := uc.RequestResetPassword(ctx, payload)
		assert.NotNil(t, err)
		assert.Nil(t, token)
		assert.Nil(t, user)
	})

	t.Run("paylod not valid", func(t *testing.T) {
		payload := &payloadentity.RequestResetPayload{Email: "gendutski.com"}

		_, _, err := uc.RequestResetPassword(ctx, payload)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, entity.ValidatorBadRequest, herr.Internal.(*entity.InternalError).Code)
	})
}

func Test_ValidateResetToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, authRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		payload := &payloadentity.ValidateResetTokenPayload{
			Email: "gendutski@gmail.com",
			Token: "token",
		}
		authRepo.EXPECT().ValidateResetToken(ctx, payload).Return(nil, nil).Times(1)

		err := uc.ValidateResetToken(ctx, payload)
		assert.Nil(t, err)
	})

	t.Run("validator failed", func(t *testing.T) {
		payload := &payloadentity.ValidateResetTokenPayload{
			Email: "gendutski.com",
			Token: "token",
		}

		err := uc.ValidateResetToken(ctx, payload)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, entity.ValidatorBadRequest, herr.Internal.(*entity.InternalError).Code)
	})
}

func Test_SubmitResetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	uc, authRepo := initUsecase(ctrl, &config.Config{})
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		payloadToken := &payloadentity.ValidateResetTokenPayload{
			Email: "gendutski@gmail.com",
			Token: "token",
		}
		payload := &payloadentity.ResetPasswordPayload{
			Email:           "gendutski@gmail.com",
			Token:           "token",
			Password:        "new password",
			ConfirmPassword: "new password",
		}
		token := &databaseentity.PasswordResetToken{
			Token: "token",
		}
		validUser := &databaseentity.User{
			Email: "gendutski@gmail.com",
		}

		authRepo.EXPECT().ValidateResetToken(ctx, payloadToken).Return(token, nil).Times(1)
		authRepo.EXPECT().GetUserByEmail(ctx, "gendutski@gmail.com").Return(validUser, nil).Times(1)
		authRepo.EXPECT().DoResetPassword(ctx, gomock.Any(), &databaseentity.PasswordResetToken{}).Return(nil).Times(1)

		_, err := uc.SubmitResetPassword(ctx, payload)
		assert.Nil(t, err)
		assert.Empty(t, token.Email)
	})

	t.Run("validator empty struct", func(t *testing.T) {
		payload := &payloadentity.ResetPasswordPayload{}

		_, err := uc.SubmitResetPassword(ctx, payload)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, entity.ValidatorBadRequest, herr.Internal.(*entity.InternalError).Code)
	})

	t.Run("validator empty TokenEmail struct", func(t *testing.T) {
		payload := &payloadentity.ResetPasswordPayload{
			Password:        "new password",
			ConfirmPassword: "new password",
		}

		_, err := uc.SubmitResetPassword(ctx, payload)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, entity.ValidatorBadRequest, herr.Internal.(*entity.InternalError).Code)
	})

	t.Run("validator TokenEmail struct not valid", func(t *testing.T) {
		payload := &payloadentity.ResetPasswordPayload{
			Email:           "gendutski.com",
			Token:           "token",
			Password:        "new password",
			ConfirmPassword: "new password",
		}

		_, err := uc.SubmitResetPassword(ctx, payload)
		assert.NotNil(t, err)
		herr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, herr.Code)
		assert.Equal(t, entity.ValidatorBadRequest, herr.Internal.(*entity.InternalError).Code)
		assert.Equal(t, map[string][]*entity.ValidatorMessage{
			"email": {{Tag: "email"}},
		}, herr.Message)
	})
}
