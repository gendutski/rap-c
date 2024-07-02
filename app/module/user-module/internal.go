package usermodule

import (
	"context"
	"errors"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// encrypt password given or generated one
func (uc *usecase) generateUserPassword(pass string) (string, error) {
	var err error
	if pass == "" {
		// generate password
		pass, err = helper.GenerateStrongPassword()
		if err != nil {
			return "", &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  http.StatusText(http.StatusInternalServerError),
				Internal: entity.NewInternalError(entity.GeneratePasswordError, err.Error()),
			}
		}
	}

	// encrypt password
	return helper.EncryptPassword(pass)
}

func (uc *usecase) getUserFromJwtClaims(ctx context.Context, claims jwt.MapClaims, guestAccepted bool) (*entity.User, error) {
	// get user
	user, err := uc.userRepo.GetUserByField(ctx, "id", claims[tokenStrID])
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// user not found
			return nil, &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  entity.ValidateTokenUserNotFoundMessage,
				Internal: entity.NewInternalError(entity.ValidateTokenUserNotFound, entity.ValidateTokenUserNotFoundMessage),
			}
		}
		// internal database error
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.ValidateTokenDBError, err.Error()),
		}
	}

	// check user
	if user.Username != claims[tokenStrName].(string) || user.Email != claims[tokenStrEmail].(string) || user.Disabled {
		return nil, &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  entity.ValidateTokenUserNotMatchMessage,
			Internal: entity.NewInternalError(entity.ValidateTokenUserNotMatch, entity.ValidateTokenUserNotMatchMessage),
		}
	}

	// if guest not accepted
	if !guestAccepted && user.IsGuest {
		return nil, &echo.HTTPError{
			Code:     http.StatusForbidden,
			Message:  entity.ValidateTokenGuestNotAcceptedMessage,
			Internal: entity.NewInternalError(entity.ValidateTokenGuestNotAccepted, entity.ValidateTokenGuestNotAcceptedMessage),
		}
	}
	return user, nil
}
