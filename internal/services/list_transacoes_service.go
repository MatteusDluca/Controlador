package services

import (
	"context"

	"controlador/backend/internal/models"
	"controlador/backend/internal/repositories"

)

type ListTransacoesService struct {
	repo repositories.TransacaoRepository
}

func NewListTransacoesService(repo repositories.TransacaoRepository) *ListTransacoesService {
	return &ListTransacoesService{repo: repo}
}

func (s *ListTransacoesService) Execute(ctx context.Context) ([]models.Transacao, error) {
	return s.repo.FindAll(ctx)
}