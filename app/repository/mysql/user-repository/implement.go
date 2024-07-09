package userrepository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/helper"
	"rap-c/app/repository/contract"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *repo) Create(ctx context.Context, user *entity.User) error {
	err := r.db.Save(user).Error
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr != nil && mysqlErr.Number == mysqlDuplicateErrorNum {
			code := entity.MysqlDuplicateKeyError
			message := entity.MysqlDuplicateKeyErrorMessage
			if strings.Contains(err.Error(), emailUniqueKeyName) {
				code = entity.UserRepoCreateEmailDuplicate
				message = fmt.Sprintf(entity.UserRepoCreateEmailDuplicateMessage, user.Email)
			} else if strings.Contains(err.Error(), usernameUniqueKeyName) {
				code = entity.UserRepoCreateUsernameDuplicate
				message = fmt.Sprintf(entity.UserRepoCreateUsernameDuplicateMessage, user.Username)
			}

			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  message,
				Internal: entity.NewInternalError(code, err.Error()),
			}
		}
	}
	return &echo.HTTPError{
		Code:     http.StatusInternalServerError,
		Message:  http.StatusText(http.StatusInternalServerError),
		Internal: entity.NewInternalError(entity.UserRepoCreateError, err.Error()),
	}
}

func (r *repo) Update(ctx context.Context, user *entity.User) error {
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

func (r *repo) GetUserByField(ctx context.Context, fieldName string, fieldValue interface{}, notFoundStatus int) (*entity.User, error) {
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

	var result entity.User
	err := r.db.Where(fmt.Sprintf("%s = ?", fieldName), fieldValue).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &echo.HTTPError{
				Code:     notFoundStatus,
				Message:  fmt.Sprintf(entity.UserRepoGetUserByFieldNotFoundMessage, fieldName, fieldValue),
				Internal: entity.NewInternalError(entity.UserRepoGetUserByFieldNotFound, err.Error()),
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

func (r *repo) GetTotalUsersByRequest(ctx context.Context, req *entity.GetUserListRequest) (int64, error) {
	var result int64
	qry := r.renderUsersQuery(req)
	err := qry.Model(entity.User{}).Count(&result).Error
	if err != nil {
		return 0, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoGetTotalUsersByRequestError, err.Error()),
		}
	}
	return result, nil
}

func (r *repo) GetUsersByRequest(ctx context.Context, req *entity.GetUserListRequest) ([]*entity.User, error) {
	var result []*entity.User
	qry := r.renderUsersQuery(req)
	// validate sort
	validSort := map[string]string{
		"userName":  "user_name",
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

func (r *repo) GenerateUserResetPassword(ctx context.Context, email string) (*entity.PasswordResetToken, error) {
	if email == "" {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoGenerateUserResetPasswordError, "email must not empty"),
		}
	}

	token, err := helper.GenerateToken(64)
	if err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoGenerateUserResetPasswordError, err.Error()),
		}
	}

	result := entity.PasswordResetToken{
		Email: email,
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
			Internal: entity.NewInternalError(entity.UserRepoGenerateUserResetPasswordError, err.Error()),
		}
	}
	return &result, nil
}

func (r *repo) ValidateResetToken(ctx context.Context, email string, token string) (*entity.PasswordResetToken, error) {
	if email == "" || token == "" {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoValidateResetTokenError, "email or token must not empty"),
		}
	}

	var result entity.PasswordResetToken
	// get token from db
	err := r.db.Where("email = ?", email).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message := fmt.Sprintf(entity.UserRepoValidateResetTokenNotFoundMessage, email, token)
			return nil, &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  message,
				Internal: entity.NewInternalError(entity.UserRepoValidateResetTokenNotFound, message),
			}
		}
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoValidateResetTokenError, err.Error()),
		}
	}
	// validate token & expired date
	if result.Token != token || time.Now().After(result.UpdatedAt.Add(resetTokenExpiration)) {
		message := fmt.Sprintf(entity.UserRepoValidateResetTokenNotFoundMessage, email, token)
		return nil, &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  message,
			Internal: entity.NewInternalError(entity.UserRepoValidateResetTokenNotFound, message),
		}
	}
	return &result, nil
}

func (r *repo) ResetPassword(ctx context.Context, user *entity.User, reset *entity.PasswordResetToken) (err error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  http.StatusText(http.StatusInternalServerError),
				Internal: entity.NewInternalError(entity.UserRepoResetPasswordError, fmt.Sprint(r)),
			}
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		err = &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoResetPasswordError, err.Error()),
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
			Internal: entity.NewInternalError(entity.UserRepoResetPasswordError, err.Error()),
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
			Internal: entity.NewInternalError(entity.UserRepoResetPasswordError, err.Error()),
		}
		return
	}

	err = tx.Commit().Error
	if err != nil {
		err = &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UserRepoResetPasswordError, err.Error()),
		}
	}
	return
}

func (r *repo) renderUsersQuery(req *entity.GetUserListRequest) *gorm.DB {
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
		if req.Show == entity.RequestShowActive {
			qry = qry.Where("disabled = ?", false)
		} else if req.Show == entity.RequestShowNotActive {
			qry = qry.Where("disabled = ?", true)
		}
	}
	if req.GuestOnly {
		qry = qry.Where("is_guest = ?", true)
	}
	return qry
}
