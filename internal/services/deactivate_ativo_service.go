package services

import (
	"context"

	"controlador/backend/internal/repositories"

)

type DeactivateAtivoService struct {
	repo repositories.AtivoRepository
}

func NewDeactivateAtivoService(repo repositories.AtivoRepository) *DeactivateAtivoService {
	return &DeactivateAtivoService{repo: repo}
}

func (s *DeactivateAtivoService) Execute(ctx context.Context, id string) error {
	ativo, err := s.repo.FindByID(ctx, id)
	if err != nil { return err }
	if ativo == nil { return ErrAtivoNaoEncontrado }

	return s.repo.Deactivate(ctx, id)
}