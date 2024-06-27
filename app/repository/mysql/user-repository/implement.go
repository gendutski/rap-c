package userrepository

import (
	"context"
	"fmt"
	"net/http"
	"rap-c/app/entity"
	"rap-c/app/repository/contract"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
