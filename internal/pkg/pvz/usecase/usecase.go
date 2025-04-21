package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	"github.com/totorialman/go-task-avito/internal/pkg/pvz"
	"github.com/totorialman/go-task-avito/models"
	"github.com/totorialman/go-task-avito/restapi/operations"
)

type PVZUsecase struct {
	repo pvz.PVZRepository
}

func NewPVZUsecase(repo pvz.PVZRepository) *PVZUsecase {
	return &PVZUsecase{repo: repo}
}

func (u *PVZUsecase) CreatePVZ(ctx context.Context, city string, id strfmt.UUID, date strfmt.DateTime) (*models.PVZ, error) {

	cityPtr := new(string)
	*cityPtr = city

	pvz := &models.PVZ{
		ID:               id,
		City:             cityPtr,
		RegistrationDate: date,
	}

	if err := u.repo.CreatePVZ(ctx, pvz); err != nil {
		return nil, err
	}

	return pvz, nil
}

func (u *PVZUsecase) CreateReception(ctx context.Context, pvzID strfmt.UUID) (*models.Reception, error) {
	activeReception, _, _ := u.repo.GetActiveReception(ctx, pvzID)
	if activeReception {
		return nil, errors.New("невозможно создать новую приемку, так как предыдущая не закрыта")
	}
	id := uuid.New()
	reception := &models.Reception{
		ID:     strfmt.UUID(id.String()),
		PvzID:  &pvzID,
		Status: swag.String("in_progress"),
	}

	err := u.repo.CreateReception(ctx, reception)
	if err != nil {
		return nil, err
	}

	return reception, nil
}

func (u *PVZUsecase) GetActiveReception(ctx context.Context, pvzID strfmt.UUID) (bool, *models.Reception, error) {
	return u.repo.GetActiveReception(ctx, pvzID)
}

func (u *PVZUsecase) AddProductToReception(ctx context.Context, pvzID strfmt.UUID, productType string) (*models.Product, error) {
	activeReception, activeReceptiont, err := u.repo.GetActiveReception(ctx, pvzID)
	if err != nil || !activeReception {
		return nil, errors.New("нет активной приемки для данного ПВЗ")
	}

	product := &models.Product{
		Type:        swag.String(productType),
		ReceptionID: &activeReceptiont.ID,
	}

	err = u.repo.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *PVZUsecase) DeleteLastProductFromReception(ctx context.Context, pvzID strfmt.UUID) error {
	activeReception, reception, err := u.repo.GetActiveReception(ctx, pvzID)
	if err != nil {
		return fmt.Errorf("failed to get active reception: %v", err)
	}
	if !activeReception || reception == nil {
		return errors.New("no active reception found")
	}

	if reception.Status == nil || *reception.Status != "in_progress" {
		return errors.New("can't delete products after reception is closed")
	}
	return u.repo.DeleteLastProduct(ctx, reception.ID)
}

func (u *PVZUsecase) CloseLastReception(ctx context.Context, pvzID strfmt.UUID) (*models.Reception, error) {
	active, reception, err := u.repo.GetActiveReception(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении активной приемки: %w", err)
	}
	if !active || reception == nil || reception.Status == nil || *reception.Status != "in_progress" {
		return nil, errors.New("приемка уже закрыта или отсутствует")
	}

	reception.Status = swag.String("closed")

	err = u.repo.UpdateReceptionStatus(ctx, *reception)
	if err != nil {
		return nil, fmt.Errorf("ошибка при закрытии приемки: %w", err)
	}
	_, reception, err = u.repo.GetCloseReception(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении закрытой приемки: %w", err)
	}
	return reception, nil
}

func (u *PVZUsecase) GetPVZs(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*operations.GetPvzOKBodyItems0, error) {
	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	pvzs, err := u.repo.GetPVZsWithReceptions(ctx, startDate, endDate, page, limit)
	if err != nil {
		return nil, fmt.Errorf("error retrieving pvzs: %w", err)
	}

	var result []*operations.GetPvzOKBodyItems0
	for _, pvzData := range pvzs {
		pvz := pvzData["pvz"].(*models.PVZ)
		receptions := pvzData["receptions"].([]map[string]interface{})

		var receptionItems []*operations.GetPvzOKBodyItems0ReceptionsItems0
		for _, receptionData := range receptions {
			reception := receptionData["reception"].(*models.Reception)
			products := receptionData["products"].([]models.Product)

			var productPtrs []*models.Product
			for i := range products {
				productPtrs = append(productPtrs, &products[i])
			}

			receptionItem := &operations.GetPvzOKBodyItems0ReceptionsItems0{
				Products:  productPtrs,
				Reception: reception,
			}
			receptionItems = append(receptionItems, receptionItem)
		}

		result = append(result, &operations.GetPvzOKBodyItems0{
			Pvz:        pvz,
			Receptions: receptionItems,
		})
	}

	return result, nil
}
