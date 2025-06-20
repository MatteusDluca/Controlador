package router

import (
	"controlador/backend/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func SetupRouter(
	ativoHandler *handlers.AtivoHandler,
	transacaoHandler *handlers.TransacaoHandler,
	categoriaHandler *handlers.CategoriaHandler,
	transacaoRecorrenteHandler *handlers.TransacaoRecorrenteHandler,
) *gin.Engine {
	router := gin.New()
	router.Use(ginZerologLogger())
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		log.Debug().Msg("Recebida requisição na rota /ping")
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	apiV1 := router.Group("/api/v1")
	{
		// Rotas de Ativos
		apiV1.POST("/ativos", ativoHandler.CreateAtivoFinanceiro)
		apiV1.GET("/ativos", ativoHandler.GetAtivosFinanceiros)
		apiV1.DELETE("/ativos/:id", ativoHandler.DeactivateAtivoFinanceiro)

		// Rotas de Transações
		apiV1.POST("/transacoes", transacaoHandler.CreateTransacao)
		apiV1.GET("/transacoes", transacaoHandler.GetTransacoes)
		// ALTERAÇÃO: Nova rota para estornar uma transação.
		apiV1.POST("/transacoes/:id/reverter", transacaoHandler.ReverseTransacao)

		// Rotas de Categorias
		apiV1.POST("/categorias", categoriaHandler.CreateCategoria)
		apiV1.GET("/categorias", categoriaHandler.GetCategorias)

		// Rotas de Transações Recorrentes
		apiV1.POST("/recorrencias", transacaoRecorrenteHandler.CreateTransacaoRecorrente)
		// CORREÇÃO: Esta rota estava causando o 404 e agora está corretamente registrada.
		apiV1.GET("/ativos/:id/recorrencias", transacaoRecorrenteHandler.ListTransacoesRecorrentesPorAtivo)
	}

	admin := router.Group("/admin")
	{
		admin.POST("/workers/processar-recorrencias", transacaoRecorrenteHandler.ProcessarRecorrencias)
	}

	return router
}

func ginZerologLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		var logEvent *zerolog.Event
		if c.Writer.Status() >= http.StatusInternalServerError {
			logEvent = log.Error()
		} else if c.Writer.Status() >= http.StatusBadRequest {
			logEvent = log.Warn()
		} else {
			logEvent = log.Info()
		}

		logEvent.
			Str("method", c.Request.Method).
			Int("status", c.Writer.Status()).
			Str("path", path).
			Str("query", query).
			Str("ip", c.ClientIP()).
			Dur("latency", latency).
			Str("user_agent", c.Request.UserAgent()).
			Msg("Requisição HTTP Recebida")
	}
}