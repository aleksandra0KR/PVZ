package controller_test

import (
	"bytes"
	"final/internal/controller"
	"final/internal/domain"
	"final/internal/usecase"
	"final/internal/usecase/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_createPVZ(t *testing.T) {
	type mockBehavior func(p *mock_usecase.MockPVZUsecase, input *domain.PVZ)
	allowedCity := "Москва"
	wrongCity := "Париж"
	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	layout := "2006-01-02T15:04:05.000Z"
	str := "2025-04-12T20:10:25.102Z"
	dateTime, err := time.Parse(layout, str)

	if err != nil {
		log.Error(err)
	}
	tests := []struct {
		name                 string
		inputBody            string
		inputPVZ             *domain.PVZ
		outputPVZ            *domain.PVZ
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success_With_All_Parameters",
			inputBody: `{  
				"id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  				"registrationDate": "2025-04-12T20:10:25.102Z",
				"city": "Москва" }`,
			inputPVZ:  &domain.PVZ{City: &allowedCity, RegistrationDate: &dateTime, ID: &id},
			outputPVZ: &domain.PVZ{City: &allowedCity, ID: &id, RegistrationDate: &dateTime},
			mockBehavior: func(p *mock_usecase.MockPVZUsecase, inputPVZ *domain.PVZ) {
				p.EXPECT().CreatePVZ(inputPVZ).Return(&domain.PVZ{City: &allowedCity, ID: &id, RegistrationDate: &dateTime}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: `{
				  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
				  "registrationDate": "2025-04-12T20:10:25.102Z",
				  "city": "Москва" }`,
		},
		{
			name: "Success_With_City_ID",
			inputBody: `{  
					"id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					"city": "Москва" }`,
			inputPVZ:  &domain.PVZ{City: &allowedCity, ID: &id},
			outputPVZ: &domain.PVZ{City: &allowedCity, ID: &id, RegistrationDate: &dateTime},
			mockBehavior: func(p *mock_usecase.MockPVZUsecase, inputPVZ *domain.PVZ) {
				p.EXPECT().CreatePVZ(inputPVZ).Return(&domain.PVZ{City: &allowedCity, ID: &id, RegistrationDate: &dateTime}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: `{
					"id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					"registrationDate": "2025-04-12T20:10:25.102Z",
					"city": "Москва" }`,
		},
		{
			name: "Success_With_City_Date",
			inputBody: `{  
					"registrationDate": "2025-04-12T20:10:25.102Z",
					"city": "Москва" }`,
			inputPVZ:  &domain.PVZ{City: &allowedCity, RegistrationDate: &dateTime},
			outputPVZ: &domain.PVZ{City: &allowedCity, ID: &id, RegistrationDate: &dateTime},
			mockBehavior: func(p *mock_usecase.MockPVZUsecase, inputPVZ *domain.PVZ) {
				p.EXPECT().CreatePVZ(inputPVZ).Return(&domain.PVZ{City: &allowedCity, ID: &id, RegistrationDate: &dateTime}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: `{
				  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
				  "registrationDate": "2025-04-12T20:10:25.102Z",
				  "city": "Москва" }`,
		},
		{
			name:      "Success_With_City",
			inputBody: `{ "city": "Москва"}`,
			inputPVZ:  &domain.PVZ{City: &allowedCity},
			outputPVZ: &domain.PVZ{City: &allowedCity, ID: &id, RegistrationDate: &dateTime},
			mockBehavior: func(p *mock_usecase.MockPVZUsecase, inputPVZ *domain.PVZ) {
				p.EXPECT().CreatePVZ(inputPVZ).Return(&domain.PVZ{City: &allowedCity, ID: &id, RegistrationDate: &dateTime}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: `{
				  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
				  "registrationDate": "2025-04-12T20:10:25.102Z",
				  "city": "Москва" }`,
		},
		{
			name: "Failure_Without_City",
			inputBody: `{  
					"id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					"registrationDate": "2025-04-12T20:10:25.102Z" }`,
			inputPVZ:  &domain.PVZ{ID: &id, RegistrationDate: &dateTime},
			outputPVZ: nil,
			mockBehavior: func(p *mock_usecase.MockPVZUsecase, input *domain.PVZ) {
				p.EXPECT().CreatePVZ(input).Return(nil, domain.ErrInvalidInputData)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{ "message":"invalid input data" }`,
		},
		{
			name: "Failure_With_Wrong_City",
			inputBody: `{  
					"city" : "Париж",
					"id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					"registrationDate": "2025-04-12T20:10:25.102Z" }`,
			inputPVZ:  &domain.PVZ{ID: &id, RegistrationDate: &dateTime, City: &wrongCity},
			outputPVZ: nil,
			mockBehavior: func(p *mock_usecase.MockPVZUsecase, input *domain.PVZ) {
				p.EXPECT().CreatePVZ(input).Return(nil, domain.ErrInvalidCity)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{ "message":"invalid city" }`,
		},
		{
			name:      "Failure_With_Empty_InputData",
			inputBody: `{}`,
			inputPVZ:  &domain.PVZ{},
			outputPVZ: nil,
			mockBehavior: func(p *mock_usecase.MockPVZUsecase, input *domain.PVZ) {
				p.EXPECT().CreatePVZ(input).Return(nil, domain.ErrInvalidCity)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{ "message":"invalid city" }`,
		},
		{
			name:      "Failure_Wrong_Json_Format",
			inputBody: `{"}`,
			inputPVZ:  nil,
			outputPVZ: nil,
			mockBehavior: func(p *mock_usecase.MockPVZUsecase, input *domain.PVZ) {
				p.EXPECT().CreatePVZ(input).Times(0)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{ "message":"invalid input data" }`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockPVZUsecase(c)
			test.mockBehavior(repo, test.inputPVZ)

			services := usecase.Usecase{PVZUsecase: repo}
			handler := controller.NewHandler(services)

			r := gin.New()
			r.POST("/pvz", handler.CreatePVZ)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/pvz",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestGetPVZList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_usecase.NewMockPVZUsecase(ctrl)

	services := usecase.Usecase{PVZUsecase: mockService}
	handler := controller.NewHandler(services)

	t.Run("Successful_Get_Pvz_List", func(t *testing.T) {
		router := gin.Default()
		router.GET("/pvz", handler.GetPVZList)

		mockService.EXPECT().GetPVZInfo("2023-01-01", "2023-12-31", "1", "10").Return([]domain.PVZWithReceptions{}, nil)

		req, _ := http.NewRequest("GET", "/pvz?startDate=2023-01-01&endDate=2023-12-31&page=1&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
