package services

import (
	"context"

	"controlador/backend/internal/models"
	"controlador/backend/internal/repositories"

)

type ListTransacoesRecorrentesService struct {
	repo repositories.TransacaoRecorrenteRepository
}

func NewListTransacoesRecorrentesService(repo repositories.TransacaoRecorrenteRepository) *ListTransacoesRecorrentesService {
	return &ListTransacoesRecorrentesService{repo: repo}
}

func (s *ListTransacoesRecorrentesService) Execute(ctx context.Context, ativoID string) ([]models.TransacaoRecorrente, error) {
	return s.repo.FindAllByAtivoID(ctx, ativoID)
}