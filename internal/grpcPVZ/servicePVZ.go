package grpcPVZ

import (
	"context"
	"final/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PVZService struct {
	UnimplementedPVZServiceServer
	repo *repository.Repository
}

func NewPVZService(repo *repository.Repository) *PVZService {
	return &PVZService{repo: repo}
}

func (s *PVZService) GetPVZList(_ context.Context, _ *GetPVZListRequest) (*GetPVZListResponse, error) {
	pvzs, err := s.repo.PvzRepository.GetAllPVZ()
	if err != nil {
		return nil, err
	}

	var grpcPVZs []*PVZ
	for _, pvz := range pvzs {
		grpcPVZs = append(grpcPVZs, &PVZ{
			Id:               *pvz.ID,
			City:             *pvz.City,
			RegistrationDate: timestamppb.New(*pvz.RegistrationDate),
		})
	}
	return &GetPVZListResponse{Pvzs: grpcPVZs}, nil
}
