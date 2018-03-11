package app

import (
	"fmt"
	cfg "m2m-lazypersistence/internal/pkg/config"
	"m2m-lazypersistence/internal/pkg/dispatch"
	"m2m-lazypersistence/internal/pkg/mensageria"
	"m2m-lazypersistence/internal/pkg/repo"
	"time"
)

var repository repo.Repository
var channelMessage = make(chan mensageria.Message)
var signalDispatcher = make(chan struct{})

var maxMessages int
var retentionTime int

var dispatcher dispatch.Dispatcher

// Bootstrap - Função Inicial do projeto
func Bootstrap(config cfg.Config) {

	maxMessages = config.MaxMessages
	retentionTime = config.RetentionSeconds

	consumer := mensageria.Consumer{
		Host:     config.RabbitMQ.Host,
		Port:     config.RabbitMQ.Port,
		User:     config.RabbitMQ.User,
		Password: config.RabbitMQ.Password,
	}

	dispatcher = dispatch.Dispatcher{
		Host:     config.MongoDB.Host,
		Port:     config.MongoDB.Port,
		Database: config.MongoDB.Database,
	}

	consumer.Connect(func(message mensageria.Message) {
		fmt.Println("------", message.Payload)
		channelMessage <- message
	})

	go eventRouter()
	go dispatcherListener()
}

func eventRouter() {
	for {
		select {
		case <-signalDispatcher:
			dispatchMassages()
		case message := <-channelMessage:
			repository.Save(message)
			evaluateDispatch()
		}
	}
}

func dispatcherListener() {
	for {
		time.Sleep(time.Second * time.Duration(retentionTime))
		signalDispatcher <- struct{}{}
	}
}

func dispatchMassages() {
	cloneRepository := repository.Clone()
	dispatcher.Dispatch(cloneRepository)
	repository.Clear()
}

func evaluateDispatch() {
	if repository.Size() >= maxMessages {
		dispatchMassages()
	}
}
