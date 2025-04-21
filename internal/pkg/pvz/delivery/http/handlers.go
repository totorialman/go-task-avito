package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"log/slog"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/totorialman/go-task-avito/internal/pkg/metrics"
	"github.com/totorialman/go-task-avito/internal/pkg/pvz/usecase"
	"github.com/totorialman/go-task-avito/internal/pkg/utils/log"
	"github.com/totorialman/go-task-avito/models"
	"github.com/totorialman/go-task-avito/restapi/operations"
)

type PVZHandler struct {
	usecase *usecase.PVZUsecase
	mt      *metrics.ProductMetrics
}

func NewPVZHandler(uc *usecase.PVZUsecase, mt *metrics.ProductMetrics) *PVZHandler {
	return &PVZHandler{usecase: uc, mt: mt}
}

func (h *PVZHandler) HandleCreatePVZ(params operations.PostPvzParams) middleware.Responder {
	logger := log.GetLoggerFromContext(params.HTTPRequest.Context()).With(slog.String("func", log.GetFuncName()))
	req := params.Body

	if req == nil {
		log.LogHandlerError(logger, errors.New("request body is nil"), http.StatusInternalServerError)
		return operations.NewPostPvzBadRequest().WithPayload(&models.Error{
			Message: swag.String("internal server error"),
		})
	}
	if req.City == nil {
		log.LogHandlerError(logger, errors.New("city is nil"), http.StatusBadRequest)
		return operations.NewPostPvzBadRequest().WithPayload(&models.Error{
			Message: swag.String("city is required"),
		})
	}

	pvz, err := h.usecase.CreatePVZ(params.HTTPRequest.Context(), *req.City, req.ID, req.RegistrationDate)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("CreatePVZ error: %w", err), http.StatusBadRequest)
		if err.Error() == "only moderators can create PVZ" {
			return operations.NewPostPvzForbidden().WithPayload(&models.Error{
				Message: swag.String(err.Error()),
			})
		}
		return operations.NewPostPvzBadRequest().WithPayload(&models.Error{
			Message: swag.String(err.Error()),
		})
	}

	if pvz == nil {
		log.LogHandlerError(logger, errors.New("pvz is nil after creation"), http.StatusInternalServerError)
		return operations.NewPostPvzBadRequest().WithPayload(&models.Error{
			Message: swag.String("failed to create PVZ"),
		})
	}

	if h.mt != nil {
		h.mt.IncreaseHitsPVZTotal()
	} else {
		log.LogHandlerError(logger, errors.New("metrics collector is nil"), http.StatusInternalServerError)
	}

	return operations.NewPostPvzCreated().WithPayload(pvz)
}

func (h *PVZHandler) HandleCreateReception(params operations.PostReceptionsParams) middleware.Responder {
	logger := log.GetLoggerFromContext(params.HTTPRequest.Context()).With(slog.String("func", log.GetFuncName()))
	req := params.Body

	activeReception, _, err := h.usecase.GetActiveReception(params.HTTPRequest.Context(), *req.PvzID)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("GetActiveReception error: %w", err), http.StatusBadRequest)
		return operations.NewPostReceptionsBadRequest().WithPayload(&models.Error{
			Message: swag.String("Ошибка при проверке активной приемки"),
		})
	}

	if activeReception {
		log.LogHandlerError(logger, errors.New("previous reception not closed"), http.StatusBadRequest)
		return operations.NewPostReceptionsBadRequest().WithPayload(&models.Error{
			Message: swag.String("Невозможно создать новую приемку, так как предыдущая не закрыта"),
		})
	}

	reception, err := h.usecase.CreateReception(params.HTTPRequest.Context(), *req.PvzID)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("CreateReception error: %w", err), http.StatusBadRequest)
		return operations.NewPostReceptionsBadRequest().WithPayload(&models.Error{
			Message: swag.String("Ошибка при создании приемки"),
		})
	}
	if h.mt != nil {
		h.mt.IncreaseHitsReTotal()
	} else {
		log.LogHandlerError(logger, errors.New("metrics collector is nil"), http.StatusInternalServerError)
	}
	return operations.NewPostReceptionsCreated().WithPayload(reception)
}

