package app

import (
	"fmt"
	cfg "m2m-lazypersistence/internal/pkg/config"
	"m2m-lazypersistence/internal/pkg/mensageria"
)

// Bootstrap - Função Inicial do projeto
func Bootstrap(config cfg.Config) {

	consumer := mensageria.Consumer{
		Host:     config.RabbitMQ.Host,
		Port:     config.RabbitMQ.Port,
		User:     config.RabbitMQ.User,
		Password: config.RabbitMQ.Password,
	}

	defer consumer.Disconnect()

	consumer.Connect(func(message mensageria.Message) {
		fmt.Println(message.Payload)
		//mensagem.Confirmar()
	})
}
