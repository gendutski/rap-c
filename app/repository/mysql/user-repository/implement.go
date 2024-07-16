package userrepository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"rap-c/app/entity"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
	"rap-c/app/repository/contract"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	minQueryLimit          int           = 10
	maxQueryLimit          int           = 100
	mysqlDuplicateErrorNum uint16        = 1062
	emailUniqueKeyName     string        = "users.uni_users_email"
	usernameUniqueKeyName  string        = "users.uni_users_username"
	resetTokenExpiration   time.Duration = time.Hour
)

type repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) contract.UserRepository {
	return &repo{db}
}

func (r *repo) Create(ctx context.Context, user *databaseentity.User) error {
	err := r.db.Save(user).Error
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr != nil && mysqlErr.Number == mysqlDuplicateErrorNum {
			if strings.Contains(err.Error(), emailUniqueKeyName) {
				return &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  fmt.Sprintf(entity.CreateUserEmailDuplicateMessage, user.Email),
					Internal: entity.NewInternalError(entity.CreateUserEmailDuplicate, err.Error()),
				}

			} else if strings.Contains(err.Error(), usernameUniqueKeyName) {
				return &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  fmt.Sprintf(entity.CreateUserUsernameDuplicateMessage, user.Username),
					Internal: entity.NewInternalError(entity.CreateUserUsernameDuplicate, err.Error()),
				}
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoCreateError, err.Error()),
		}
	}
	return nil
}

func (r *repo) Update(ctx context.Context, user *databaseentity.User) error {
	if user.ID == 0 {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoUpdateError, "data not found, empty primary key"),
		}
	}
	err := r.db.Save(user).Error
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoUpdateError, err.Error()),
		}
	}
	return nil
}

func (r *repo) GetUserByField(ctx context.Context, fieldName string, fieldValue interface{}, notFoundStatus int) (*databaseentity.User, error) {
	acceptedFields := map[string]bool{
		"id":       true,
		"username": true,
		"email":    true,
	}
	if !acceptedFields[fieldName] {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoGetUserByFieldError, "invalid field name for query"),
		}
	}

	var result databaseentity.User
	err := r.db.Where(fmt.Sprintf("%s = ?", fieldName), fieldValue).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &echo.HTTPError{
				Code:     notFoundStatus,
				Message:  fmt.Sprintf(entity.SearchSingleUserNotFOundMessage, fieldName, fieldValue),
				Internal: entity.NewInternalError(entity.SearchSingleUserNotFOund, err.Error()),
			}
		}
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoGetUserByFieldError, err.Error()),
		}
	}
	return &result, nil
}

func (r *repo) GetTotalUsersByRequest(ctx context.Context, req *payloadentity.GetUserListRequest) (int64, error) {
	var result int64
	qry := r.renderUsersQuery(req)
	err := qry.Model(databaseentity.User{}).Count(&result).Error
	if err != nil {
		return 0, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoGetTotalUsersByRequestError, err.Error()),
		}
	}
	return result, nil
}

func (r *repo) GetUsersByRequest(ctx context.Context, req *payloadentity.GetUserListRequest) ([]*databaseentity.User, error) {
	var result []*databaseentity.User
	qry := r.renderUsersQuery(req)
	// validate sort
	validSort := map[string]string{
		"username":  "user_name",
		"fullName":  "full_name",
		"email":     "email",
		"role":      "role",
		"createdAt": "creted_at",
		"updatedAt": "updated_at",
	}
	sort := validSort[req.SortField]
	if sort == "" {
		req.SortField = "createdAt"
		sort = "created_at"
	}
	order := "asc"
	if req.DescendingOrder {
		order = "desc"
	}
	// validate limit
	if req.Limit < minQueryLimit {
		req.Limit = minQueryLimit
	} else if req.Limit > maxQueryLimit {
		req.Limit = maxQueryLimit
	}
	err := qry.Order(fmt.Sprintf("%s %s", sort, order)).
		Limit(req.Limit).
		Offset(req.Page.GetOffset(req.Limit)).
		Find(&result).
		Error
	if err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoGetUsersByRequestError, err.Error()),
		}
	}
	return result, nil
}

func (r *repo) renderUsersQuery(req *payloadentity.GetUserListRequest) *gorm.DB {
	qry := r.db
	if req.UserName != "" {
		qry = qry.Where("user_name = ?", req.UserName)
	}
	if req.Email != "" {
		qry = qry.Where("email like ?", fmt.Sprintf("%%%s%%", req.Email))
	}
	if req.FullName != "" {
		qry = qry.Where("name like ?", fmt.Sprintf("%%%s%%", req.FullName))
	}
	if req.Show != "" {
		if req.Show == databaseentity.RequestShowActive {
			qry = qry.Where("disabled = ?", false)
		} else if req.Show == databaseentity.RequestShowNotActive {
			qry = qry.Where("disabled = ?", true)
		}
	}
	if req.GuestOnly {
		qry = qry.Where("is_guest = ?", true)
	}
	return qry
}