func (h *PVZHandler) HandleAddProductToReception(params operations.PostProductsParams) middleware.Responder {
	logger := log.GetLoggerFromContext(params.HTTPRequest.Context()).With(slog.String("func", log.GetFuncName()))
	req := params.Body

	activeReception, _, err := h.usecase.GetActiveReception(params.HTTPRequest.Context(), *req.PvzID)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("GetActiveReception error: %w", err), http.StatusBadRequest)
		return operations.NewPostProductsBadRequest().WithPayload(&models.Error{
			Message: swag.String("Ошибка при проверке активной приемки"),
		})
	}

	if !activeReception {
		log.LogHandlerError(logger, errors.New("no active reception for PVZ"), http.StatusBadRequest)
		return operations.NewPostProductsBadRequest().WithPayload(&models.Error{
			Message: swag.String("Нет активной приемки для данного ПВЗ"),
		})
	}

	productType := ""
	if req.Type != nil {
		productType = *req.Type
	}

	product, err := h.usecase.AddProductToReception(params.HTTPRequest.Context(), *req.PvzID, productType)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("AddProductToReception error: %w", err), http.StatusBadRequest)
		return operations.NewPostProductsBadRequest().WithPayload(&models.Error{
			Message: swag.String("Ошибка при добавлении товара в приемку"),
		})
	}
	if h.mt != nil {
		h.mt.IncreaseHitsProductTotal()
	} else {
		log.LogHandlerError(logger, errors.New("metrics collector is nil"), http.StatusInternalServerError)
	}
	return operations.NewPostProductsCreated().WithPayload(product)
}

func (h *PVZHandler) HandleDeleteLastProduct(params operations.PostPvzPvzIDDeleteLastProductParams) middleware.Responder {
	logger := log.GetLoggerFromContext(params.HTTPRequest.Context()).With(slog.String("func", log.GetFuncName()))

	err := h.usecase.DeleteLastProductFromReception(params.HTTPRequest.Context(), params.PvzID)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("DeleteLastProductFromReception error: %w", err), http.StatusBadRequest)
		if err.Error() == "lol" {
			return operations.NewPostPvzPvzIDDeleteLastProductForbidden().WithPayload(&models.Error{
				Message: swag.String(err.Error()),
			})
		}
		return operations.NewPostPvzPvzIDDeleteLastProductBadRequest().WithPayload(&models.Error{
			Message: swag.String(err.Error()),
		})
	}

	return operations.NewPostPvzPvzIDDeleteLastProductOK()
}

func (h *PVZHandler) HandleCloseLastReception(params operations.PostPvzPvzIDCloseLastReceptionParams) middleware.Responder {
	logger := log.GetLoggerFromContext(params.HTTPRequest.Context()).With(slog.String("func", log.GetFuncName()))

	reception, err := h.usecase.CloseLastReception(params.HTTPRequest.Context(), params.PvzID)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("CloseLastReception error: %w", err), http.StatusBadRequest)
		if err.Error() == "приемка уже закрыта или отсутствует" {
			return operations.NewPostPvzPvzIDCloseLastReceptionBadRequest().WithPayload(&models.Error{
				Message: swag.String(err.Error()),
			})
		}
		return operations.NewPostPvzPvzIDCloseLastReceptionForbidden().WithPayload(&models.Error{
			Message: swag.String("Нет прав для закрытия приемки"),
		})
	}

	return operations.NewPostPvzPvzIDCloseLastReceptionOK().WithPayload(reception)
}

func (h *PVZHandler) HandleGetPVZs(params operations.GetPvzParams) middleware.Responder {
	logger := log.GetLoggerFromContext(params.HTTPRequest.Context()).With(slog.String("func", log.GetFuncName()))

	var startDate, endDate *time.Time

	if params.StartDate != nil {
		t, err := time.Parse(time.RFC3339, params.StartDate.String())
		if err != nil {
			log.LogHandlerError(logger, fmt.Errorf("Error parsing start date: %w", err), http.StatusBadRequest)
			return operations.NewGetPvzOK()
		}
		startDate = &t
	}

	if params.EndDate != nil {
		t, err := time.Parse(time.RFC3339, params.EndDate.String())
		if err != nil {
			log.LogHandlerError(logger, fmt.Errorf("Error parsing end date: %w", err), http.StatusBadRequest)
			return operations.NewGetPvzOK()
		}
		endDate = &t
	}

	page := int(*params.Page)
	limit := int(*params.Limit)

	pvzs, err := h.usecase.GetPVZs(params.HTTPRequest.Context(), startDate, endDate, page, limit)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("GetPVZs error: %w", err), http.StatusInternalServerError)
		return operations.NewGetPvzOK()
	}

	return operations.NewGetPvzOK().WithPayload(pvzs)
}
