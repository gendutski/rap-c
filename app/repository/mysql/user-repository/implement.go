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

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	minQueryLimit          int    = 10
	maxQueryLimit          int    = 100
	mysqlDuplicateErrorNum uint16 = 1062
	emailUniqueKeyName     string = "users.uni_users_email"
	usernameUniqueKeyName  string = "users.uni_users_username"
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
				code = entity.CreateUserEmailDuplicate
				message = fmt.Sprintf(entity.CreateUserEmailDuplicateMessage, user.Email)
			} else if strings.Contains(err.Error(), usernameUniqueKeyName) {
				code = entity.CreateUserUsernameDuplicate
				message = fmt.Sprintf(entity.CreateUserUsernameDuplicateMessage, user.Username)
			}

			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  message,
				Internal: entity.NewInternalError(code, err.Error()),
			}
		}
	}
	return err
}

func (r *repo) Update(ctx context.Context, user *entity.User) error {
	if user.ID == 0 {
		return errors.New("data not found, empty primary key")
	}
	return r.db.Save(user).Error
}

func (r *repo) GetUserByField(ctx context.Context, fieldName string, fieldValue interface{}) (*entity.User, error) {
	acceptedFields := map[string]bool{
		"id":       true,
		"username": true,
		"email":    true,
	}
	if !acceptedFields[fieldName] {
		return nil, errors.New("invalid field name for query")
	}

	var result entity.User
	err := r.db.Where(fmt.Sprintf("%s = ?", fieldName), fieldValue).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *repo) GetTotalUsersByRequest(ctx context.Context, req *entity.GetUserListRequest) (int64, error) {
	var result int64
	qry := r.renderUsersQuery(req)
	err := qry.Model(entity.User{}).Count(&result).Error
	if err != nil {
		return 0, err
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
		return nil, err
	}
	return result, nil
}

func (r *repo) GenerateUserResetPassword(ctx context.Context, email string) (*entity.PasswordResetToken, error) {
	token, err := helper.GenerateToken(64)
	if err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.GenerateTokenError, err.Error()),
		}
	}

	result := entity.PasswordResetToken{
		Email: email,
		Token: token,
	}
	err = r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{"token"}),
		}).
		Create(&result).Error
	if err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.RepoGenerateUserResetPasswordError, err.Error()),
		}
	}
	return &result, nil
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
