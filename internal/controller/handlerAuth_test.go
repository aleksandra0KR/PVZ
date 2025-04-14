package controller

import (
	"bytes"
	"final/internal/domain"
	auth "final/internal/midleware"
	"final/internal/usecase"
	"final/internal/usecase/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_registerUser(t *testing.T) {
	type mockBehavior func(u *mock_usecase.MockUserUsecase, input *domain.InputUser)
	id := "123e4567-e89b-12d3-a456-426614174000"
	email := "user@example.com"
	role := "employee"
	password := "password"

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            *domain.InputUser
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Failure_Without_Password",
			inputBody: `{"email": "user@example.com", "role": "employee"}`,
			inputUser: &domain.InputUser{Email: &email, Role: &role},
			mockBehavior: func(u *mock_usecase.MockUserUsecase, input *domain.InputUser) {
				u.EXPECT().Register(input).Return(nil, domain.ErrInvalidCredentials)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: `{
			    "message": "invalid email or password" }`,
		},
		{
			name:      "Failure_Without_Email",
			inputBody: `{"password": "password", "role": "employee"}`,
			inputUser: &domain.InputUser{Password: &password, Role: &role},
			mockBehavior: func(u *mock_usecase.MockUserUsecase, input *domain.InputUser) {
				u.EXPECT().Register(input).Return(nil, domain.ErrInvalidCredentials)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: `{
			    "message": "invalid email or password" }`,
		},
		{
			name:      "Failure_Without_Role",
			inputBody: `{"email": "user@example.com", "password": "password"}`,
			inputUser: &domain.InputUser{Email: &email, Password: &password},
			mockBehavior: func(u *mock_usecase.MockUserUsecase, input *domain.InputUser) {
				u.EXPECT().Register(input).Times(0)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: `{
			    "message": "invalid role" }`,
		},
		{
			name:      "Success",
			inputBody: `{"email": "user@example.com", "password": "password", "role": "employee"}`,
			inputUser: &domain.InputUser{Email: &email, Role: &role, Password: &password},
			mockBehavior: func(u *mock_usecase.MockUserUsecase, input *domain.InputUser) {
				u.EXPECT().Register(input).Return(&domain.User{
					ID:    &id,
					Email: &email,
					Role:  &role,
				}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: `{
				"id": "123e4567-e89b-12d3-a456-426614174000",
				"email": "user@example.com",
				"role": "employee"
			}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userUsecase := mock_usecase.NewMockUserUsecase(c)
			test.mockBehavior(userUsecase, test.inputUser)

			services := usecase.Usecase{UserUsecase: userUsecase}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/register", handler.Register)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_dummyLogin(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Error("error loading .env file")
	}
	err := auth.InitAuthFromConfig()
	if err != nil {
		log.Error(err)
		return
	}
	tokenModerator, err := auth.GenerateJWT("0", "moderator")
	if err != nil {
		log.Error(err)
		return
	}
	tokenEmployee, err := auth.GenerateJWT("0", "employee")
	if err != nil {
		log.Error(err)
		return
	}

	tests := []struct {
		name                 string
		inputBody            string
		expectedResponseBody string
		expectedStatusCode   int
	}{
		{
			name:                 "Moderator role",
			inputBody:            `{"role": "moderator"}`,
			expectedResponseBody: `"` + tokenModerator + `"`,
			expectedStatusCode:   http.StatusOK,
		},
		{
			name:                 "Employee role",
			inputBody:            `{"role": "employee"}`,
			expectedResponseBody: `"` + tokenEmployee + `"`,
			expectedStatusCode:   http.StatusOK,
		},
		{
			name:                 "Invalid role",
			inputBody:            `{"role": "hacker"}`,
			expectedResponseBody: `{"message":"invalid role"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name:                 "Failure_Wrong_Json_Format",
			inputBody:            `{"role": "hacker}`,
			expectedResponseBody: `{"message":"invalid input data"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userUsecase := mock_usecase.NewMockUserUsecase(c)

			services := usecase.Usecase{UserUsecase: userUsecase}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/dummyLogin", handler.DummyLogin)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/dummyLogin", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestRegister(t *testing.T) {
	type mockBehavior func(p *mock_usecase.MockUserUsecase, input *domain.InputUser)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	email := "user@example.com"
	password := "string"
	role := "employee"
	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	userInput := domain.InputUser{
		Email:    &email,
		Password: &password,
		Role:     &role,
	}

	user := domain.User{
		Email: &email,
		ID:    &id,
		Role:  &role,
	}

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            *domain.InputUser
		mockBehavior         mockBehavior
		expectedResponseBody string
		expectedStatusCode   int
	}{
		{
			name: "Successful register",
			inputBody: `{
				  "email": "user@example.com",
				  "password": "string",
				  "role": "employee" }`,
			inputUser: &userInput,
			mockBehavior: func(p *mock_usecase.MockUserUsecase, input *domain.InputUser) {
				p.EXPECT().Register(input).Return(&user, nil)
			},
			expectedResponseBody: `{
				  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
				  "email": "user@example.com",
				  "role": "employee" }`,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "Failure without password",
			inputBody: `{
				  "email": "user@example.com",
				  "role": "employee" }`,
			inputUser: &domain.InputUser{
				Email: &email,
				Role:  &role},
			mockBehavior: func(p *mock_usecase.MockUserUsecase, input *domain.InputUser) {
				p.EXPECT().Register(input).Return(nil, domain.ErrInvalidCredentials)
			},
			expectedResponseBody: `{"message":"invalid email or password"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name: "Failure without email",
			inputBody: `{
				  "password": "string",
				  "role": "employee" }`,
			inputUser: &domain.InputUser{
				Password: &password,
				Role:     &role},
			mockBehavior: func(p *mock_usecase.MockUserUsecase, input *domain.InputUser) {
				p.EXPECT().Register(input).Return(nil, domain.ErrInvalidCredentials)
			},
			expectedResponseBody: `{"message":"invalid email or password"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name: "Failure without role",
			inputBody: `{
				  "email": "user@example.com",
				   "password": "string"}`,
			inputUser: &domain.InputUser{
				Email:    &email,
				Password: &password},
			mockBehavior: func(p *mock_usecase.MockUserUsecase, input *domain.InputUser) {
				p.EXPECT().Register(input).Times(0)
			},
			expectedResponseBody: `{"message":"invalid role"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name: "Failure_Wrong_Json_Format",
			inputBody: `{
				  "email": user@example.com",
				   "password": "string"}`,
			inputUser: nil,
			mockBehavior: func(p *mock_usecase.MockUserUsecase, input *domain.InputUser) {
				p.EXPECT().Register(input).Times(0)
			},
			expectedResponseBody: `{ "message": "invalid input data" }`,
			expectedStatusCode:   http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockUserUsecase(c)
			test.mockBehavior(repo, test.inputUser)

			services := usecase.Usecase{UserUsecase: repo}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/register", handler.Register)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestLogin(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Error("error loading .env file")
	}
	err := auth.InitAuthFromConfig()
	if err != nil {
		log.Error(err)
		return
	}
	_, err = auth.GenerateJWT("0", "moderator")
	if err != nil {
		log.Error(err)
		return
	}
	tokenEmployee, err := auth.GenerateJWT("0", "employee")
	if err != nil {
		log.Error(err)
		return
	}
	id := "0"
	email := "user@example.com"
	password := "string"
	role := "employee"
	userInput := domain.InputUser{
		Email:    &email,
		Password: &password,
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*userInput.Password), bcrypt.DefaultCost)
	hashStr := string(hash)
	if err != nil {
		log.Error(err)
		return
	}
	user := domain.User{
		Email:    &email,
		Password: &hashStr,
		ID:       &id,
		Role:     &role,
	}
	type mockBehavior func(p *mock_usecase.MockUserUsecase, input *domain.InputUser)
	tests := []struct {
		name                 string
		inputBody            string
		inputUser            *domain.InputUser
		outputReception      *domain.User
		mockBehavior         mockBehavior
		expectedResponseBody string
		expectedStatusCode   int
	}{
		{
			name: "Successful login",
			inputBody: `{
				  "email": "user@example.com",
				  "password": "string"}`,
			inputUser: &userInput,
			mockBehavior: func(p *mock_usecase.MockUserUsecase, input *domain.InputUser) {
				p.EXPECT().Login(input).Return(&user, nil)
			},
			expectedResponseBody: `"` + tokenEmployee + `"`,
			expectedStatusCode:   http.StatusOK,
		},
		{
			name: "Failure_Wrong_Json_Format",
			inputBody: `{
				  "email": "user@example.com",
				  "password": string"}`,
			inputUser: nil,
			mockBehavior: func(p *mock_usecase.MockUserUsecase, input *domain.InputUser) {
				p.EXPECT().Login(input).Times(0)
			},
			expectedResponseBody: `{"message":"invalid input data"}`,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name: "Failure_Invalid_Credentials",
			inputBody: `{
				  "email": "user@example.com",
				  "password": "wrongpassword"}`,
			inputUser: &domain.InputUser{
				Email:    &email,
				Password: func() *string { s := "wrongpassword"; return &s }(),
			},
			mockBehavior: func(p *mock_usecase.MockUserUsecase, input *domain.InputUser) {
				p.EXPECT().Login(input).Return(nil, domain.ErrInvalidCredentials)
			},
			expectedResponseBody: `{"message":"invalid email or password"}`,
			expectedStatusCode:   http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockUserUsecase(c)
			test.mockBehavior(repo, test.inputUser)

			services := usecase.Usecase{UserUsecase: repo}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/login", handler.Login)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
