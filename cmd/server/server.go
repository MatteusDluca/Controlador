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

	// Repositórios
	ativoRepo := repositories.NewPgAtivoRepository(database.DB)
	transacaoRepo := repositories.NewPgTransacaoRepository(database.DB)
	categoriaRepo := repositories.NewPgCategoriaRepository(database.DB)
	transacaoRecorrenteRepo := repositories.NewPgTransacaoRecorrenteRepository(database.DB)

	// Serviços
	createAtivoSvc := services.NewCreateAtivoService(ativoRepo)
	listAtivoSvc := services.NewListAtivosService(ativoRepo)
	deactivateAtivoSvc := services.NewDeactivateAtivoService(ativoRepo)
	createTransacaoSvc := services.NewCreateTransacaoService(database.DB, transacaoRepo, ativoRepo, categoriaRepo)
	listTransacoesSvc := services.NewListTransacoesService(transacaoRepo)
	reverseTransacaoSvc := services.NewReverseTransacaoService(database.DB, transacaoRepo, ativoRepo)
	createCategoriaSvc := services.NewCreateCategoriaService(categoriaRepo)
	listCategoriaSvc := services.NewListCategoriasService(categoriaRepo)
	
	createRecorrenciaSvc := services.NewCreateTransacaoRecorrenteService(transacaoRecorrenteRepo, ativoRepo, categoriaRepo)
	// ALTERAÇÃO: Corrigido para instanciar o serviço a partir do pacote 'services'.
	listRecorrenciasSvc := services.NewListTransacoesRecorrentesService(transacaoRecorrenteRepo)
	processarRecorrenciasSvc := services.NewProcessarRecorrenciasService(transacaoRecorrenteRepo, createTransacaoSvc)

	// Handlers
	ativoHandler := handlers.NewAtivoHandler(createAtivoSvc, listAtivoSvc, deactivateAtivoSvc)
	transacaoHandler := handlers.NewTransacaoHandler(createTransacaoSvc, listTransacoesSvc, reverseTransacaoSvc)
	categoriaHandler := handlers.NewCategoriaHandler(createCategoriaSvc, listCategoriaSvc)
	transacaoRecorrenteHandler := handlers.NewTransacaoRecorrenteHandler(createRecorrenciaSvc, listRecorrenciasSvc, processarRecorrenciasSvc)


	// --- SETUP DO SERVIDOR ---
	r := router.SetupRouter(ativoHandler, transacaoHandler, categoriaHandler, transacaoRecorrenteHandler)

	log.Info().Msg("Servidor iniciado na porta :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Falha ao iniciar o servidor Gin")
	}
}