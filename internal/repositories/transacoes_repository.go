package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"controlador/backend/internal/models"

)

type TransacaoRepository interface {
	Create(ctx context.Context, tx pgx.Tx, transacao *models.Transacao) error
	FindAll(ctx context.Context) ([]models.Transacao, error)
	FindByID(ctx context.Context, id string) (*models.Transacao, error)
}

type pgTransacaoRepository struct {
	db *pgxpool.Pool
}

func NewPgTransacaoRepository(db *pgxpool.Pool) TransacaoRepository {
	return &pgTransacaoRepository{db: db}
}

func (r *pgTransacaoRepository) Create(ctx context.Context, tx pgx.Tx, transacao *models.Transacao) error {
	sql := `INSERT INTO transacoes (id, ativo_financeiro_id, descricao, valor, tipo, created_at, reversal_of) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := tx.Exec(ctx, sql, transacao.ID, transacao.AtivoFinanceiroID, transacao.Descricao, transacao.Valor, transacao.Tipo, transacao.CreatedAt, transacao.ReversalOf)
	return err
}

func (r *pgTransacaoRepository) FindByID(ctx context.Context, id string) (*models.Transacao, error) {
	var t models.Transacao
	sql := `SELECT id, ativo_financeiro_id, descricao, valor, tipo, created_at, reversal_of FROM transacoes WHERE id = $1`
	err := r.db.QueryRow(ctx, sql).Scan(&t.ID, &t.AtivoFinanceiroID, &t.Descricao, &t.Valor, &t.Tipo, &t.CreatedAt, &t.ReversalOf)
	if err != nil {
		if err == pgx.ErrNoRows { return nil, nil }
		return nil, err
	}
	return &t, nil
}

func (r *pgTransacaoRepository) FindAll(ctx context.Context) ([]models.Transacao, error) {
	var transacoes []models.Transacao
	sql := `SELECT id, ativo_financeiro_id, descricao, valor, tipo, created_at, reversal_of FROM transacoes ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, sql)
	if err != nil { return nil, err }
	defer rows.Close()
	for rows.Next() {
		var t models.Transacao
		if err := rows.Scan(&t.ID, &t.AtivoFinanceiroID, &t.Descricao, &t.Valor, &t.Tipo, &t.CreatedAt, &t.ReversalOf); err != nil {
			return nil, err
		}
		transacoes = append(transacoes, t)
	}
	return transacoes, nil
}