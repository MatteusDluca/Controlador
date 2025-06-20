package services

import (
	"context"

	"controlador/backend/internal/models"
	"controlador/backend/internal/repositories"

)

type ListCategoriasService struct {
	repo repositories.CategoriaRepository
}

func NewListCategoriasService(repo repositories.CategoriaRepository) *ListCategoriasService {
	return &ListCategoriasService{repo: repo}
}

// Execute retorna a lista de todas as categorias.
func (s *ListCategoriasService) Execute(ctx context.Context) ([]models.Categoria, error) {
	return s.repo.FindAll(ctx)
}