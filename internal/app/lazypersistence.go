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
var channelRequest chan mensageria.RequestPersistence
var signalDispatcher chan struct{}

var maxMessages int
var retentionTime int

var dispatcher dispatch.Dispatcher

// Bootstrap - Função Inicial do projeto
func Bootstrap(config cfg.Config) {

	maxMessages = config.MaxMessages
	retentionTime = config.RetentionSeconds

	channelRequest = make(chan mensageria.RequestPersistence, config.MaxMessages)
	signalDispatcher = make(chan struct{})

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

	if repository.Size() == 0 {
		return
	}

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

	consumer := mensageria.Consumer{
		Host:          config.RabbitMQ.Host,
		Port:          config.RabbitMQ.Port,
		User:          config.RabbitMQ.User,
		PrefetchCount: config.MaxMessages * 2,
		Password:      config.RabbitMQ.Password,
		Queues:        config.RabbitMQ.Queues,
	}

	consumer.Connect(listener)
}

func listener(request mensageria.RequestPersistence) {
	log.Println("Mensagem recebida:", request.Message.Payload)
	channelRequest <- request
}
