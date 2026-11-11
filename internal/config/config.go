// Package config centralizará a leitura e a validação das variáveis de ambiente.
package config

// Config define as configurações previstas para a execução.
type Config struct {
	BaseURL   string
	Token     string
	TimeoutMS float64
}
