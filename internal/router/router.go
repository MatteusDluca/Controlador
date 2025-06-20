package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"controlador/backend/internal/handlers"

)

// SetupRouter configura o motor Gin com todos os middlewares e rotas.
// Ele recebe os handlers como dependências para conectar as rotas à lógica correta.
func SetupRouter(
	ativoHandler *handlers.AtivoHandler,
	transacaoHandler *handlers.TransacaoHandler,
) *gin.Engine {
	// Cria um novo motor Gin, sem os middlewares padrão.
	router := gin.New()

	// Adiciona nossos middlewares customizados de forma explícita.
	router.Use(ginZerologLogger()) // Nosso logger estruturado.
	router.Use(gin.Recovery())   // Recupera de panics para não derrubar o servidor.

	// Rota de verificação de saúde do sistema.
	router.GET("/ping", func(c *gin.Context) {
		log.Debug().Msg("Recebida requisição na rota /ping")
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Agrupa todas as rotas da nossa API sob o prefixo /api/v1 para versionamento.
	apiV1 := router.Group("/api/v1")
	{
		// Rotas para Ativos Financeiros
		apiV1.POST("/ativos", ativoHandler.CreateAtivoFinanceiro)
		apiV1.GET("/ativos", ativoHandler.GetAtivosFinanceiros)

		// Rotas para Transações
		apiV1.POST("/transacoes", transacaoHandler.CreateTransacao)
		apiV1.GET("/transacoes", transacaoHandler.GetTransacoes)
	}

	// Retorna o motor Gin totalmente configurado.
	return router
}

// ginZerologLogger é um middleware para o Gin que usa o Zerolog para logar cada requisição.
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