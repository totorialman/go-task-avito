package http

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/totorialman/go-task-avito/internal/pkg/auth"
	"github.com/totorialman/go-task-avito/internal/pkg/metrics"
	utils "github.com/totorialman/go-task-avito/internal/pkg/utils/sendError"
	"github.com/totorialman/go-task-avito/models"
	"github.com/totorialman/go-task-avito/restapi/operations"
)

type AuthHandler struct {
	authUsecase auth.AuthUsecase
	secret      string
	mt          *metrics.ProductMetrics
}

func NewAuthHandler(authUsecase auth.AuthUsecase, mt *metrics.ProductMetrics) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase, secret: os.Getenv("JWT_SECRET"),mt: mt}
}

func (h *AuthHandler) HandleDummyLogin(params operations.PostDummyLoginParams) middleware.Responder {
	h.mt.IncreaseHits()
	if params.Body.Role == nil {
		
		return operations.NewPostDummyLoginBadRequest().WithPayload(
			&models.Error{Message: swag.String("role is required")},
		)
	}

	token, err := h.authUsecase.GenerateDummyToken(params.HTTPRequest.Context(), *params.Body.Role)
	if err != nil {
		return operations.NewPostDummyLoginBadRequest().WithPayload(
			&models.Error{Message: swag.String(err.Error())},
		)
	}
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
		_, err := w.Write([]byte("Успешная авторизация"))
		if err != nil {
			http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
		}
		//_ = p.Produce(w, models.Token(token))
	})
}

func (h *AuthHandler) HandleLogin(params operations.PostLoginParams) middleware.Responder {

	if params.Body.Email == nil || params.Body.Password == nil {
		return operations.NewPostLoginUnauthorized().WithPayload(
			&models.Error{Message: swag.String("email and password are required")},
		)
	}

	email := string(*params.Body.Email)
	_, token, csrfToken, err := h.authUsecase.Login(params.HTTPRequest.Context(), email, *params.Body.Password)
	if err != nil {
		return operations.NewPostLoginUnauthorized().WithPayload(
			&models.Error{Message: swag.String(err.Error())},
		)
	}

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

		http.SetCookie(w, &http.Cookie{
			Name:     "CSRF-Token",
			Value:    csrfToken,
			HttpOnly: false,
			Secure:   false,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		})
		w.Header().Set("X-CSRF-Token", csrfToken)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Успешная авторизация"))
		if err != nil {
			http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
		}
	})
}

func (h *AuthHandler) HandleSignUp(params operations.PostRegisterParams) middleware.Responder {
	if params.Body.Email == nil || params.Body.Password == nil || params.Body.Role == nil {
		return operations.NewPostRegisterBadRequest().WithPayload(
			&models.Error{Message: swag.String("email, password and role are required")},
		)
	}

	email := string(*params.Body.Email)
	password := string(*params.Body.Password)
	role := string(*params.Body.Role)

	user, token, csrfToken, err := h.authUsecase.SignUp(params.HTTPRequest.Context(), email, password, role)
	if err != nil {
		return operations.NewPostRegisterBadRequest().WithPayload(
			&models.Error{Message: swag.String(err.Error())},
		)
	}

	return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
		http.SetCookie(w, &http.Cookie{
			Name:     "JWT",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "CSRF-Token",
			Value:    csrfToken,
			HttpOnly: false,
			Secure:   true,
			Expires:  time.Now().Add(24 * time.Hour),
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		w.Header().Set("X-CSRF-Token", csrfToken)
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(user); err != nil {
			utils.SendError(w, "Ошибка формирования JSON", http.StatusInternalServerError)
		}
	})
}
