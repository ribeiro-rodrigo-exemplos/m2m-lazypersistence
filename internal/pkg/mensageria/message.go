package mensageria

import (
	"github.com/streadway/amqp"
)

// Headers - Cabeçalhos da mensagem
type Headers struct {
	Action     string `json:"action"`
	Collection string `json:"collection"`
}

// Message - mensagem consumida do rabbitmq
type Message struct {
	Headers  Headers     `json:"headers"`
	Payload  interface{} `json:"payload"`
	delivery amqp.Delivery
}

// Confirm - confirma processamento da mensagem para o rabbitmq
func (m *Message) Confirm() {
	m.delivery.Ack(false)
}

// Reject - rejeita processamento da mensagem
func (m *Message) Reject() {
	m.delivery.Nack(false, true)
}
