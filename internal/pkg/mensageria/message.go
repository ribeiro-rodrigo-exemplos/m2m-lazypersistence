package mensageria

import (
	"github.com/streadway/amqp"
)

// Message - mensagem consumida do rabbitmq
type Message struct {
	Payload  interface{}
	delivery amqp.Delivery
}

// RequestPersistence - estrutura que representa uma requisição de persistencia
type RequestPersistence struct {
	Message Message
	Headers map[string]interface{}
}

// Confirm - confirma processamento da mensagem para o rabbitmq
func (m *Message) Confirm() {
	m.delivery.Ack(false)
}

// Reject - rejeita processamento da mensagem, fazendo requeue
func (m *Message) Reject() {
	m.delivery.Nack(false, true)
}

/*Discard - descarta mensagem, enviando a mesma para ao dead letter exchange
caso o mesmo esteja configurado*/
func (m *Message) Discard() {
	m.delivery.Nack(false, false)
}
