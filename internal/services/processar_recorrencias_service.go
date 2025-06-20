package services

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"controlador/backend/internal/models"
	"controlador/backend/internal/repositories"

)

type ProcessarRecorrenciasService struct {
	trRepo             repositories.TransacaoRecorrenteRepository
	createTransacaoSvc *CreateTransacaoService
}

func NewProcessarRecorrenciasService(trr repositories.TransacaoRecorrenteRepository, cts *CreateTransacaoService) *ProcessarRecorrenciasService {
	return &ProcessarRecorrenciasService{
		trRepo:             trr,
		createTransacaoSvc: cts,
	}
}

// RelatorioProcessamento contém o resultado da execução do worker.
type RelatorioProcessamento struct {
	TotalParaProcessar int
	Sucesso            int
	Falhas             int
	Erros              []string
}

func (s *ProcessarRecorrenciasService) Execute(ctx context.Context) (*RelatorioProcessamento, error) {
	diaAtual := time.Now().Day()
	log.Info().Int("dia", diaAtual).Msg("Iniciando processamento de transações recorrentes.")

	recorrencias, err := s.trRepo.FindActiveByDay(ctx, diaAtual)
	if err != nil {
		log.Error().Err(err).Msg("Erro ao buscar transações recorrentes do dia.")
		return nil, err
	}

	relatorio := &RelatorioProcessamento{
		TotalParaProcessar: len(recorrencias),
	}

	if relatorio.TotalParaProcessar == 0 {
		log.Info().Msg("Nenhuma transação recorrente para processar hoje.")
		return relatorio, nil
	}

	log.Info().Int("total", relatorio.TotalParaProcessar).Msg("Transações encontradas para processamento.")

	for _, recorrencia := range recorrencias {
		transacao := models.Transacao{
			AtivoFinanceiroID: recorrencia.AtivoFinanceiroID,
			CategoriaID:       recorrencia.CategoriaID,
			Descricao:         fmt.Sprintf("Recorrência: %s", recorrencia.Descricao),
			Valor:             recorrencia.Valor,
			Tipo:              recorrencia.Tipo,
		}

		_, err := s.createTransacaoSvc.Execute(ctx, transacao)
		if err != nil {
			log.Error().Err(err).Str("recorrencia_id", recorrencia.ID).Msg("Falha ao processar transação recorrente.")
			relatorio.Falhas++
			relatorio.Erros = append(relatorio.Erros, err.Error())
		} else {
			log.Info().Str("recorrencia_id", recorrencia.ID).Msg("Transação recorrente processada com sucesso.")
			relatorio.Sucesso++
		}
	}

	log.Info().Interface("relatorio", relatorio).Msg("Processamento de transações recorrentes concluído.")
	return relatorio, nil
}