package services

import (
	"context"

	"controlador/backend/internal/models"
	"controlador/backend/internal/repositories"

)

type ListAtivosService struct {
	repo repositories.AtivoRepository
}

func NewListAtivosService(repo repositories.AtivoRepository) *ListAtivosService {
	return &ListAtivosService{repo: repo}
}

func (s *ListAtivosService) Execute(ctx context.Context) ([]models.AtivoFinanceiro, error) {
	return s.repo.FindAll(ctx)
}