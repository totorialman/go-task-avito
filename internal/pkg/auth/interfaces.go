package auth

import (
	"context"
	"errors"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/totorialman/go-task-avito/models"
)

var (
	ErrCreatingUser       = errors.New("Ошибка в создании пользователя")
	ErrUserNotFound       = errors.New("Пользователь не найден")
	ErrInvalidLogin       = errors.New("Неверный формат логина")
	ErrInvalidPassword    = errors.New("Неверный формат пароля")
	ErrInvalidCredentials = errors.New("Неверный логин или пароль")
	ErrGeneratingToken    = errors.New("Ошибка генерации токена")
	ErrInvalidName        = errors.New("Имя и фамилия должны содержать только русские буквы и быть от 2 до 25 символов")
	ErrInvalidPhone       = errors.New("Некорректный номер телефона")
	ErrUUID               = errors.New("Ошибка создания UUID")
	ErrSamePassword       = errors.New("Новый пароль совпадает со старым")
	ErrBasePath           = errors.New("Базовый путь для картинок не установлен")
	ErrFileCreation       = errors.New("Ошибка при создании файла")
	ErrFileSaving         = errors.New("Ошибка при сохранении файла")
	ErrFileDeletion       = errors.New("Ошибка при удалении файла")
	ErrDBError            = errors.New("Ошибка БД")
	ErrAddressNotFound    = errors.New("Ошибка поиска адреса")
)

type AuthRepo interface {
	InsertUser(ctx context.Context, userID strfmt.UUID, email string, hashedPassword string, role string) error

	GetUserCredsByEmail(ctx context.Context, email string) (userID uuid.UUID, role string, passwordHash string, err error)
}

type AuthUsecase interface {
	GenerateDummyToken(ctx context.Context, role string) (string, error)
	SignUp(ctx context.Context, email string, password string, role string) (*models.User, string, string, error)

	Login(ctx context.Context, email, password string) (*models.User, string, string, error)
}
