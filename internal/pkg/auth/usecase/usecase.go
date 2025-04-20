package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"os"
	"strings"
	"time"
	"unicode"

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

const (
	minNameLength  = 2
	maxNameLength  = 25
	minPhoneLength = 7
	maxPhoneLength = 15
	maxLoginLength = 20
	minLoginLength = 3
	minPassLength  = 8
	maxPassLength  = 25
)

const allowedRunes = "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

func isValidName(name string) bool {
	if len(name) < minNameLength || len(name) > maxNameLength {
		return false
	}
	for _, r := range name {
		if !strings.ContainsRune(allowedRunes, r) {
			return false
		}
	}
	return true
}

func isValidPhone(phone string) bool {
	if len(phone) < minPhoneLength || len(phone) > maxPhoneLength {
		return false
	}
	for _, r := range phone {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func validLogin(login string) bool {
	if len(login) < minLoginLength || len(login) > maxLoginLength {
		return false
	}
	for _, char := range login {
		if !strings.Contains(allowedChars, string(char)) {
			return false
		}
	}
	return true
}

func validPassword(password string) bool {
	var up, low, digit, special bool

	if len(password) < minPassLength || len(password) > maxPassLength {
		return false
	}

	for _, char := range password {

		switch {
		case unicode.IsUpper(char):
			up = true
		case unicode.IsLower(char):
			low = true
		case unicode.IsDigit(char):
			digit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			special = true
		default:
			return false
		}
	}

	return up && low && digit && special
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
		logger.Error(auth.ErrGeneratingToken.Error())
		return "", auth.ErrGeneratingToken
	}

	return token, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (*models.User, string, string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))
	userID, role, passwordHash, err := uc.authRepo.GetUserCredsByEmail(ctx, email)
	if err != nil {
		return nil, "", "", auth.ErrInvalidLogin
	}
	logger.Error("%v, checkPassword([]byte(passwordHash), password): %s", checkPassword(passwordHash, password))
	if !checkPassword(passwordHash, password) {
		return nil, "", "", auth.ErrInvalidPassword
	}

	token, err := generateToken(role)
	if err != nil {
		return nil, "", "", auth.ErrGeneratingToken
	}

	csrf := uuid.Must(uuid.NewV4()).String()

	var emailFmt strfmt.Email
	emailFmt = strfmt.Email(email)

	user := &models.User{
		Email: &emailFmt,
		ID:    strfmt.UUID(userID.String()),
		Role:  &role,
	}

	return user, token, csrf, nil
}

var (
	ErrGeneratingSalt = errors.New("ошибка генерации соли")
)

func (uc *AuthUsecase) SignUp(ctx context.Context, email string, password string, role string) (*models.User, string, string, error) {
	logger := log.GetLoggerFromContext(ctx).With(slog.String("func", log.GetFuncName()))

	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		logger.Error("Error generating salt: " + err.Error())
		return nil, "", "", ErrGeneratingSalt
	}

	hashedPassword := HashPassword(salt, password)

	userID, err := uuid.NewV4()
	if err != nil {
		logger.Error("Error generating UUID: " + err.Error())
		return nil, "", "", err
	}
	emailFmt := strfmt.Email(email)
	newUser := &models.User{
		Email: &emailFmt,
		Role:  &role,
		ID:    strfmt.UUID(userID.String()),
	}

	err = uc.authRepo.InsertUser(ctx, newUser.ID, email, hashedPassword, role)
	if err != nil {
		logger.Error("Error creating user: " + err.Error())
		return nil, "", "", auth.ErrCreatingUser
	}

	token, err := generateToken(role)
	if err != nil {
		logger.Error("Error generating token: " + err.Error())
		return nil, "", "", auth.ErrGeneratingToken
	}

	csrfToken := uuid.Must(uuid.NewV4()).String()

	logger.Info("User successfully registered")
	return newUser, token, csrfToken, nil
}
