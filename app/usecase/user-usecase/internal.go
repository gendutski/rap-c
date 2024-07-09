package userusecase

import (
	"context"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// encrypt password given or generated one
func (uc *usecase) generateUserPassword(pass string) (password string, encryptPassword string, err error) {
	if pass == "" {
		// generate password
		password, err = helper.GenerateStrongPassword()
		if err != nil {
			err = &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  http.StatusText(http.StatusInternalServerError),
				Internal: entity.NewInternalError(entity.UserUsecaseGenerateStrongPasswordError, err.Error()),
			}
			return
		}
	} else {
		password = pass
	}

	// encrypt password
	encryptPassword, err = helper.EncryptPassword(password)
	if err != nil {
		err = &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserUsecaseEncryptPasswordError, err.Error()),
		}
	}
	return
}

func (uc *usecase) getUserFromJwtClaims(ctx context.Context, claims jwt.MapClaims, guestAccepted bool) (*entity.User, error) {
	// get user
	user, err := uc.userRepo.GetUserByField(ctx, "id", claims[tokenStrID], http.StatusUnauthorized)
	if err != nil {
		return nil, err
	}

	// check user
	if user.Username != claims[tokenStrName].(string) || user.Email != claims[tokenStrEmail].(string) || user.Disabled {
		return nil, &echo.HTTPError{
			Code:     http.StatusUnauthorized,
			Message:  entity.UserUsecaseGetUserFromJwtClaimsUnauthorizedMessage,
			Internal: entity.NewInternalError(entity.UserUsecaseGetUserFromJwtClaimsUnauthorized, entity.UserUsecaseGetUserFromJwtClaimsUnauthorizedMessage),
		}
	}

	// if guest not accepted
	if !guestAccepted && user.IsGuest {
		return nil, &echo.HTTPError{
			Code:     http.StatusForbidden,
			Message:  entity.UserUsecaseGetUserFromJwtClaimsForbidGuestMessage,
			Internal: entity.NewInternalError(entity.UserUsecaseGetUserFromJwtClaimsForbidGuest, entity.UserUsecaseGetUserFromJwtClaimsForbidGuestMessage),
		}
	}
	return user, nil
}
