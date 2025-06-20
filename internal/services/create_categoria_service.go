package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"controlador/backend/internal/models"
	"controlador/backend/internal/repositories"

)

var ErrCategoriaJaExiste = errors.New("uma categoria com este nome já existe")

type CreateCategoriaService struct {
	repo repositories.CategoriaRepository
}

func NewCreateCategoriaService(repo repositories.CategoriaRepository) *CreateCategoriaService {
	return &CreateCategoriaService{repo: repo}
}

// Execute orquestra a criação de uma nova categoria.
func (s *CreateCategoriaService) Execute(ctx context.Context, input models.Categoria) (*models.Categoria, error) {
	// 1. Validação da regra de negócio: Não permitir categorias com nomes duplicados.
	existing, err := s.repo.FindByName(ctx, input.Nome)
	if err != nil {
		// Erro ao consultar o banco de dados.
		return nil, err
	}
	if existing != nil {
		// Uma categoria com este nome já foi encontrada.
		return nil, ErrCategoriaJaExiste
	}

	// 2. Preparação do modelo para persistência.
	input.ID = uuid.New().String()

	// 3. Chamada ao repositório para salvar a nova categoria.
	if err := s.repo.Create(ctx, &input); err != nil {
		return nil, err
	}

	// 4. Retorna a categoria criada com sucesso.
	return &input, nil
}