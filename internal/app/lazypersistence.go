package app

import (
	"fmt"
	cfg "m2m-lazypersistence/internal/pkg/config"
	"m2m-lazypersistence/internal/pkg/mensageria"
	"time"
)

var repositorio = make(map[string][]mensageria.Message)

const maxMessages = 10

// Bootstrap - Função Inicial do projeto
func Bootstrap(config cfg.Config) {

	consumer := mensageria.Consumer{
		Host:     config.RabbitMQ.Host,
		Port:     config.RabbitMQ.Port,
		User:     config.RabbitMQ.User,
		Password: config.RabbitMQ.Password,
	}

	//defer consumer.Disconnect()

	var canal = make(chan mensageria.Message)
	var signal = make(chan struct{})

	consumer.Connect(func(message mensageria.Message) {
		canal <- message
	})

	go func() {

		for {

			select {
			case <-signal:
				dispatcherMessages(repositorio)
			case message := <-canal:
				saveMessage(message)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * 3)
			signal <- struct{}{}
		}
	}()
}

func saveMessage(message mensageria.Message) {
	fmt.Println(message.Payload)
	messagesPendings := repositorio[message.Headers.Action]
	messagesPendings = append(messagesPendings, message)
	repositorio[message.Headers.Action] = messagesPendings
	if len(repositorio) >= maxMessages {
		dispatcherMessages(copyRepository())
	}
}

func dispatcherMessages(repository map[string][]mensageria.Message) {
	fmt.Println("signal - gravando dados no mongo")
}

func copyRepository() map[string][]mensageria.Message {
	newRepository := make(map[string][]mensageria.Message)

	for chave, valor := range repositorio {
		newRepository[chave] = valor
	}

	return newRepository
}

func logarRepositorio() {
	for chave, valor := range repositorio {
		fmt.Println(chave + "**********************")
		for _, message := range valor {
			fmt.Println(message.Payload)
		}
	}
}
