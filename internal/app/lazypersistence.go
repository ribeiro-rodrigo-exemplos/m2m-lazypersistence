package app

import (
	"log"
	cfg "m2m-lazypersistence/internal/pkg/config"
	"m2m-lazypersistence/internal/pkg/dispatch"
	"m2m-lazypersistence/internal/pkg/mensageria"
	"m2m-lazypersistence/internal/pkg/repo"
	"time"
)

var repository repo.Repository
var channelRequest = make(chan mensageria.RequestPersistence)
var signalDispatcher = make(chan struct{})

var maxMessages int
var retentionTime int

var dispatcher dispatch.Dispatcher

// Bootstrap - Função Inicial do projeto
func Bootstrap(config cfg.Config) {

	maxMessages = config.MaxMessages
	retentionTime = config.RetentionSeconds

	createConsumers(config)

	dispatcher = dispatch.Dispatcher{
		Host:     config.MongoDB.Host,
		Port:     config.MongoDB.Port,
		Database: config.MongoDB.Database,
	}

	go eventRouter()
	go dispatcherListener()
}

func eventRouter() {
	for {
		select {
		case <-signalDispatcher:
			dispatchMassages()
		case request := <-channelRequest:
			repository.Save(request)
			evaluateDispatch()
		}
	}
}

func dispatcherListener() {
	for {
		time.Sleep(time.Second * time.Duration(retentionTime))

		log.Println("Tempo de retenção atingido")

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

func createConsumers(config cfg.Config) {

	for _, queueNames := range config.RabbitMQ.Queues {
		consumer := mensageria.Consumer{
			Host:     config.RabbitMQ.Host,
			Port:     config.RabbitMQ.Port,
			User:     config.RabbitMQ.User,
			Password: config.RabbitMQ.Password,
			Queue:    queueNames,
		}

		consumer.Connect(listener)
	}
}

func listener(request mensageria.RequestPersistence) {
	log.Println("Mensagem recebida:", request.Message.Payload)
	channelRequest <- request
}
