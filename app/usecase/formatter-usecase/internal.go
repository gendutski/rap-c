package formatterusecase

import (
	databaseentity "rap-c/app/entity/database-entity"
	responseentity "rap-c/app/entity/response-entity"
)

func (uc *usecase) formatUser(user *databaseentity.User, mapUsers map[int]string) *responseentity.UserResponse {
	createdBy, ok := mapUsers[user.CreatedByDB]
	if !ok {
		createdBy = defaultUserUsername
	}
	updatedBy, ok := mapUsers[user.UpdatedByDB]
	if !ok {
		updatedBy = defaultUserUsername
	}
	return &responseentity.UserResponse{
		Username:           user.Username,
		FullName:           user.FullName,
		Email:              user.Email,
		PasswordMustChange: user.PasswordMustChange,
		Disabled:           user.Disabled,
		IsGuest:            user.IsGuest,
		CreatedAt:          user.CreatedAt,
		CreatedBy:          createdBy,
		UpdatedAt:          user.UpdatedAt,
		UpdatedBy:          updatedBy,
	}
}
