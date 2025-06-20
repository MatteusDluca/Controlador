package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"controlador/backend/internal/models"
	"controlador/backend/internal/repositories"

)

type CreateAtivoService struct {
	repo repositories.AtivoRepository
}

func NewCreateAtivoService(repo repositories.AtivoRepository) *CreateAtivoService {
	return &CreateAtivoService{repo: repo}
}

func (s *CreateAtivoService) Execute(ctx context.Context, input models.AtivoFinanceiro) (*models.AtivoFinanceiro, error) {
	input.ID = uuid.New().String()
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now

	err := s.repo.Save(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &input, nil
}