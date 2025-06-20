package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"controlador/backend/internal/models"
	"controlador/backend/internal/services"

)

type TransacaoRecorrenteHandler struct {
	createService    *services.CreateTransacaoRecorrenteService
	listService      *services.ListTransacoesRecorrentesService // Corrigido para usar o serviço importado
	processarService *services.ProcessarRecorrenciasService
}

func NewTransacaoRecorrenteHandler(createSvc *services.CreateTransacaoRecorrenteService, listSvc *services.ListTransacoesRecorrentesService, processarSvc *services.ProcessarRecorrenciasService) *TransacaoRecorrenteHandler {
	return &TransacaoRecorrenteHandler{
		createService:    createSvc,
		listService:      listSvc,
		processarService: processarSvc,
	}
}

func (h *TransacaoRecorrenteHandler) CreateTransacaoRecorrente(c *gin.Context) {
	var input models.TransacaoRecorrente
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Msg("Erro no bind do JSON para criar transação recorrente")
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição inválido"})
		return
	}

	novaRecorrencia, err := h.createService.Execute(c.Request.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Erro no serviço de criação de transação recorrente")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, novaRecorrencia)
}

func (h *TransacaoRecorrenteHandler) ListTransacoesRecorrentesPorAtivo(c *gin.Context) {
	ativoID := c.Param("id")
	if ativoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do ativo é obrigatório"})
		return
	}

	recorrencias, err := h.listService.Execute(c.Request.Context(), ativoID)
	if err != nil {
		log.Error().Err(err).Msg("Erro ao listar transações recorrentes")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao listar transações recorrentes"})
		return
	}

	c.JSON(http.StatusOK, recorrencias)
}

func (h *TransacaoRecorrenteHandler) ProcessarRecorrencias(c *gin.Context) {
	log.Info().Msg("Requisição para acionar o worker de processamento de recorrências recebida.")
	relatorio, err := h.processarService.Execute(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("Erro na execução do worker de processamento de recorrências")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao processar recorrências"})
		return
	}

	c.JSON(http.StatusOK, relatorio)
}