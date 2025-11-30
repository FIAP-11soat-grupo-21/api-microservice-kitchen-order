package interfaces

import (
	"context"
)

// Message representa uma mensagem genérica do broker
type Message struct {
	ID      string
	Body    []byte
	Headers map[string]string
}

// MessageHandler é a função que processa uma mensagem
type MessageHandler func(ctx context.Context, message Message) error

// MessageBroker interface genérica para diferentes implementações de message broker
type MessageBroker interface {
	// Connect conecta ao broker
	Connect(ctx context.Context) error
	
	// Close fecha a conexão com o broker
	Close() error
	
	// Publish publica uma mensagem em uma fila/tópico
	Publish(ctx context.Context, queue string, message Message) error
	
	// Subscribe inscreve-se em uma fila/tópico para receber mensagens
	Subscribe(ctx context.Context, queue string, handler MessageHandler) error
	
	// Start inicia o consumo de mensagens
	Start(ctx context.Context) error
	
	// Stop para o consumo de mensagens
	Stop() error
}

