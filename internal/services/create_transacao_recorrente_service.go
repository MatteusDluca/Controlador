package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"controlador/backend/internal/models"
	"controlador/backend/internal/repositories"

)

var ErrDiaInvalido = errors.New("o dia do vencimento deve estar entre 1 e 31")

type CreateTransacaoRecorrenteService struct {
	trRepo        repositories.TransacaoRecorrenteRepository
	ativoRepo     repositories.AtivoRepository
	categoriaRepo repositories.CategoriaRepository
}

func NewCreateTransacaoRecorrenteService(tr repositories.TransacaoRecorrenteRepository, ar repositories.AtivoRepository, cr repositories.CategoriaRepository) *CreateTransacaoRecorrenteService {
	return &CreateTransacaoRecorrenteService{
		trRepo:        tr,
		ativoRepo:     ar,
		categoriaRepo: cr,
	}
}

func (s *CreateTransacaoRecorrenteService) Execute(ctx context.Context, input models.TransacaoRecorrente) (*models.TransacaoRecorrente, error) {
	// 1. Validações de dados de entrada
	if input.DiaDoVencimento < 1 || input.DiaDoVencimento > 31 {
		return nil, ErrDiaInvalido
	}

	// 2. Validações de existência e compatibilidade
	ativo, err := s.ativoRepo.FindByID(ctx, input.AtivoFinanceiroID)
	if err != nil {
		return nil, err
	}
	if ativo == nil {
		return nil, ErrAtivoNaoEncontrado
	}

	categoria, err := s.categoriaRepo.FindByID(ctx, input.CategoriaID)
	if err != nil {
		return nil, err
	}
	if categoria == nil {
		return nil, ErrCategoriaNaoEncontrada
	}

	if (input.Tipo == models.TransacaoRecebimento || input.Tipo == models.TransacaoDebito) && ativo.Tipo != models.AtivoContaCorrente {
		return nil, ErrTipoTransacaoInvalido
	}

	// 3. Preparar o modelo para persistência
	input.ID = uuid.New().String()
	input.Ativa = true // Uma nova recorrência sempre começa ativa.
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now

	// 4. Salvar no banco de dados
	if err := s.trRepo.Create(ctx, &input); err != nil {
		return nil, err
	}

	return &input, nil
}