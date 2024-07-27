package authrepository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"rap-c/app/entity"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
	"rap-c/app/helper"
	"rap-c/app/repository/contract"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	resetTokenExpiration time.Duration = time.Hour
)

type repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) contract.AuthRepository {
	return &repo{db}
}

func (r *repo) DoUserLogin(ctx context.Context, payload *payloadentity.AttemptLoginPayload) (*databaseentity.User, error) {
	var user databaseentity.User

	// get from database
	err := r.db.Where("email = ?", payload.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  fmt.Sprintf(entity.AttemptLoginFailedMessage),
				Internal: entity.NewInternalError(entity.AttemptLoginFailed, err.Error()),
			}
		}
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AuthRepoGetUserLoginError, err.Error()),
		}
	}

	// validate password
	if !helper.ValidateEncryptedPassword(user.Password, payload.Password) {
		return nil, &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  fmt.Sprintf(entity.AttemptLoginFailedMessage),
			Internal: entity.NewInternalError(entity.AttemptLoginFailed, "password not match"),
		}
	}

	// validate active status
	if user.Disabled {
		return nil, &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  entity.AttemptLoginUserDeactivatedMessage,
			Internal: entity.NewInternalError(entity.AttemptLoginUserDeactivated, entity.AttemptLoginUserDeactivatedMessage),
		}
	}

	return &user, nil
}

func (r *repo) DoRenewPassword(ctx context.Context, user *databaseentity.User, payload *payloadentity.RenewPasswordPayload) error {
	if user.ID == 0 {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoUpdateError, "data not found, empty primary key"),
		}
	}

	// encrypt password
	encryptPassword, err := helper.EncryptPassword(payload.Password)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.HelperEncryptPasswordError, err.Error()),
		}
	}

	user.Password = encryptPassword
	user.PasswordMustChange = false
	user.UpdatedBy = user.ID
	err = r.db.Save(user).Error
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AuthRepoDoRenewPasswordError, err.Error()),
		}
	}
	return nil
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*databaseentity.User, error) {
	var user databaseentity.User

	// get from database
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  fmt.Sprintf(entity.SearchSingleUserNotFOundMessage, "email", email),
				Internal: entity.NewInternalError(entity.SearchSingleUserNotFOund, err.Error()),
			}
		}
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AuthRepoGetUserLoginError, err.Error()),
		}
	}
	return &user, nil
}

func (r *repo) GenerateUserResetPassword(ctx context.Context, payload *payloadentity.RequestResetPayload) (*databaseentity.PasswordResetToken, error) {
	token, err := helper.GenerateToken(64)
	if err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.HelperGenerateTokenError, err.Error()),
		}
	}

	result := databaseentity.PasswordResetToken{
		Email: payload.Email,
		Token: token,
	}
	err = r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{"token", "updated_at"}),
		}).
		Create(&result).Error
	if err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AuthRepoGenerateUserResetPasswordError, err.Error()),
		}
	}
	return &result, nil
}

func (r *repo) ValidateResetToken(ctx context.Context, payload *payloadentity.ValidateResetTokenPayload) (*databaseentity.PasswordResetToken, error) {
	var result databaseentity.PasswordResetToken

	// get token from db
	err := r.db.Where("email = ?", payload.Email).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  entity.ResetPasswordRequestNotFoundMessage,
				Internal: entity.NewInternalError(entity.ResetPasswordRequestNotFound, err.Error()),
			}
		}
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AuthRepoValidateResetTokenError, err.Error()),
		}
	}
	// validate token & expired date
	if result.Token == "" || result.Token != payload.Token || time.Now().After(result.UpdatedAt.Add(resetTokenExpiration)) {
		return nil, &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  entity.ResetPasswordRequestNotFoundMessage,
			Internal: entity.NewInternalError(entity.ResetPasswordRequestNotFound, "token expired or not match"),
		}
	}
	return &result, nil
}

func (r *repo) DoResetPassword(ctx context.Context, user *databaseentity.User, reset *databaseentity.PasswordResetToken) (err error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  http.StatusText(http.StatusInternalServerError),
				Internal: entity.NewInternalError(entity.AuthRepoDoResetPasswordError, fmt.Sprint(r)),
			}
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		err = &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AuthRepoDoResetPasswordError, err.Error()),
		}
		return
	}

	// save user
	err = tx.Save(user).Error
	if err != nil {
		tx.Rollback()
		err = &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AuthRepoDoResetPasswordError, err.Error()),
		}
		return
	}

	// save reset password
	err = tx.Save(reset).Error
	if err != nil {
		tx.Rollback()
		err = &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AuthRepoDoResetPasswordError, err.Error()),
		}
		return
	}

	err = tx.Commit().Error
	if err != nil {
		err = &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AuthRepoDoResetPasswordError, err.Error()),
		}
	}
	return
}
