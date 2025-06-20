package repositories

import (
	"context"

	"controlador/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransacaoRecorrenteRepository interface {
	Create(ctx context.Context, tr *models.TransacaoRecorrente) error
	FindByID(ctx context.Context, id string) (*models.TransacaoRecorrente, error)
	FindAllByAtivoID(ctx context.Context, ativoID string) ([]models.TransacaoRecorrente, error)
	FindActiveByDay(ctx context.Context, dia int) ([]models.TransacaoRecorrente, error)
	Update(ctx context.Context, tr *models.TransacaoRecorrente) error
	Delete(ctx context.Context, id string) error
}

type pgTransacaoRecorrenteRepository struct {
	db *pgxpool.Pool
}

func NewPgTransacaoRecorrenteRepository(db *pgxpool.Pool) TransacaoRecorrenteRepository {
	return &pgTransacaoRecorrenteRepository{db: db}
}

func (r *pgTransacaoRecorrenteRepository) Create(ctx context.Context, tr *models.TransacaoRecorrente) error {
	sql := `
		INSERT INTO transacoes_recorrentes 
		(id, ativo_financeiro_id, categoria_id, descricao, valor, tipo, dia_do_vencimento, ativa, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Exec(ctx, sql, tr.ID, tr.AtivoFinanceiroID, tr.CategoriaID, tr.Descricao, tr.Valor, tr.Tipo, tr.DiaDoVencimento, tr.Ativa, tr.CreatedAt, tr.UpdatedAt)
	return err
}

func (r *pgTransacaoRecorrenteRepository) FindByID(ctx context.Context, id string) (*models.TransacaoRecorrente, error) {
	var tr models.TransacaoRecorrente
	sql := `
		SELECT id, ativo_financeiro_id, categoria_id, descricao, valor, tipo, dia_do_vencimento, ativa, created_at, updated_at 
		FROM transacoes_recorrentes WHERE id = $1`
	err := r.db.QueryRow(ctx, sql, id).Scan(
		&tr.ID, &tr.AtivoFinanceiroID, &tr.CategoriaID, &tr.Descricao, &tr.Valor, &tr.Tipo, &tr.DiaDoVencimento, &tr.Ativa, &tr.CreatedAt, &tr.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &tr, nil
}

func (r *pgTransacaoRecorrenteRepository) FindAllByAtivoID(ctx context.Context, ativoID string) ([]models.TransacaoRecorrente, error) {
	var recorrentes []models.TransacaoRecorrente
	sql := `
		SELECT id, ativo_financeiro_id, categoria_id, descricao, valor, tipo, dia_do_vencimento, ativa, created_at, updated_at 
		FROM transacoes_recorrentes WHERE ativo_financeiro_id = $1 ORDER BY dia_do_vencimento ASC`
	rows, err := r.db.Query(ctx, sql, ativoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tr models.TransacaoRecorrente
		if err := rows.Scan(
			&tr.ID, &tr.AtivoFinanceiroID, &tr.CategoriaID, &tr.Descricao, &tr.Valor, &tr.Tipo, &tr.DiaDoVencimento, &tr.Ativa, &tr.CreatedAt, &tr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		recorrentes = append(recorrentes, tr)
	}
	return recorrentes, nil
}

func (r *pgTransacaoRecorrenteRepository) FindActiveByDay(ctx context.Context, dia int) ([]models.TransacaoRecorrente, error) {
	var recorrentes []models.TransacaoRecorrente
	sql := `
		SELECT id, ativo_financeiro_id, categoria_id, descricao, valor, tipo, dia_do_vencimento, ativa, created_at, updated_at 
		FROM transacoes_recorrentes WHERE dia_do_vencimento = $1 AND ativa = TRUE`
	rows, err := r.db.Query(ctx, sql, dia)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tr models.TransacaoRecorrente
		if err := rows.Scan(
			&tr.ID, &tr.AtivoFinanceiroID, &tr.CategoriaID, &tr.Descricao, &tr.Valor, &tr.Tipo, &tr.DiaDoVencimento, &tr.Ativa, &tr.CreatedAt, &tr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		recorrentes = append(recorrentes, tr)
	}
	return recorrentes, nil
}

func (r *pgTransacaoRecorrenteRepository) Update(ctx context.Context, tr *models.TransacaoRecorrente) error {
	sql := `
		UPDATE transacoes_recorrentes SET 
		ativo_financeiro_id = $1, categoria_id = $2, descricao = $3, valor = $4, tipo = $5, dia_do_vencimento = $6, ativa = $7, updated_at = $8
		WHERE id = $9`
	_, err := r.db.Exec(ctx, sql, tr.AtivoFinanceiroID, tr.CategoriaID, tr.Descricao, tr.Valor, tr.Tipo, tr.DiaDoVencimento, tr.Ativa, tr.UpdatedAt, tr.ID)
	return err
}

func (r *pgTransacaoRecorrenteRepository) Delete(ctx context.Context, id string) error {
	sql := `DELETE FROM transacoes_recorrentes WHERE id = $1`
	_, err := r.db.Exec(ctx, sql, id)
	return err
}