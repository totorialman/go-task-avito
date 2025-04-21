package pvz

import (
	"context"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/totorialman/go-task-avito/models"
	"github.com/totorialman/go-task-avito/restapi/operations"
)

type PVZRepository interface {
	CreatePVZ(ctx context.Context, pvz *models.PVZ) error
	CreateReception(ctx context.Context, reception *models.Reception) error
	GetActiveReception(ctx context.Context, pvzID strfmt.UUID) (bool, *models.Reception, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	DeleteLastProduct(ctx context.Context, receptionID strfmt.UUID) error
	UpdateReceptionStatus(ctx context.Context, reception models.Reception) error
	GetCloseReception(ctx context.Context, pvzID strfmt.UUID) (bool, *models.Reception, error)
	GetPVZsWithReceptions(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]map[string]interface{}, error)
}

type PVZUsecase interface {
	CreatePVZ(ctx context.Context, city string, id strfmt.UUID, date strfmt.DateTime) (*models.PVZ, error) 
	CreateReception(ctx context.Context, pvzID, createdBy strfmt.UUID) (*models.Reception, error)
	GetPVZs(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*operations.GetPvzOKBodyItems0, error)
	CloseLastReception(ctx context.Context, pvzID strfmt.UUID) (*models.Reception, error) 
	DeleteLastProductFromReception(ctx context.Context, pvzID strfmt.UUID) error 
	AddProductToReception(ctx context.Context, pvzID strfmt.UUID, productType string) (*models.Product, error)
	GetActiveReception(ctx context.Context, pvzID strfmt.UUID) (*models.Reception, error)
}
