package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"controlador/backend/internal/database"
	"controlador/backend/internal/handlers"
	"controlador/backend/internal/repositories"
	"controlador/backend/internal/router"
	"controlador/backend/internal/services"

)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msg("Iniciando o Controlador...")

	database.Connect()
	database.Migrate()

	// --- INJEÇÃO DE DEPENDÊNCIAS ---
	ativoRepo := repositories.NewPgAtivoRepository(database.DB)
	transacaoRepo := repositories.NewPgTransacaoRepository(database.DB)

	createAtivoSvc := services.NewCreateAtivoService(ativoRepo)
	listAtivoSvc := services.NewListAtivosService(ativoRepo)
	deactivateAtivoSvc := services.NewDeactivateAtivoService(ativoRepo)

	// CORREÇÃO: Passa o database.DB para o construtor do serviço
	createTransacaoSvc := services.NewCreateTransacaoService(database.DB, transacaoRepo, ativoRepo)

	listTransacoesSvc := services.NewListTransacoesService(transacaoRepo)
	reverseTransacaoSvc := services.NewReverseTransacaoService(database.DB, transacaoRepo, ativoRepo)

	ativoHandler := handlers.NewAtivoHandler(createAtivoSvc, listAtivoSvc, deactivateAtivoSvc)
	transacaoHandler := handlers.NewTransacaoHandler(createTransacaoSvc, listTransacoesSvc, reverseTransacaoSvc)

	// --- SETUP DO SERVIDOR ---
	r := router.SetupRouter(ativoHandler, transacaoHandler)

	log.Info().Msg("Servidor iniciado na porta :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Falha ao iniciar o servidor Gin")
	}
}