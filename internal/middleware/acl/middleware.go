package acl

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/casbin/casbin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/golang-jwt/jwt"
)

func NewAclMiddleware(next http.Handler) http.Handler {

	e, err := casbin.NewEnforcerSafe("internal/middleware/acl/model.conf", "internal/middleware/acl/policy.csv")
	if err != nil {
		log.Fatalf("failed to load enforcer: %v", err)
	}

	skipPaths := map[string]bool{
		"/dummyLogin": true,
		"/login":      true,
		"/register":   true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/pvz" {
			log.Printf("Skipping %v %v %v", r.Method, r.URL.Path, r)
			next.ServeHTTP(w, r)
			return
		}
		// Пропускаем указанные пути
		if skipPaths[r.URL.Path] {
			log.Printf("%v %v %v", next, w, r)
			next.ServeHTTP(w, r)
			return
		}

		var token string
		// Получаем токен из куки или заголовка Authorization
		cookie, err := r.Cookie("JWT")
		if err == nil {
			token = cookie.Value
		} else {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// Если токен не найден, отправляем ошибку Unauthorized
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Получаем роль пользователя из токена
		role, _ := getRoleFromToken(cookie.Value)

		// Проверяем права доступа через Casbin
		res, _ := e.EnforceSafe(role, r.URL.Path, r.Method)
		log.Printf("path=%s method=%s role=%s access=%v", r.URL.Path, r.Method, role, res)

		if res {
			// Если доступ разрешен, передаем запрос дальше
			log.Printf("Passing to next handler: %s", r.URL.Path)
			next.ServeHTTP(w, r)
		} else {
			// Если доступ запрещен, отправляем кастомное сообщение в ответе
			w.WriteHeader(http.StatusForbidden)
			jsonResponse := struct {
				Message string `json:"message"`
			}{
				Message: "Доступ запрещен",
			}
			responseBody, _ := json.Marshal(jsonResponse)
			w.Header().Set("Content-Type", "application/json")
			w.Write(responseBody)
			return
		}
	})
}

func getRoleFromToken(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("missing JWT_SECRET")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role, ok := claims["role"].(string)
		if !ok {
			return "", errors.New("role not found or invalid in token")
		}
		return role, nil
	}

	return "", errors.New("invalid token")
}

func GetEntMw() middleware.Builder {
	e, _ := casbin.NewEnforcerSafe("./model.conf", "./policy.csv")

	skipPaths := map[string]bool{
		"/dummyLogin": true,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if skipPaths[r.URL.Path] {
				log.Printf("%v %v %v", next, w, r)
				next.ServeHTTP(w, r)
				return
			}

			role := ""
			res, _ := e.EnforceSafe(role, r.URL.Path, r.Method)
			log.Printf("path=%s role=%s access=%v", r.URL.Path, role, res)
			if res {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		})
	}
}
