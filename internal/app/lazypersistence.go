package app

import (
	"fmt"
	cfg "m2m-lazypersistence/internal/pkg/config"
	"m2m-lazypersistence/internal/pkg/dispatcher"
	"m2m-lazypersistence/internal/pkg/mensageria"
	"time"
)

var repository = make(map[string][]mensageria.Message)
var channelMessage = make(chan mensageria.Message)
var signalDispatcher = make(chan struct{})

var maxMessages int

// Bootstrap - Função Inicial do projeto
func Bootstrap(config cfg.Config) {

	maxMessages = config.MaxMessages

	consumer := mensageria.Consumer{
		Host:     config.RabbitMQ.Host,
		Port:     config.RabbitMQ.Port,
		User:     config.RabbitMQ.User,
		Password: config.RabbitMQ.Password,
	}

	consumer.Connect(func(message mensageria.Message) {
		channelMessage <- message
	})

	go eventRouter()
	go dispatcherListener()
}

func eventRouter() {
	for {

		select {
		case <-signalDispatcher:
			dispatcher.Dispatch(repository)
		case message := <-channelMessage:
			saveMessage(message)
			evaluateDispatch()
		}
	}
}

func dispatcherListener() {
	for {
		time.Sleep(time.Second * 30)
		signalDispatcher <- struct{}{}
	}
}

func saveMessage(message mensageria.Message) {
	fmt.Println("------", message.Payload)
	messagesPendings := repository[message.Headers.Action]
	messagesPendings = append(messagesPendings, message)
	repository[message.Headers.Action] = messagesPendings
}

func evaluateDispatch() {
	if len(repository) >= maxMessages {
		dispatcher.Dispatch(repository)
	}
}

func logarRepositorio() {
	for chave, valor := range repository {
		fmt.Println(chave + "**********************")
		for _, message := range valor {
			fmt.Println(message.Payload)
		}
	}
}
