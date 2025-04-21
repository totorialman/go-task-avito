package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/totorialman/go-task-avito/internal/pkg/auth"
	"github.com/totorialman/go-task-avito/internal/pkg/utils/log"
	"github.com/totorialman/go-task-avito/models"
	"golang.org/x/crypto/argon2"
)

func HashPassword(salt []byte, plainPassword string) string {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	return hex.EncodeToString(append(salt, hashedPass...))
}

func checkPassword(passHash string, plainPassword string) bool {
	passHashBytes, err := hex.DecodeString(passHash)
	if err != nil {
		return false
	}

	salt := make([]byte, 8)
	copy(salt, passHashBytes[:8])

	userPassHash := HashPassword(salt, plainPassword)

	return userPassHash == passHash
}

type AuthUsecase struct {
	authRepo auth.AuthRepo
}

func NewAuthUsecase(authRepo auth.AuthRepo) *AuthUsecase {
	return &AuthUsecase{
		authRepo: authRepo,
	}
}

func generateToken(role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", auth.ErrGeneratingToken
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func (u *AuthUsecase) GenerateDummyToken(ctx context.Context, role string) (string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	token, err := generateToken(role)
	if err != nil {
		log.LogHandlerError(logger, auth.ErrGeneratingToken, http.StatusInternalServerError)
		return "", auth.ErrGeneratingToken
	}

	return token, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	userID, role, passwordHash, err := uc.authRepo.GetUserCredsByEmail(ctx, email)
	if err != nil {
		log.LogHandlerError(logger, auth.ErrInvalidLogin, http.StatusUnauthorized)
		return nil, "", auth.ErrInvalidLogin
	}

	if !checkPassword(passwordHash, password) {
		log.LogHandlerError(logger, auth.ErrInvalidPassword, http.StatusUnauthorized)
		return nil, "", auth.ErrInvalidPassword
	}

	token, err := generateToken(role)
	if err != nil {
		log.LogHandlerError(logger, auth.ErrGeneratingToken, http.StatusInternalServerError)
		return nil, "", auth.ErrGeneratingToken
	}

	var emailFmt strfmt.Email
	emailFmt = strfmt.Email(email)

	user := &models.User{
		Email: &emailFmt,
		ID:    strfmt.UUID(userID.String()),
		Role:  &role,
	}

	return user, token, nil
}

var (
	ErrGeneratingSalt = errors.New("ошибка генерации соли")
)

func (uc *AuthUsecase) SignUp(ctx context.Context, email string, password string, role string) (*models.User, string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		log.LogHandlerError(logger, ErrGeneratingSalt, http.StatusInternalServerError)
		return nil, "", ErrGeneratingSalt
	}

	hashedPassword := HashPassword(salt, password)

	userID, err := uuid.NewV4()
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("ошибка генерации UUID: %w", err), http.StatusInternalServerError)
		return nil, "", err
	}

	emailFmt := strfmt.Email(email)
	newUser := &models.User{
		Email: &emailFmt,
		Role:  &role,
		ID:    strfmt.UUID(userID.String()),
	}

	err = uc.authRepo.InsertUser(ctx, newUser.ID, email, hashedPassword, role)
	if err != nil {
		log.LogHandlerError(logger, auth.ErrCreatingUser, http.StatusInternalServerError)
		return nil, "", auth.ErrCreatingUser
	}

	token, err := generateToken(role)
	if err != nil {
		log.LogHandlerError(logger, auth.ErrGeneratingToken, http.StatusInternalServerError)
		return nil, "", auth.ErrGeneratingToken
	}

	return newUser, token, nil
}
