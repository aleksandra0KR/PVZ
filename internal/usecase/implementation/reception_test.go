package implementation

import (
	"context"
	"final/internal/domain"
	"final/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	PvzId := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	layout := "2006-01-02T15:04:05.000Z"
	str := "2025-04-13T00:46:47.524Z"
	status := "in_progress"
	dateTime, err := time.Parse(layout, str)

	if err != nil {
		log.Error(err)
	}

	mockRepo := mock_repository.NewMockReceptionRepository(ctrl)
	useCase := NewReceptionUseCase(mockRepo)

	ctx := context.Background()
	t.Run("Successful_Create_Reception", func(t *testing.T) {
		reception := &domain.Reception{PvzId: &PvzId}
		outputReception := &domain.Reception{PvzId: &PvzId, ID: &id, DateTime: &dateTime, Status: &status}

		mockRepo.EXPECT().CreateReception(ctx, reception).Return(outputReception, nil)

		result, err := useCase.CreateReception(ctx, reception)
		assert.NoError(t, err)
		assert.Equal(t, outputReception, result)
	})

	t.Run("Failure_nil_pvzId", func(t *testing.T) {
		reception := &domain.Reception{}

		result, err := useCase.CreateReception(ctx, reception)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidInputData, err)
		assert.Nil(t, result)
	})
}

func TestCloseReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	PvzId := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	layout := "2006-01-02T15:04:05.000Z"
	str := "2025-04-13T00:46:47.524Z"
	status := "in_progress"
	dateTime, err := time.Parse(layout, str)

	if err != nil {
		log.Error(err)
	}
	mockRepo := mock_repository.NewMockReceptionRepository(ctrl)
	useCase := NewReceptionUseCase(mockRepo)

	ctx := context.Background()
	t.Run("Successful_Close_Reception", func(t *testing.T) {
		outputReception := &domain.Reception{PvzId: &PvzId, ID: &id, DateTime: &dateTime, Status: &status}

		mockRepo.EXPECT().CloseReception(ctx, &PvzId).Return(outputReception, nil)

		result, err := useCase.CloseReception(ctx, &PvzId)
		assert.NoError(t, err)
		assert.Equal(t, outputReception, result)
	})

	t.Run("Failure_nil_pvzId", func(t *testing.T) {
		var pvzID *string

		result, err := useCase.CloseReception(ctx, pvzID)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidInputData, err)
		assert.Nil(t, result)
	})

	t.Run("Failure_Empty_pvzId", func(t *testing.T) {
		pvzID := ""

		result, err := useCase.CloseReception(ctx, &pvzID)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidInputData, err)
		assert.Nil(t, result)
	})
}
