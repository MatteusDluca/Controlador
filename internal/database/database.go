package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

)

var DB *pgxpool.Pool

func Connect() {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	var err error
	for i := 0; i < 15; i++ {
		DB, err = pgxpool.New(context.Background(), connString)
		if err == nil && DB.Ping(context.Background()) == nil {
			log.Info().Msg("Conexão com o banco de dados estabelecida com sucesso.")
			return
		}
		time.Sleep(2 * time.Second)
	}
	log.Fatal().Err(err).Msg("Não foi possível conectar ao banco de dados.")
}

func Migrate() {
	// Migração de Ativos
	createAtivosSQL := `
	CREATE TABLE IF NOT EXISTS ativos_financeiros (
		id UUID PRIMARY KEY,
		usuario_id UUID,
		nome VARCHAR(255) NOT NULL,
		tipo VARCHAR(50) NOT NULL,
		saldo_atual NUMERIC(15, 2) DEFAULT 0.00,
		limite_disponivel NUMERIC(15, 2) DEFAULT 0.00,
		is_active BOOLEAN NOT NULL DEFAULT TRUE, -- NOVO CAMPO
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`
	if _, err := DB.Exec(context.Background(), createAtivosSQL); err != nil {
		log.Fatal().Err(err).Msg("Falha ao migrar tabela 'ativos_financeiros'.")
	}
	log.Info().Msg("Migração da tabela 'ativos_financeiros' concluída.")

	// Migração de Transações
	createTransacoesSQL := `
	CREATE TABLE IF NOT EXISTS transacoes (
		id UUID PRIMARY KEY,
		ativo_financeiro_id UUID NOT NULL REFERENCES ativos_financeiros(id) ON DELETE CASCADE,
		descricao VARCHAR(255) NOT NULL,
		valor NUMERIC(15, 2) NOT NULL,
		tipo VARCHAR(50) NOT NULL,
		reversal_of UUID NULL REFERENCES transacoes(id), -- NOVO CAMPO
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`
	if _, err := DB.Exec(context.Background(), createTransacoesSQL); err != nil {
		log.Fatal().Err(err).Msg("Falha ao migrar tabela 'transacoes'.")
	}
	log.Info().Msg("Migração da tabela 'transacoes' concluída.")
}
