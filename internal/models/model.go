package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type TipoAtivo string

const (
	AtivoContaCorrente TipoAtivo = "CONTA_CORRENTE"
	AtivoCartaoCredito TipoAtivo = "CARTÃO_CREDITO"
)

type TipoTransacao string

const (
	TransacaoRecebimento TipoTransacao = "RECEBIMENTO"
	TransacaoDebito      TipoTransacao = "DEBITO"
	TransacaoCredito     TipoTransacao = "CREDITO"
	TransacaoEstorno     TipoTransacao = "ESTORNO"
)

type Categoria struct {
	ID    string `json:"id" db:"id"`
	Nome  string `json:"nome" db:"nome"`
	Icone string `json:"icone" db:"icone"`
}

type AtivoFinanceiro struct {
	ID               string    `json:"id" db:"id"`
	Instituicao      string    `json:"instituicao" db:"instituicao"`
	Nome             string    `json:"nome" db:"nome"`
	Tipo             TipoAtivo `json:"tipo" db:"tipo"`
	SaldoAtual       float64   `json:"saldo_atual" db:"saldo_atual"`
	LimiteDisponivel float64   `json:"limite_disponivel" db:"limite_disponivel"`
	IsActive         bool      `json:"is_active" db:"is_active"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type Transacao struct {
	ID                string        `json:"id" db:"id"`
	AtivoFinanceiroID string        `json:"ativo_financeiro_id" db:"ativo_financeiro_id"`
	CategoriaID       string        `json:"categoria_id" db:"categoria_id"`
	Descricao         string        `json:"descricao" db:"descricao"`
	Valor             float64       `json:"valor" db:"valor"`
	Tipo              TipoTransacao `json:"tipo" db:"tipo"`
	ReversalOf        *string       `json:"reversal_of,omitempty" db:"reversal_of"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
}

type TransacaoRecorrente struct {
	ID                string        `json:"id" db:"id"`
	AtivoFinanceiroID string        `json:"ativo_financeiro_id" db:"ativo_financeiro_id"`
	CategoriaID       string        `json:"categoria_id" db:"categoria_id"`
	Descricao         string        `json:"descricao" db:"descricao"`
	Valor             float64       `json:"valor" db:"valor"`
	Tipo              TipoTransacao `json:"tipo" db:"tipo"`
	DiaDoVencimento   int           `json:"dia_do_vencimento" db:"dia_do_vencimento"`
	Ativa             bool          `json:"ativa" db:"ativa"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at" db:"updated_at"`
}

func (t *TipoAtivo) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch TipoAtivo(s) {
	case AtivoContaCorrente, AtivoCartaoCredito:
		*t = TipoAtivo(s)
		return nil
	default:
		return fmt.Errorf("tipo de ativo inválido: %s", s)
	}
}

func (t *TipoTransacao) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch TipoTransacao(s) {
	case TransacaoRecebimento, TransacaoDebito, TransacaoCredito, TransacaoEstorno:
		*t = TipoTransacao(s)
		return nil
	default:
		return fmt.Errorf("tipo de transação inválido: %s", s)
	}
}