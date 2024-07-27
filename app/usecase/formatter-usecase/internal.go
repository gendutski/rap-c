package formatterusecase

import (
	databaseentity "rap-c/app/entity/database-entity"
	responseentity "rap-c/app/entity/response-entity"
)

const (
	defaultUserUsername string = "SYSTEM"
	errorUserUsername   string = "DELETED USER"
	createdByString     string = "createdBy"
	updatedByString     string = "updatedBy"
)

func (uc *usecase) formatUser(user *databaseentity.User, mapUsers map[int]string) *responseentity.UserResponse {
	createdBy, ok := mapUsers[user.CreatedBy]
	if !ok {
		createdBy = defaultUserUsername
	}
	updatedBy, ok := mapUsers[user.UpdatedBy]
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

func (uc *usecase) formatUnit(unit *databaseentity.Unit, mapUsers map[int]string) *responseentity.UnitResponse {
	createdBy, ok := mapUsers[unit.CreatedBy]
	if !ok {
		createdBy = errorUserUsername
	}

	return &responseentity.UnitResponse{
		Name:      unit.Name,
		CreatedAt: unit.CreatedAt,
		CreatedBy: createdBy,
	}
}
