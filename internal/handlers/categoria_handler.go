package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"controlador/backend/internal/models"
	"controlador/backend/internal/services"

)

type CategoriaHandler struct {
	createService *services.CreateCategoriaService
	listService   *services.ListCategoriasService
	// Futuramente, podemos adicionar aqui os serviços de update e delete.
}

func NewCategoriaHandler(createSvc *services.CreateCategoriaService, listSvc *services.ListCategoriasService) *CategoriaHandler {
	return &CategoriaHandler{
		createService: createSvc,
		listService:   listSvc,
	}
}

func (h *CategoriaHandler) CreateCategoria(c *gin.Context) {
	var input models.Categoria
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Msg("Erro no bind do JSON para criar categoria")
		c.JSON(http.StatusBadRequest, gin.H{"error": "corpo da requisição inválido"})
		return
	}

	novaCategoria, err := h.createService.Execute(c.Request.Context(), input)
	if err != nil {
		log.Error().Err(err).Msg("Erro no serviço de criação de categoria")
		if errors.Is(err, services.ErrCategoriaJaExiste) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar categoria"})
		return
	}

	c.JSON(http.StatusCreated, novaCategoria)
}

func (h *CategoriaHandler) GetCategorias(c *gin.Context) {
	categorias, err := h.listService.Execute(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("Erro ao buscar categorias")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar categorias"})
		return
	}
	c.JSON(http.StatusOK, categorias)
}