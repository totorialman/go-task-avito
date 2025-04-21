package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/totorialman/go-task-avito/models"
	"github.com/totorialman/go-task-avito/restapi/operations"
)

type DummyAuthUsecase struct {
	TokenResult struct {
		Called bool
		Role   string
		Token  string
		Err    error
	}
	LoginResult struct {
		Called bool
		Email  string
		Pass   string
		User   *models.User
		Token  string
		Err    error
	}
	SignUpResult struct {
		Called bool
		Email  string
		Pass   string
		Role   string
		User   *models.User
		Token  string
		Err    error
	}
}

func (m *DummyAuthUsecase) GenerateDummyToken(ctx context.Context, role string) (string, error) {
	m.TokenResult.Called = true
	m.TokenResult.Role = role
	return m.TokenResult.Token, m.TokenResult.Err
}

func (m *DummyAuthUsecase) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	m.LoginResult.Called = true
	m.LoginResult.Email = email
	m.LoginResult.Pass = password
	return m.LoginResult.User, m.LoginResult.Token, m.LoginResult.Err
}

func (m *DummyAuthUsecase) SignUp(ctx context.Context, email, password, role string) (*models.User, string, error) {
	m.SignUpResult.Called = true
	m.SignUpResult.Email = email
	m.SignUpResult.Pass = password
	m.SignUpResult.Role = role
	return m.SignUpResult.User, m.SignUpResult.Token, m.SignUpResult.Err
}

func TestAuthHandler_HandleDummyLogin(t *testing.T) {
	tests := []struct {
		name           string
		role           string
		mockToken      string
		mockError      error
		expectedStatus int
	}{
		{"Success", "admin", "test-token", nil, http.StatusOK},
		{"Missing role", "", "", errors.New("role is required"), http.StatusBadRequest},
		{"Token generation error", "user", "", errors.New("generation error"), http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &DummyAuthUsecase{}
			mock.TokenResult.Token = tt.mockToken
			mock.TokenResult.Err = tt.mockError

			handler := NewAuthHandler(mock)

			var jsonBody []byte
			if tt.role != "" {
				jsonBody, _ = json.Marshal(map[string]string{"role": tt.role})
			}

			req := httptest.NewRequest("POST", "/dummy-login", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			var body *operations.PostDummyLoginBody
			if tt.role != "" {
				body = &operations.PostDummyLoginBody{Role: swag.String(tt.role)}
			}

			// Проверка на nil
			var params operations.PostDummyLoginParams
			if body != nil {
				params = operations.PostDummyLoginParams{
					HTTPRequest: req,
					Body:        *body, // Теперь безопасно разыменовываем, потому что body не nil
				}
			} else {
				params = operations.PostDummyLoginParams{
					HTTPRequest: req,
					Body:        operations.PostDummyLoginBody{}, // Если role пустая, создаём пустую структуру
				}
			}

			resp := handler.HandleDummyLogin(params)
			rr := httptest.NewRecorder()
			resp.WriteResponse(rr, runtime.JSONProducer())

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestAuthHandler_HandleLogin(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		password       string
		mockUser       *models.User
		mockToken      string
		mockError      error
		expectedStatus int
	}{
		{
			name:     "Success",
			email:    "test@example.com",
			password: "password",
			mockUser: func() *models.User {
				email := strfmt.Email("test@example.com")
				return &models.User{Email: &email}
			}(),
			mockToken:      "test-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing credentials",
			email:          "",
			password:       "",
			mockUser:       nil,
			mockToken:      "",
			mockError:      errors.New("missing credentials"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Login error",
			email:          "test@example.com",
			password:       "wrong",
			mockUser:       nil,
			mockToken:      "",
			mockError:      errors.New("invalid credentials"),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &DummyAuthUsecase{}
			mock.LoginResult.User = tt.mockUser
			mock.LoginResult.Token = tt.mockToken
			mock.LoginResult.Err = tt.mockError

			handler := NewAuthHandler(mock)

			email := strfmt.Email(tt.email)
			password := swag.String(tt.password)

			reqBody := operations.PostLoginBody{
				Email:    &email,
				Password: password,
			}

			if tt.email == "" || tt.password == "" {
				reqBody.Email = nil
				reqBody.Password = nil
			}

			jsonBody, _ := json.Marshal(reqBody)

			req := httptest.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			params := operations.PostLoginParams{
				HTTPRequest: req,
				Body:        reqBody,
			}

			resp := handler.HandleLogin(params)
			rr := httptest.NewRecorder()
			resp.WriteResponse(rr, runtime.JSONProducer())

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if rr.Code == http.StatusOK {
				found := false
				for _, c := range rr.Result().Cookies() {
					if c.Name == "JWT" && c.Value == tt.mockToken {
						found = true
					}
				}
				if !found {
					t.Error("JWT cookie not found or value mismatch")
				}
			}
		})
	}
}


func TestAuthHandler_HandleSignUp(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		password       string
		role           string
		mockUser       *models.User
		mockToken      string
		mockError      error
		expectedStatus int
	}{
		{
			"Success",
			"new@example.com",
			"password",
			"user",
			&models.User{
				Email: func() *strfmt.Email {
					e := strfmt.Email("new@example.com")
					return &e
				}(),
				Role: swag.String("user"),
			},
			"test-token",
			nil,
			http.StatusOK,
		},
		{
			"Missing fields", 
			"", "", "", 
			nil, "", nil, 
			http.StatusBadRequest,
		},
		{
			"Signup error",
			"exists@example.com",
			"password",
			"user",
			nil, "", errors.New("user exists"), http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &DummyAuthUsecase{}
			mock.SignUpResult.User = tt.mockUser
			mock.SignUpResult.Token = tt.mockToken
			mock.SignUpResult.Err = tt.mockError

			handler := NewAuthHandler(mock)

			reqBody := operations.PostRegisterBody{
				Email:    func() *strfmt.Email { e := strfmt.Email(tt.email); return &e }(),
				Password: swag.String(tt.password),
				Role:     swag.String(tt.role),
			}
			
			// Если поля пустые, указываем nil
			if tt.email == "" || tt.password == "" || tt.role == "" {
				reqBody.Email = nil
				reqBody.Password = nil
				reqBody.Role = nil
			}

			jsonBody, _ := json.Marshal(reqBody)
			
			req := httptest.NewRequest("POST", "/signup", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			params := operations.PostRegisterParams{
				HTTPRequest: req,
				Body:        reqBody,
			}

			resp := handler.HandleSignUp(params)
			rr := httptest.NewRecorder()
			resp.WriteResponse(rr, runtime.JSONProducer())

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if rr.Code == http.StatusOK {
				found := false
				for _, c := range rr.Result().Cookies() {
					if c.Name == "JWT" && c.Value == tt.mockToken {
						found = true
					}
				}
				if !found {
					t.Error("JWT cookie not found or mismatch")
				}

				var user models.User
				err := json.NewDecoder(rr.Body).Decode(&user)
				if err != nil {
					t.Errorf("failed to decode response body: %v", err)
				}
				if user.Email == nil || *user.Email != *tt.mockUser.Email {
					t.Errorf("expected email %s, got %v", *tt.mockUser.Email, user.Email)
				}
			}
		})
	}
}


func TestNewAuthHandler(t *testing.T) {
	mock := &DummyAuthUsecase{}
	handler := NewAuthHandler(mock)
	if handler == nil {
		t.Fatal("handler is nil")
	}
	if handler.authUsecase != mock {
		t.Fatal("handler usecase mismatch")
	}
}
