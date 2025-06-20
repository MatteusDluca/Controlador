package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"controlador/backend/internal/models"
	"controlador/backend/internal/services"

)

type AtivoHandler struct {
	createService     *services.CreateAtivoService
	listService       *services.ListAtivosService
	deactivateService *services.DeactivateAtivoService
}

// CORREÇÃO: Adicionado o deactivateSvc como parâmetro.
func NewAtivoHandler(createSvc *services.CreateAtivoService, listSvc *services.ListAtivosService, deactivateSvc *services.DeactivateAtivoService) *AtivoHandler {
	return &AtivoHandler{
		createService:     createSvc,
		listService:       listSvc,
		deactivateService: deactivateSvc,
	}
}

func (h *AtivoHandler) DeactivateAtivo(c *gin.Context) {
	id := c.Param("id")
	err := h.deactivateService.Execute(c.Request.Context(), id)
	if err != nil {
		if err == services.ErrAtivoNaoEncontrado {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Error().Err(err).Msg("Erro ao desativar ativo")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao desativar ativo"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AtivoHandler) CreateAtivoFinanceiro(c *gin.Context) {
	var input models.AtivoFinanceiro
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Msg("Erro ao fazer bind do JSON para criar ativo")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novoAtivo, err := h.createService.Execute(c.Request.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Erro ao criar ativo financeiro")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar o ativo financeiro"})
		return
	}

	c.JSON(http.StatusCreated, novoAtivo)
}

func (h *AtivoHandler) GetAtivosFinanceiros(c *gin.Context) {
	ativos, err := h.listService.Execute(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("Erro ao buscar ativos financeiros")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar ativos financeiros"})
		return
	}

	c.JSON(http.StatusOK, ativos)
}