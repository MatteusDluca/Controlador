package repositories

import (
	"context"

	"controlador/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AtivoRepository interface {
	Save(ctx context.Context, ativo *models.AtivoFinanceiro) error
	FindAll(ctx context.Context) ([]models.AtivoFinanceiro, error)
	FindByID(ctx context.Context, id string) (*models.AtivoFinanceiro, error)
	UpdateBalance(ctx context.Context, tx pgx.Tx, ativoID string, valor float64, tipo models.TipoTransacao) error
	Deactivate(ctx context.Context, id string) error
}

type pgAtivoRepository struct {
	db *pgxpool.Pool
}

func NewPgAtivoRepository(db *pgxpool.Pool) AtivoRepository {
	return &pgAtivoRepository{db: db}
}

func (r *pgAtivoRepository) Deactivate(ctx context.Context, id string) error {
	sql := `UPDATE ativos_financeiros SET is_active = FALSE, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, sql, id)
	return err
}

func (r *pgAtivoRepository) UpdateBalance(ctx context.Context, tx pgx.Tx, ativoID string, valor float64, tipo models.TipoTransacao) error {
	var sqlUpdate string
	if tipo == models.TransacaoRecebimento {
		sqlUpdate = `UPDATE ativos_financeiros SET saldo_atual = saldo_atual + $1, updated_at = NOW() WHERE id = $2`
	} else if tipo == models.TransacaoDebito {
		sqlUpdate = `UPDATE ativos_financeiros SET saldo_atual = saldo_atual - $1, updated_at = NOW() WHERE id = $2`
	} else if tipo == models.TransacaoCredito {
		sqlUpdate = `UPDATE ativos_financeiros SET limite_disponivel = limite_disponivel - $1, updated_at = NOW() WHERE id = $2`
	} else if tipo == models.TransacaoEstorno {
		ativo, err := r.FindByID(ctx, ativoID)
		if err != nil {
			return err
		}
		if ativo.Tipo == models.AtivoContaCorrente {
			sqlUpdate = `UPDATE ativos_financeiros SET saldo_atual = saldo_atual + $1, updated_at = NOW() WHERE id = $2`
		} else {
			sqlUpdate = `UPDATE ativos_financeiros SET limite_disponivel = limite_disponivel + $1, updated_at = NOW() WHERE id = $2`
		}
	}
	_, err := tx.Exec(ctx, sqlUpdate, valor, ativoID)
	return err
}

func (r *pgAtivoRepository) Save(ctx context.Context, ativo *models.AtivoFinanceiro) error {
	sql := `INSERT INTO ativos_financeiros (id, instituicao, nome, tipo, saldo_atual, limite_disponivel, created_at, updated_at, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE)`
	_, err := r.db.Exec(ctx, sql, ativo.ID, ativo.Instituicao, ativo.Nome, ativo.Tipo, ativo.SaldoAtual, ativo.LimiteDisponivel, ativo.CreatedAt, ativo.UpdatedAt)
	return err
}

func (r *pgAtivoRepository) FindAll(ctx context.Context) ([]models.AtivoFinanceiro, error) {
	var ativos []models.AtivoFinanceiro
	sql := `SELECT id, instituicao, nome, tipo, saldo_atual, limite_disponivel, is_active, created_at, updated_at FROM ativos_financeiros ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a models.AtivoFinanceiro
		if err := rows.Scan(&a.ID, &a.Instituicao, &a.Nome, &a.Tipo, &a.SaldoAtual, &a.LimiteDisponivel, &a.IsActive, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		ativos = append(ativos, a)
	}
	return ativos, nil
}

func (r *pgAtivoRepository) FindByID(ctx context.Context, id string) (*models.AtivoFinanceiro, error) {
	var ativo models.AtivoFinanceiro
	sql := `SELECT id, instituicao, nome, tipo, saldo_atual, limite_disponivel, is_active, created_at, updated_at FROM ativos_financeiros WHERE id = $1`
	err := r.db.QueryRow(ctx, sql, id).Scan(&ativo.ID, &ativo.Instituicao, &ativo.Nome, &ativo.Tipo, &ativo.SaldoAtual, &ativo.LimiteDisponivel, &ativo.IsActive, &ativo.CreatedAt, &ativo.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ativo, nil
}