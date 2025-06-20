package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"controlador/backend/internal/models"

)

// CategoriaRepository define a interface para operações de persistência de Categoria.
type CategoriaRepository interface {
	Create(ctx context.Context, categoria *models.Categoria) error
	FindAll(ctx context.Context) ([]models.Categoria, error)
	FindByID(ctx context.Context, id string) (*models.Categoria, error)
	Update(ctx context.Context, categoria *models.Categoria) error
	Delete(ctx context.Context, id string) error
	FindByName(ctx context.Context, nome string) (*models.Categoria, error)
}

type pgCategoriaRepository struct {
	db *pgxpool.Pool
}

// NewPgCategoriaRepository cria uma nova instância de CategoriaRepository.
func NewPgCategoriaRepository(db *pgxpool.Pool) CategoriaRepository {
	return &pgCategoriaRepository{db: db}
}

func (r *pgCategoriaRepository) Create(ctx context.Context, categoria *models.Categoria) error {
	sql := `INSERT INTO categorias (id, nome, icone) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, sql, categoria.ID, categoria.Nome, categoria.Icone)
	return err
}

func (r *pgCategoriaRepository) FindAll(ctx context.Context) ([]models.Categoria, error) {
	var categorias []models.Categoria
	sql := `SELECT id, nome, icone FROM categorias ORDER BY nome ASC`
	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Categoria
		if err := rows.Scan(&c.ID, &c.Nome, &c.Icone); err != nil {
			return nil, err
		}
		categorias = append(categorias, c)
	}

	return categorias, nil
}

func (r *pgCategoriaRepository) FindByID(ctx context.Context, id string) (*models.Categoria, error) {
	var c models.Categoria
	sql := `SELECT id, nome, icone FROM categorias WHERE id = $1`
	err := r.db.QueryRow(ctx, sql, id).Scan(&c.ID, &c.Nome, &c.Icone)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Retorna nil, nil se não encontrar, para ser tratado no serviço.
		}
		return nil, err
	}
	return &c, nil
}

func (r *pgCategoriaRepository) FindByName(ctx context.Context, nome string) (*models.Categoria, error) {
	var c models.Categoria
	sql := `SELECT id, nome, icone FROM categorias WHERE nome = $1`
	err := r.db.QueryRow(ctx, sql, nome).Scan(&c.ID, &c.Nome, &c.Icone)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Retorna nil, nil se não encontrar, para ser tratado no serviço.
		}
		return nil, err
	}
	return &c, nil
}

func (r *pgCategoriaRepository) Update(ctx context.Context, categoria *models.Categoria) error {
	sql := `UPDATE categorias SET nome = $1, icone = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, sql, categoria.Nome, categoria.Icone, categoria.ID)
	return err
}

func (r *pgCategoriaRepository) Delete(ctx context.Context, id string) error {
	sql := `DELETE FROM categorias WHERE id = $1`
	_, err := r.db.Exec(ctx, sql, id)
	return err
}