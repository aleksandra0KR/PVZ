package implementation

import (
	"final/internal/domain"
	"final/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockProductRepository(ctrl)
	useCase := NewProductUseCase(mockRepo)

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

	t.Run("Successful_Create_Product", func(t *testing.T) {
		inputProduct := &domain.InputProduct{
			Type:  &productType,
			PvzId: &pvzId,
		}

		product := &domain.Product{
			ID:          &id,
			Type:        &productType,
			DateTime:    &dateTime,
			ReceptionID: &receptionId,
		}

		mockRepo.EXPECT().CreateProduct(inputProduct).Return(product, nil)

		result, err := useCase.CreateProduct(inputProduct)
		assert.NoError(t, err)
		assert.Equal(t, product, result)
	})

	t.Run("Failure_Invalid_Type", func(t *testing.T) {
		inputProduct := &domain.InputProduct{
			Type:  &wrongProductType,
			PvzId: &pvzId,
		}

		result, err := useCase.CreateProduct(inputProduct)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidProductType, err)
		assert.Nil(t, result)
	})

	t.Run("Failure_Empty_Type", func(t *testing.T) {
		inputProduct := &domain.InputProduct{
			Type:  nil,
			PvzId: &pvzId,
		}

		result, err := useCase.CreateProduct(inputProduct)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidInputData, err)
		assert.Nil(t, result)
	})

	t.Run("Failure_Empty_PvzId", func(t *testing.T) {
		inputProduct := &domain.InputProduct{
			Type:  &productType,
			PvzId: nil,
		}

		result, err := useCase.CreateProduct(inputProduct)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrEmptyPvzID, err)
		assert.Nil(t, result)
	})
}

func TestDeleteLastProductForPVZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockProductRepository(ctrl)
	useCase := NewProductUseCase(mockRepo)

	t.Run("Successful_Delete_Product", func(t *testing.T) {
		pvzId := "1"
		mockRepo.EXPECT().DeleteLastProductForPVZ(&pvzId).Return(nil)
		err := useCase.DeleteLastProductForPVZ(&pvzId)
		assert.NoError(t, err)
	})

	t.Run("Failure_Empty_pvzId", func(t *testing.T) {
		var pvzId *string
		err := useCase.DeleteLastProductForPVZ(pvzId)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrEmptyPvzID, err)
	})

	t.Run("Failure_Nil_pvzId", func(t *testing.T) {
		err := useCase.DeleteLastProductForPVZ(nil)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrEmptyPvzID, err)
	})
}
