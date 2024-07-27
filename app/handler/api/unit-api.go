package api

import (
	"fmt"
	"net/http"
	"rap-c/app/entity"
	payloadentity "rap-c/app/entity/payload-entity"
	responseentity "rap-c/app/entity/response-entity"
	"rap-c/app/handler"
	"rap-c/app/usecase/contract"
	"rap-c/config"

	"github.com/labstack/echo/v4"
)

type UnitAPI interface {
	// get unit list
	GetUnitList(e echo.Context) error
	// get total unit list
	GetTotalUnitList(e echo.Context) error
	// create unit
	Create(e echo.Context) error
	// delete not used unit
	Delete(e echo.Context) error
}

func NewUnitHandler(cfg *config.Config, router *config.Route,
	unitUsecase contract.UnitUsecase, formatterUsecase contract.FormatterUsecase) UnitAPI {
	return &unitHandler{
		cfg:              cfg,
		router:           router,
		unitUsecase:      unitUsecase,
		formatterUsecase: formatterUsecase,
		BaseHandler:      handler.NewBaseHandler(cfg, router),
	}
}

type unitHandler struct {
	cfg              *config.Config
	router           *config.Route
	unitUsecase      contract.UnitUsecase
	formatterUsecase contract.FormatterUsecase
	BaseHandler      *handler.BaseHandler
}

func (h *unitHandler) GetUnitList(e echo.Context) error {
	req := new(payloadentity.GetUnitListRequest)
	err := e.Bind(req)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AllHandlerBindError, fmt.Sprintf("unit-api.GetUnitList bind error: %v", err)),
		}
	}
	ctx := e.Request().Context()

	units, err := h.unitUsecase.GetUnitList(ctx, req)
	if err != nil {
		return err
	}

	resp, err := h.formatterUsecase.FormatUnits(ctx, units)
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, &responseentity.GetUnitListResponse{
		Units:   resp,
		Request: req,
	})
}

func (h *unitHandler) GetTotalUnitList(e echo.Context) error {
	req := new(payloadentity.GetUnitListRequest)
	err := e.Bind(req)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AllHandlerBindError, fmt.Sprintf("unit-api.GetTotalUnitList bind error: %v", err)),
		}
	}
	ctx := e.Request().Context()

	total, err := h.unitUsecase.GetTotalUnitList(ctx, req)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"total":   total,
		"request": req,
	})
}

func (h *unitHandler) Create(e echo.Context) error {
	payload := new(payloadentity.CreateUnitPayload)
	err := e.Bind(payload)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AllHandlerBindError, fmt.Sprintf("unit-api.Create bind error: %v", err)),
		}
	}
	ctx := e.Request().Context()

	// get author
	author, err := h.BaseHandler.GetAuthor(e)
	if err != nil {
		return err
	}

	// create unit
	unit, err := h.unitUsecase.Create(ctx, payload, author)
	if err != nil {
		return err
	}

	resp, err := h.formatterUsecase.FormatUnit(ctx, unit, map[int]string{author.ID: author.Username})
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, resp)
}

func (h *unitHandler) Delete(e echo.Context) error {
	payload := new(payloadentity.DeleteUnitPayload)
	err := e.Bind(payload)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  http.StatusText(http.StatusInternalServerError),
			Internal: entity.NewInternalError(entity.AllHandlerBindError, fmt.Sprintf("unit-api.Delete bind error: %v", err)),
		}
	}
	ctx := e.Request().Context()

	// delete unit
	err = h.unitUsecase.Delete(ctx, payload)
	if err != nil {
		return err
	}
	return e.JSON(http.StatusNoContent, map[string]interface{}{
		"status": "ok",
	})
}
