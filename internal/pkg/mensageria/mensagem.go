package mensageria

import "github.com/streadway/amqp"

// Mensagem - mensagem consumida do rabbitmq
type Mensagem struct {
	delivery amqp.Delivery
}

// Confirmar - confirma processamento da mensagem para o rabbitmq
func (m *Mensagem) Confirmar() {
	m.delivery.Ack(false)
}

// Get - obtem mensagem consumida
func (m *Mensagem) Get() string {
	return string(m.delivery.Body)
}
