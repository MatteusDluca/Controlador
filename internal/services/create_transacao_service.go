package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool" // <-- Import necessário

	"controlador/backend/internal/models"
	"controlador/backend/internal/repositories"

)

var (
	ErrAtivoNaoEncontrado   = errors.New("ativo financeiro não encontrado")
	ErrSaldoInsuficiente    = errors.New("saldo ou limite insuficiente para a transação")
	ErrTipoTransacaoInvalido  = errors.New("tipo de transação inválido ou incompatível com o ativo")
)

type CreateTransacaoService struct {
	db            *pgxpool.Pool // <-- CORREÇÃO: Adicionado o pool de conexão
	transacaoRepo repositories.TransacaoRepository
	ativoRepo     repositories.AtivoRepository
}

// CORREÇÃO: O construtor agora aceita o pool de conexão
func NewCreateTransacaoService(db *pgxpool.Pool, tRepo repositories.TransacaoRepository, aRepo repositories.AtivoRepository) *CreateTransacaoService {
	return &CreateTransacaoService{
		db:            db,
		transacaoRepo: tRepo,
		ativoRepo:     aRepo,
	}
}

func (s *CreateTransacaoService) Execute(ctx context.Context, input models.Transacao) (*models.Transacao, error) {
	// Inicia uma transação de banco de dados. Essencial para consistência.
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	// Garante que a transação seja desfeita (ROLLBACK) se algo der errado.
	defer tx.Rollback(ctx)

	// 1. Validar o ativo
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

	// 2. Validar o tipo de transação e o saldo/limite
	switch input.Tipo {
	case models.TransacaoDebito:
		if ativo.Tipo != models.AtivoContaCorrente {
			return nil, ErrTipoTransacaoInvalido
		}
		if ativo.SaldoAtual < input.Valor {
			return nil, ErrSaldoInsuficiente
		}
	case models.TransacaoCredito:
		if ativo.Tipo != models.AtivoCartaoCredito {
			return nil, ErrTipoTransacaoInvalido
		}
		if ativo.LimiteDisponivel < input.Valor {
			return nil, ErrSaldoInsuficiente
		}
	default:
		return nil, ErrTipoTransacaoInvalido
	}

	// 3. Preparar a transação
	input.ID = uuid.New().String()
	input.CreatedAt = time.Now()

	// 4. Chamar os repositórios, passando a transação (tx)
	if err := s.transacaoRepo.Create(ctx, tx, &input); err != nil {
		return nil, err
	}
	if err := s.ativoRepo.UpdateBalance(ctx, tx, input.AtivoFinanceiroID, input.Valor, input.Tipo); err != nil {
		return nil, err
	}

	// Se tudo correu bem, confirma a transação (COMMIT).
	return &input, tx.Commit(ctx)
}