package repository

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"log/slog"

	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/totorialman/go-task-avito/internal/pkg/utils/log"
	"github.com/totorialman/go-task-avito/models"
)

type PVZRepo struct {
	db *pgxpool.Pool
}

func NewPVZRepo(db *pgxpool.Pool) *PVZRepo {
	return &PVZRepo{db: db}
}

func (r *PVZRepo) CreatePVZ(ctx context.Context, pvz *models.PVZ) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	query := `INSERT INTO pvz (id, city, registration_date) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, pvz.ID, pvz.City, pvz.RegistrationDate)

	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("db exec error: %w", err), http.StatusInternalServerError)
		return fmt.Errorf("db exec error: %w", err)
	}

	return nil
}

func (r *PVZRepo) CreateReception(ctx context.Context, reception *models.Reception) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	currentTime := strfmt.DateTime(time.Now())
	reception.DateTime = &currentTime

	_, err := r.db.Exec(ctx, `
        INSERT INTO receptions (id, pvz_id, status, date_time)
        VALUES ($1, $2, $3, $4)`,
		reception.ID, reception.PvzID, reception.Status, reception.DateTime)

	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("failed to create reception: %w", err), http.StatusInternalServerError)
		return fmt.Errorf("failed to create reception: %w", err)
	}

	return nil
}

func (r *PVZRepo) GetActiveReception(ctx context.Context, pvzID strfmt.UUID) (bool, *models.Reception, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	var reception models.Reception
	err := r.db.QueryRow(ctx, `
        SELECT id, pvz_id, status, date_time
        FROM receptions
        WHERE pvz_id = $1 AND status = 'in_progress'
        LIMIT 1`, pvzID).Scan(&reception.ID, &reception.PvzID, &reception.Status, &reception.DateTime)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("failed to get active reception: %w", err), http.StatusInternalServerError)
		return false, nil, nil
	}

	return true, &reception, nil
}

func (r *PVZRepo) CreateProduct(ctx context.Context, product *models.Product) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	currentTime := strfmt.DateTime(time.Now())
	product.DateTime = currentTime
	_, err := r.db.Exec(ctx, `INSERT INTO products (reception_id, type, date_time) VALUES ($1, $2, $3)`,
		product.ReceptionID, product.Type, product.DateTime)

	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("failed to create product: %w", err), http.StatusInternalServerError)
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *PVZRepo) DeleteLastProduct(ctx context.Context, receptionID strfmt.UUID) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	var productID strfmt.UUID
	err := r.db.QueryRow(ctx, `
        SELECT id
        FROM products
        WHERE reception_id = $1
        ORDER BY date_time DESC
        LIMIT 1`, receptionID).Scan(&productID)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("failed to get the last product: %w", err), http.StatusInternalServerError)
		return fmt.Errorf("failed to get the last product: %w", err)
	}

	_, err = r.db.Exec(ctx, `DELETE FROM products WHERE id = $1`, productID)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("failed to delete the last product: %w", err), http.StatusInternalServerError)
		return fmt.Errorf("failed to delete the last product: %w", err)
	}

	return nil
}

func (r *PVZRepo) UpdateReceptionStatus(ctx context.Context, reception models.Reception) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	currentTime := strfmt.DateTime(time.Now())
	reception.DateTime = &currentTime
	query := `UPDATE receptions SET status = $1, date_time = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, reception.Status, reception.DateTime, reception.ID)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("failed to update reception status: %w", err), http.StatusInternalServerError)
		return fmt.Errorf("failed to update reception status: %w", err)
	}
	return nil
}

func (r *PVZRepo) GetCloseReception(ctx context.Context, pvzID strfmt.UUID) (bool, *models.Reception, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	var reception models.Reception
	err := r.db.QueryRow(ctx, `
        SELECT id, pvz_id, status, date_time
        FROM receptions
        WHERE pvz_id = $1 AND status = 'closed'
        LIMIT 1`, pvzID).Scan(&reception.ID, &reception.PvzID, &reception.Status, &reception.DateTime)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("failed to get closed reception: %w", err), http.StatusInternalServerError)
		return false, nil, nil
	}
	return true, &reception, nil
}

func (r *PVZRepo) GetPVZsWithReceptions(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]map[string]interface{}, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	offset := (page - 1) * limit

	query := `
		SELECT pvz.id, pvz.city, pvz.registration_date, 
		       r.id AS reception_id, r.date_time AS reception_date, r.status AS reception_status,
		       p.id AS product_id, p.type AS product_type, p.date_time AS product_date_time
		FROM pvz
		LEFT JOIN receptions r ON pvz.id = r.pvz_id
		LEFT JOIN products p ON r.id = p.reception_id
		WHERE ($1::timestamp IS NULL OR r.date_time >= $1)
		  AND ($2::timestamp IS NULL OR r.date_time <= $2)
		ORDER BY r.date_time
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(ctx, query, startDate, endDate, limit, offset)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("error executing query: %w", err), http.StatusInternalServerError)
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	pvzsMap := make(map[strfmt.UUID]map[string]interface{})

	for rows.Next() {
		var pvzID strfmt.UUID
		var city string
		var registrationDate time.Time
		var receptionID strfmt.UUID
		var receptionDate time.Time
		var receptionStatus string
		var productID strfmt.UUID
		var productType string
		var productDateTime time.Time

		if err := rows.Scan(&pvzID, &city, &registrationDate, &receptionID, &receptionDate, &receptionStatus, &productID, &productType, &productDateTime); err != nil {
			log.LogHandlerError(logger, fmt.Errorf("error scanning row: %w", err), http.StatusInternalServerError)
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		if _, exists := pvzsMap[pvzID]; !exists {
			pvzsMap[pvzID] = map[string]interface{}{
				"pvz": &models.PVZ{
					ID:               pvzID,
					City:             &city,
					RegistrationDate: strfmt.DateTime(registrationDate),
				},
				"receptions": make(map[strfmt.UUID]map[string]interface{}),
			}
		}
		pvzData := pvzsMap[pvzID]
		receptions := pvzData["receptions"].(map[strfmt.UUID]map[string]interface{})

		if _, exists := receptions[receptionID]; !exists {
			receptions[receptionID] = map[string]interface{}{
				"reception": &models.Reception{
					ID:       receptionID,
					PvzID:    &pvzID,
					DateTime: (*strfmt.DateTime)(&receptionDate),
					Status:   &receptionStatus,
				},
				"products": []models.Product{},
			}
		}

		if productID != "" {
			product := models.Product{
				ID:          productID,
				ReceptionID: &receptionID,
				DateTime:    strfmt.DateTime(productDateTime),
				Type:        &productType,
			}
			receptions[receptionID]["products"] = append(receptions[receptionID]["products"].([]models.Product), product)
		}
	}

	if err := rows.Err(); err != nil {
		log.LogHandlerError(logger, fmt.Errorf("row iteration error: %w", err), http.StatusInternalServerError)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	var result []map[string]interface{}
	for _, pvzData := range pvzsMap {
		receptionMap := pvzData["receptions"].(map[strfmt.UUID]map[string]interface{})
		var receptionList []map[string]interface{}
		for _, v := range receptionMap {
			receptionList = append(receptionList, v)
		}
		pvzData["receptions"] = receptionList
		result = append(result, pvzData)
	}

	return result, nil
}
