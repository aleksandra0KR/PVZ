package controller_test

import (
	"bytes"
	"context"
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

func TestHandler_createProduct(t *testing.T) {
	type mockBehavior func(p *mock_usecase.MockProductUsecase, input *domain.InputProduct)
	productType := "электроника"
	wrongProductType := "косметика"
	pvzId := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	receptionId := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	layout := "2006-01-02T15:04:05.000Z"
	str := "2025-04-13T01:28:43.606Z"
	dateTime, err := time.Parse(layout, str)

	if err != nil {
		log.Error(err)
	}
	ctx := context.Background()
	tests := []struct {
		name                 string
		inputBody            string
		inputProduct         *domain.InputProduct
		outputProduct        *domain.Product
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Success_With_All_Parameters",
			inputBody: `{
				  "type": "электроника",
				  "pvzId": "3fa85f64-5717-4562-b3fc-2c963f66afa6" }`,
			inputProduct:  &domain.InputProduct{Type: &productType, PvzId: &pvzId},
			outputProduct: &domain.Product{Type: &productType, ID: &id, DateTime: &dateTime, ReceptionID: &receptionId},
			mockBehavior: func(p *mock_usecase.MockProductUsecase, inputProduct *domain.InputProduct) {
				p.EXPECT().CreateProduct(ctx, inputProduct).Return(&domain.Product{Type: &productType, ID: &id, DateTime: &dateTime, ReceptionID: &receptionId}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: `{
				  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
				  "dateTime": "2025-04-13T01:28:43.606Z",
				  "type": "электроника",
				  "receptionId": "3fa85f64-5717-4562-b3fc-2c963f66afa6" }`,
		},
		{
			name: "Failure_With_Wrong_Type",
			inputBody: `{
				  "type": "косметика",
				  "pvzId": "3fa85f64-5717-4562-b3fc-2c963f66afa6" }`,
			inputProduct:  &domain.InputProduct{Type: &wrongProductType, PvzId: &pvzId},
			outputProduct: nil,
			mockBehavior: func(p *mock_usecase.MockProductUsecase, inputProduct *domain.InputProduct) {
				p.EXPECT().CreateProduct(ctx, inputProduct).Return(nil, domain.ErrInvalidProductType)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid product type"}`,
		},
		{
			name: "Failure_With_Empty_Type",
			inputBody: `{
				  "pvzId": "3fa85f64-5717-4562-b3fc-2c963f66afa6" }`,
			inputProduct:  &domain.InputProduct{Type: nil, PvzId: &pvzId},
			outputProduct: nil,
			mockBehavior: func(p *mock_usecase.MockProductUsecase, inputProduct *domain.InputProduct) {
				p.EXPECT().CreateProduct(ctx, inputProduct).Return(nil, domain.ErrInvalidProductType)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid product type"}`,
		},
		{
			name:          "Failure_With_Empty_PvzId",
			inputBody:     `{  "type": "электроника"}`,
			inputProduct:  &domain.InputProduct{Type: &productType},
			outputProduct: nil,
			mockBehavior: func(p *mock_usecase.MockProductUsecase, inputProduct *domain.InputProduct) {
				p.EXPECT().CreateProduct(ctx, inputProduct).Return(nil, domain.ErrEmptyPvzID)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"empty pvz id"}`,
		},
		{
			name: "Failure_Wrong_Json_Format",
			inputBody: `{
				  "type": "косметика,
				  "pvzId": "3fa85f64-5717-4562-b3fc-2c963f66afa6" }`,
			inputProduct:  nil,
			outputProduct: nil,
			mockBehavior: func(p *mock_usecase.MockProductUsecase, inputProduct *domain.InputProduct) {
				p.EXPECT().CreateProduct(ctx, inputProduct).Times(0)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"invalid input data"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockProductUsecase(c)
			test.mockBehavior(repo, test.inputProduct)

			services := usecase.Usecase{ProductUsecase: repo}
			handler := controller.NewHandler(services)

			r := gin.New()
			r.POST("/products", handler.CreateProduct)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/products",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteReception(t *testing.T) {
	type mockBehavior func(p *mock_usecase.MockProductUsecase, input *string)
	PvzId := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	emptyId := ""

	ctx := context.Background()
	tests := []struct {
		name                 string
		inputParam           string
		inputPvzID           *string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "Success",
			inputParam: PvzId,
			inputPvzID: &PvzId,

			mockBehavior: func(p *mock_usecase.MockProductUsecase, input *string) {
				p.EXPECT().DeleteLastProductForPVZ(ctx, input).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: ``},
		{
			name:       "Failure_With_Empty_PvzId",
			inputParam: emptyId,
			inputPvzID: &emptyId,
			mockBehavior: func(p *mock_usecase.MockProductUsecase, inputProduct *string) {
				p.EXPECT().DeleteLastProductForPVZ(ctx, inputProduct).Return(domain.ErrEmptyPvzID)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"empty pvz id"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_usecase.NewMockProductUsecase(c)
			test.mockBehavior(repo, test.inputPvzID)

			services := usecase.Usecase{ProductUsecase: repo}
			handler := controller.NewHandler(services)

			r := gin.New()
			r.POST("/pvz/:pvzId/delete_last_product", handler.DeleteLastProductForPVZ)

			w := httptest.NewRecorder()
			requestURL := fmt.Sprintf("/pvz/%s/delete_last_product", test.inputParam)
			req := httptest.NewRequest("POST", requestURL, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			if test.expectedResponseBody == `` {
				assert.Equal(t, 0, w.Body.Len())

			} else {
				assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
			}
		})
	}
}
