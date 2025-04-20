package pg

import (
	"context"
	"log/slog"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype/pgxtype"
	"github.com/totorialman/go-task-avito/internal/pkg/utils/log"
)

const (
	getFieldProduct   = "SELECT id, name, price, image_url, weight FROM products WHERE id = ANY($1)"
	getRestaurantName = "SELECT name FROM restaurants WHERE id = $1"
	insertOrder       = `INSERT INTO orders (id, user_id, status, address_id, order_products,
		apartment_or_office, intercom, entrance, floor,
		courier_comment, leave_at_door, created_at, final_price) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
)

type AuthRepo struct {
	db pgxtype.Querier
}

func NewAuthRepo(db pgxtype.Querier) *AuthRepo {
	return &AuthRepo{db: db}
}

const getUserCredsQuery = `SELECT id, role, password_hash FROM users WHERE email = $1`

func (r *AuthRepo) GetUserCredsByEmail(ctx context.Context, email string) (uuid.UUID, string, string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	var id uuid.UUID
	var role, hash string
	err := r.db.QueryRow(ctx, getUserCredsQuery, email).Scan(&id, &role, &hash)
	if err != nil {
		logger.Error("GetUserCredsByEmail error: " + err.Error())
		return uuid.Nil, "", "", err
	}

	return id, role, hash, nil
}

func (repo *AuthRepo) InsertUser(ctx context.Context, userID strfmt.UUID, email string, hashedPassword string, role string) error {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	query := `INSERT INTO users (id, email, role, password_hash) VALUES ($1, $2, $3, $4)`
	_, err := repo.db.Exec(ctx, query, userID, email, role, hashedPassword)
	if err != nil {
		logger.Error("Failed to insert user into database: %w", err)
		return err
	}

	return nil
}
