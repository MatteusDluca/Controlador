package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"controlador/backend/internal/models"
	"controlador/backend/internal/repositories"

)

var ErrTransacaoJaEstornada = errors.New("transação já foi estornada")

type ReverseTransacaoService struct {
	db            *pgxpool.Pool
	transacaoRepo repositories.TransacaoRepository
	ativoRepo     repositories.AtivoRepository
}

func NewReverseTransacaoService(db *pgxpool.Pool, tRepo repositories.TransacaoRepository, aRepo repositories.AtivoRepository) *ReverseTransacaoService {
	return &ReverseTransacaoService{db: db, transacaoRepo: tRepo, ativoRepo: aRepo}
}

func (s *ReverseTransacaoService) Execute(ctx context.Context, transacaoID string) (*models.Transacao, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil { return nil, err }
	defer tx.Rollback(ctx)

	original, err := s.transacaoRepo.FindByID(ctx, transacaoID)
	if err != nil { return nil, err }
	if original == nil { return nil, errors.New("transação original não encontrada") }
	if original.ReversalOf != nil { return nil, ErrTransacaoJaEstornada }

	estorno := &models.Transacao{
		ID:                uuid.New().String(),
		AtivoFinanceiroID: original.AtivoFinanceiroID,
		Descricao:         fmt.Sprintf("Estorno de: %s", original.Descricao),
		Valor:             original.Valor,
		Tipo:              models.TransacaoEstorno,
		ReversalOf:        &original.ID,
		CreatedAt:         time.Now(),
	}
	
	if err := s.transacaoRepo.Create(ctx, tx, estorno); err != nil { return nil, err }
	if err := s.ativoRepo.UpdateBalance(ctx, tx, estorno.AtivoFinanceiroID, estorno.Valor, estorno.Tipo); err != nil { return nil, err }

	return estorno, tx.Commit(ctx)
}