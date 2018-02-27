package manager

import (
	"fmt"
	"m2m-lazypersistence/internal/pkg/mensageria"
)

func Init() {

	consumer := mensageria.Consumer{
		Host:     "localhost",
		Port:     5672,
		User:     "guest",
		Password: "guest",
	}

	defer consumer.Disconnect()

	consumer.Connect(func(message mensageria.Message) {
		fmt.Println(message.Payload)
		//mensagem.Confirmar()
	})
}
