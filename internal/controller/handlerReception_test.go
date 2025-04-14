package controller_test

import (
	"bytes"
	"final/internal/controller"
	"final/internal/domain"
	"final/internal/usecase"
	"final/internal/usecase/mocks"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_createReception(t *testing.T) {
	type mockBehavior func(p *mock_usecase.MockReceptionUsecase, reception *domain.Reception)
	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	PvzId := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	layout := "2006-01-02T15:04:05.000Z"
	str := "2025-04-13T00:46:47.524Z"
	status := "in_progress"
	dateTime, err := time.Parse(layout, str)

	if err != nil {
		log.Error(err)
	}

	tests := []struct {
		name                 string
		inputBody            string
		inputReception       *domain.Reception
		outputReception      *domain.Reception
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success",
			inputBody: `{
				"pvzId": "3fa85f64-5717-4562-b3fc-2c963f66afa6"
			}`,
			inputReception:  &domain.Reception{PvzId: &PvzId},
			outputReception: &domain.Reception{PvzId: &PvzId, ID: &id, DateTime: &dateTime, Status: &status},

			mockBehavior: func(p *mock_usecase.MockReceptionUsecase, input *domain.Reception) {
				p.EXPECT().CreateReception(input).Return(&domain.Reception{PvzId: &PvzId, ID: &id, DateTime: &dateTime, Status: &status}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: `{
					  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					  "dateTime": "2025-04-13T00:46:47.524Z",
					  "pvzId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					  "status": "in_progress"} `,
		},
		{
			name:            "Failure_With_Empty_PvzId",
			inputBody:       `{ }`,
			inputReception:  &domain.Reception{},
			outputReception: nil,

			mockBehavior: func(p *mock_usecase.MockReceptionUsecase, input *domain.Reception) {
				p.EXPECT().CreateReception(input).Return(nil, domain.ErrInvalidInputData)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{ "message": "invalid input data" }`,
		},
		{
			name:            "Failure_Wrong_Json_Format",
			inputBody:       `{ " }`,
			inputReception:  nil,
			outputReception: nil,
			mockBehavior: func(p *mock_usecase.MockReceptionUsecase, input *domain.Reception) {
				p.EXPECT().CreateReception(input).Times(0)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{ "message": "invalid input data" }`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockReceptionUsecase(c)
			test.mockBehavior(repo, test.inputReception)

			services := usecase.Usecase{ReceptionUsecase: repo}
			handler := controller.NewHandler(services)

			r := gin.New()
			r.POST("/receptions", handler.CreateReception)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/receptions", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_closeReception(t *testing.T) {
	type mockBehavior func(p *mock_usecase.MockReceptionUsecase, pvzID *string)
	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	PvzId := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	layout := "2006-01-02T15:04:05.000Z"
	str := "2025-04-13T00:46:47.524Z"
	status := "close"
	emptyId := ""
	dateTime, err := time.Parse(layout, str)

	if err != nil {
		log.Error(err)
	}

	tests := []struct {
		name                 string
		inputParam           string
		inputPvzID           *string
		outputReception      *domain.Reception
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:            "Success",
			inputParam:      PvzId,
			inputPvzID:      &PvzId,
			outputReception: &domain.Reception{PvzId: &PvzId, ID: &id, DateTime: &dateTime, Status: &status},

			mockBehavior: func(p *mock_usecase.MockReceptionUsecase, input *string) {
				p.EXPECT().CloseReception(input).Return(&domain.Reception{PvzId: &PvzId, ID: &id, DateTime: &dateTime, Status: &status}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{
					  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					  "dateTime": "2025-04-13T00:46:47.524Z",
					  "pvzId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					  "status": "close"} `,
		},
		{
			name:            "Failure_Without_PvzId",
			inputParam:      emptyId,
			inputPvzID:      &emptyId,
			outputReception: nil,

			mockBehavior: func(p *mock_usecase.MockReceptionUsecase, input *string) {
				p.EXPECT().CloseReception(input).Return(nil, domain.ErrInvalidInputData)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{ "message":"invalid input data" }`,
		},
		{
			name:            "Failure_Without_PvzId",
			inputParam:      emptyId,
			inputPvzID:      &emptyId,
			outputReception: nil,

			mockBehavior: func(p *mock_usecase.MockReceptionUsecase, input *string) {
				p.EXPECT().CloseReception(input).Return(nil, domain.ErrInvalidInputData)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{ "message":"invalid input data" }`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockReceptionUsecase(c)
			test.mockBehavior(repo, test.inputPvzID)

			services := usecase.Usecase{ReceptionUsecase: repo}
			handler := controller.NewHandler(services)

			r := gin.New()
			r.POST("/pvz/:pvzId/close_last_reception", handler.CloseLastReception)

			w := httptest.NewRecorder()
			requestURL := fmt.Sprintf("/pvz/%s/close_last_reception", test.inputParam)
			req := httptest.NewRequest("POST", requestURL, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
