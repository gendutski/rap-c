package unitrepository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"rap-c/app/entity"
	databaseentity "rap-c/app/entity/database-entity"
	payloadentity "rap-c/app/entity/payload-entity"
	"rap-c/app/repository/contract"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	minQueryLimit             int    = 10
	maxQueryLimit             int    = 100
	mysqlDuplicateErrorNum    uint16 = 1062
	mysqlForeignKeyReferences uint16 = 1451
)

type repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) contract.UnitRepository {
	return &repo{db}
}

func (r *repo) Create(ctx context.Context, unit *databaseentity.Unit) error {
	err := r.db.Save(unit).Error
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr != nil && mysqlErr.Number == mysqlDuplicateErrorNum {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  fmt.Sprintf(entity.CreateUnitNameDuplicateMessage, unit.Name),
				Internal: entity.NewInternalError(entity.CreateUnitNameDuplicate, err.Error()),
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UnitRepoCreateError, err.Error()),
		}
	}
	return nil
}

func (r *repo) GetUnitByName(ctx context.Context, name string) (*databaseentity.Unit, error) {
	var result databaseentity.Unit
	err := r.db.Where("name = ?", name).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  fmt.Sprintf(entity.UnitNameNotFoundMessage, name),
				Internal: entity.NewInternalError(entity.UnitNameNotFound, err.Error()),
			}
		}
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UnitRepoGetUnitByNameError, err.Error()),
		}
	}
	return &result, nil
}

func (r *repo) GetTotalUnitsByRequest(ctx context.Context, req *payloadentity.GetUnitListRequest) (int64, error) {
	var result int64
	qry := r.renderUnitsQuery(req)
	err := qry.Model(databaseentity.Unit{}).Count(&result).Error
	if err != nil {
		return 0, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UnitRepoGetTotalUnitsByRequestError, err.Error()),
		}
	}
	return result, nil
}

func (r *repo) GetUnitsByRequest(ctx context.Context, req *payloadentity.GetUnitListRequest) ([]*databaseentity.Unit, error) {
	var result []*databaseentity.Unit
	qry := r.renderUnitsQuery(req)
	// validate sort
	validSort := map[string]string{
		"name":      "name",
		"createdAt": "creted_at",
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
			Internal: entity.NewInternalError(entity.UnitRepoGetUnitsByRequestError, err.Error()),
		}
	}
	return result, nil
}

func (r *repo) Delete(ctx context.Context, unit *databaseentity.Unit) error {
	err := r.db.Delete(unit).Error
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr != nil && mysqlErr.Number == mysqlForeignKeyReferences {
			return &echo.HTTPError{
				Code:     http.StatusForbidden,
				Message:  entity.DeleteUsedUnitForbiddenMessage,
				Internal: entity.NewInternalError(entity.DeleteUsedUnitForbidden, err.Error()),
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.UnitRepoDeleteError, err.Error()),
		}
	}
	return nil
}

func (r *repo) renderUnitsQuery(req *payloadentity.GetUnitListRequest) *gorm.DB {
	qry := r.db
	if req.Name != "" {
		qry = qry.Where("name like ?", fmt.Sprintf("%%%s%%", req.Name))
	}
	return qry
}
