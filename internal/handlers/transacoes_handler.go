package handlers

import (
	"controlador/backend/internal/models"
	"controlador/backend/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type TransacaoHandler struct {
	createService  *services.CreateTransacaoService
	listService    *services.ListTransacoesService
	reverseService *services.ReverseTransacaoService
}

func NewTransacaoHandler(createSvc *services.CreateTransacaoService, listSvc *services.ListTransacoesService, reverseSvc *services.ReverseTransacaoService) *TransacaoHandler {
	return &TransacaoHandler{
		createService:  createSvc,
		listService:    listSvc,
		reverseService: reverseSvc,
	}
}

// ALTERAÇÃO: Este método foi adicionado para lidar com a rota de estorno.
func (h *TransacaoHandler) ReverseTransacao(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da transação é obrigatório"})
		return
	}

	estorno, err := h.reverseService.Execute(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrTransacaoJaEstornada) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Error().Err(err).Msg("Erro ao estornar transação")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao estornar transação"})
		return
	}
	c.JSON(http.StatusCreated, estorno)
}

func (h *TransacaoHandler) CreateTransacao(c *gin.Context) {
	var input models.Transacao
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Msg("Erro no bind do JSON para criar transação")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novaTransacao, err := h.createService.Execute(c.Request.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Erro no serviço de criação de transação")
		if errors.Is(err, services.ErrSaldoInsuficiente) || errors.Is(err, services.ErrAtivoNaoEncontrado) || errors.Is(err, services.ErrTipoTransacaoInvalido) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar a transação"})
		return
	}

	c.JSON(http.StatusCreated, novaTransacao)
}

func (h *TransacaoHandler) GetTransacoes(c *gin.Context) {
	transacoes, err := h.listService.Execute(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("Erro ao buscar transações")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar transações"})
		return
	}
	c.JSON(http.StatusOK, transacoes)
}