package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"controlador/backend/internal/models"
	"controlador/backend/internal/repositories"

)

var (
	ErrAtivoNaoEncontrado     = errors.New("ativo financeiro não encontrado")
	ErrSaldoInsuficiente      = errors.New("saldo ou limite insuficiente para a transação")
	ErrTipoTransacaoInvalido  = errors.New("tipo de transação inválido ou incompatível com o ativo")
	ErrCategoriaNaoEncontrada = errors.New("categoria não encontrada")
)

type CreateTransacaoService struct {
	db            *pgxpool.Pool
	transacaoRepo repositories.TransacaoRepository
	ativoRepo     repositories.AtivoRepository
	categoriaRepo repositories.CategoriaRepository
}

func NewCreateTransacaoService(db *pgxpool.Pool, tRepo repositories.TransacaoRepository, aRepo repositories.AtivoRepository, cRepo repositories.CategoriaRepository) *CreateTransacaoService {
	return &CreateTransacaoService{
		db:            db,
		transacaoRepo: tRepo,
		ativoRepo:     aRepo,
		categoriaRepo: cRepo,
	}
}

func (s *CreateTransacaoService) Execute(ctx context.Context, input models.Transacao) (*models.Transacao, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 1. Validar a categoria
	categoria, err := s.categoriaRepo.FindByID(ctx, input.CategoriaID)
	if err != nil {
		return nil, err
	}
	if categoria == nil {
		return nil, ErrCategoriaNaoEncontrada
	}

	// 2. Validar o ativo
	ativo, err := s.ativoRepo.FindByID(ctx, input.AtivoFinanceiroID)
	if err != nil {
		return nil, err
	}
	if ativo == nil {
		return nil, ErrAtivoNaoEncontrado
	}
	if !ativo.IsActive {
		return nil, errors.New("ativo financeiro está desativado")
	}

	// 3. Validar o tipo de transação e o saldo/limite
	// ALTERAÇÃO: Lógica de validação refatorada com o novo tipo 'RECEBIMENTO'.
	switch input.Tipo {
	case models.TransacaoRecebimento:
		// Um recebimento só pode ocorrer em uma conta corrente.
		if ativo.Tipo != models.AtivoContaCorrente {
			return nil, ErrTipoTransacaoInvalido
		}
		// Nenhuma verificação de saldo é necessária para recebimentos.
	case models.TransacaoDebito:
		// Um débito (gasto) só pode ocorrer em uma conta corrente.
		if ativo.Tipo != models.AtivoContaCorrente {
			return nil, ErrTipoTransacaoInvalido
		}
		if ativo.SaldoAtual < input.Valor {
			return nil, ErrSaldoInsuficiente
		}
	case models.TransacaoCredito:
		// Um crédito (gasto) só pode ocorrer em um cartão de crédito.
		if ativo.Tipo != models.AtivoCartaoCredito {
			return nil, ErrTipoTransacaoInvalido
		}
		if ativo.LimiteDisponivel < input.Valor {
			return nil, ErrSaldoInsuficiente
		}
	default:
		return nil, ErrTipoTransacaoInvalido
	}

	// 4. Preparar a transação
	input.ID = uuid.New().String()
	input.CreatedAt = time.Now()

	// 5. Chamar os repositórios, passando a transação (tx)
	if err := s.transacaoRepo.Create(ctx, tx, &input); err != nil {
		return nil, err
	}
	if err := s.ativoRepo.UpdateBalance(ctx, tx, input.AtivoFinanceiroID, input.Valor, input.Tipo); err != nil {
		return nil, err
	}

	return &input, tx.Commit(ctx)
}