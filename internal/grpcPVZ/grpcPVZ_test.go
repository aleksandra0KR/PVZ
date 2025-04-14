package grpcPVZ

import (
	"context"
	"final/internal/domain"
	"final/internal/repository"
	"final/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetPVZList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPVZRepository(ctrl)
	service := NewPVZService(&repository.Repository{PvzRepository: mockRepo})
	allowedCity := "Москва"
	id := "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	id2 := "3fa85f64-5717-4562-b3fc-2c963f66afa4"

	layout := "2006-01-02T15:04:05.000Z"
	str := "2025-04-12T20:10:25.102Z"
	dateTime, err := time.Parse(layout, str)
	if err != nil {
		log.Error(err)
	}
	t.Run("Successful_Get_Pvz", func(t *testing.T) {
		pvz1 := &domain.PVZ{
			ID:               &id,
			City:             &allowedCity,
			RegistrationDate: &dateTime,
		}
		pvz2 := &domain.PVZ{
			ID:               &id2,
			City:             &allowedCity,
			RegistrationDate: &dateTime,
		}

		mockRepo.EXPECT().GetAllPVZ().Return([]*domain.PVZ{pvz1, pvz2}, nil)

		resp, err := service.GetPVZList(context.Background(), &GetPVZListRequest{})
		assert.NoError(t, err)
		assert.Len(t, resp.Pvzs, 2)
		assert.Equal(t, id, resp.Pvzs[0].Id)
		assert.Equal(t, allowedCity, resp.Pvzs[0].City)
		assert.Equal(t, id2, resp.Pvzs[1].Id)
		assert.Equal(t, allowedCity, resp.Pvzs[1].City)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().GetAllPVZ().Return(nil, domain.ErrGetPVZ)

		resp, err := service.GetPVZList(context.Background(), &GetPVZListRequest{})
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestGetPVZList_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPVZRepository(ctrl)
	mockRepo.EXPECT().GetAllPVZ().Return(nil, domain.ErrGetPVZ)

	repo := &repository.Repository{PvzRepository: mockRepo}
	service := NewPVZService(repo)

	resp, err := service.GetPVZList(context.Background(), &GetPVZListRequest{})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrGetPVZ.Error(), err.Error())
}

func TestGetPVZList_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPVZRepository(ctrl)
	mockRepo.EXPECT().GetAllPVZ().Return([]*domain.PVZ{}, nil)

	repo := &repository.Repository{PvzRepository: mockRepo}
	service := NewPVZService(repo)

	resp, err := service.GetPVZList(context.Background(), &GetPVZListRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.Pvzs)
}

func TestGetPVZList_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockPVZRepository(ctrl)
	mockRepo.EXPECT().GetAllPVZ().Return([]*domain.PVZ{}, nil)

	repo := &repository.Repository{PvzRepository: mockRepo}
	service := NewPVZService(repo)
	resp, err := service.GetPVZList(context.Background(), &GetPVZListRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.Pvzs)
}
