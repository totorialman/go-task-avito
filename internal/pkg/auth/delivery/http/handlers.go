package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"log/slog"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/totorialman/go-task-avito/internal/pkg/auth"
	"github.com/totorialman/go-task-avito/internal/pkg/utils/log"
	utils "github.com/totorialman/go-task-avito/internal/pkg/utils/sendError"
	"github.com/totorialman/go-task-avito/models"
	"github.com/totorialman/go-task-avito/restapi/operations"
)

type AuthHandler struct {
	authUsecase auth.AuthUsecase
	secret      string
}

func NewAuthHandler(authUsecase auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase, secret: os.Getenv("JWT_SECRET")}
}

func (h *AuthHandler) HandleDummyLogin(params operations.PostDummyLoginParams) middleware.Responder {
	logger := log.GetLoggerFromContext(params.HTTPRequest.Context()).With(slog.String("func", log.GetFuncName()))

	if params.Body.Role == nil {
		log.LogHandlerError(logger, errors.New("role is required"), http.StatusBadRequest)
		return operations.NewPostDummyLoginBadRequest().WithPayload(
			&models.Error{Message: swag.String("role is required")},
		)
	}

	token, err := h.authUsecase.GenerateDummyToken(params.HTTPRequest.Context(), *params.Body.Role)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("error generating dummy token: %w", err), http.StatusBadRequest)
		return operations.NewPostDummyLoginBadRequest().WithPayload(
			&models.Error{Message: swag.String(err.Error())},
		)
	}

	logger.Info("Dummy token generated successfully")
	return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
		http.SetCookie(w, &http.Cookie{
			Name:     "JWT",
			Value:    token,
			HttpOnly: true,
			Secure:   false,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})

		w.WriteHeader(http.StatusOK)
		_ = p.Produce(w, models.Token(token))
	})
}

func (h *AuthHandler) HandleLogin(params operations.PostLoginParams) middleware.Responder {
	logger := log.GetLoggerFromContext(params.HTTPRequest.Context()).With(slog.String("func", log.GetFuncName()))

	if params.Body.Email == nil || params.Body.Password == nil {
		log.LogHandlerError(logger, errors.New("email and password are required"), http.StatusUnauthorized)
		return operations.NewPostLoginUnauthorized().WithPayload(
			&models.Error{Message: swag.String("email and password are required")},
		)
	}

	email := string(*params.Body.Email)
	_, token, err := h.authUsecase.Login(params.HTTPRequest.Context(), email, *params.Body.Password)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("login failed: %w", err), http.StatusUnauthorized)
		return operations.NewPostLoginUnauthorized().WithPayload(
			&models.Error{Message: swag.String(err.Error())},
		)
	}

	logger.Info("User logged in successfully", slog.String("email", email))
	return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
		http.SetCookie(w, &http.Cookie{
			Name:     "JWT",
			Value:    token,
			HttpOnly: true,
			Secure:   false,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		})
		w.WriteHeader(http.StatusOK)
		_ = p.Produce(w, models.Token(token))
	})
}

func (h *AuthHandler) HandleSignUp(params operations.PostRegisterParams) middleware.Responder {
	logger := log.GetLoggerFromContext(params.HTTPRequest.Context()).With(slog.String("func", log.GetFuncName()))

	if params.Body.Email == nil || params.Body.Password == nil || params.Body.Role == nil {
		log.LogHandlerError(logger, errors.New("email, password and role are required"), http.StatusBadRequest)
		return operations.NewPostRegisterBadRequest().WithPayload(
			&models.Error{Message: swag.String("email, password and role are required")},
		)
	}

	email := string(*params.Body.Email)
	password := string(*params.Body.Password)
	role := string(*params.Body.Role)

	user, token, err := h.authUsecase.SignUp(params.HTTPRequest.Context(), email, password, role)
	if err != nil {
		log.LogHandlerError(logger, fmt.Errorf("signup failed: %w", err), http.StatusBadRequest)
		return operations.NewPostRegisterBadRequest().WithPayload(
			&models.Error{Message: swag.String(err.Error())},
		)
	}

	logger.Info("User signed up successfully", slog.String("email", email), slog.String("role", role))
	return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
		http.SetCookie(w, &http.Cookie{
			Name:     "JWT",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
		})
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.LogHandlerError(logger, fmt.Errorf("JSON encoding error: %w", err), http.StatusInternalServerError)
			utils.SendError(w, "Ошибка формирования JSON", http.StatusInternalServerError)
		}
	})
}
