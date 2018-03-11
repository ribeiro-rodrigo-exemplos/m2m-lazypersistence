package dispatch

import (
	"fmt"
	"m2m-lazypersistence/internal/pkg/mensageria"
	"m2m-lazypersistence/internal/pkg/repo"
	"strconv"

	"gopkg.in/mgo.v2"
)

var methods = map[string]method{
	"insert": insert,
}

// Dispatcher - Responsável pelo despacho das mensagens
type Dispatcher struct {
	Host     string
	Port     int
	Database string
	session  *mgo.Session
}

type method func(*Dispatcher, repo.Operation) error

// Dispatch - Salva as mensagens no mongodb
func (d *Dispatcher) Dispatch(repository repo.Repository) {
	fmt.Println("Gravando dados no mongo")

	if d.session == nil {
		err := d.openSession()

		if err != nil {
			fmt.Println("Erro ao abrir sessão com o mongodb")
			repository.Reject()
			return
		}
	}

	go func() {

		repository.Each(func(key string, operation repo.Operation) {
			method := methods[operation.Action]
			err := method(d, operation)

			if err != nil {
				operation.Reject()
			} else {
				operation.Confirm()
			}
		})
	}()
}

func (d *Dispatcher) openSession() error {
	session, err := mgo.Dial(d.Host + ":" + strconv.Itoa(d.Port))
	d.session = session

	return err
}

func insert(dispatcher *Dispatcher, operation repo.Operation) error {

	var payloads []interface{}

	operation.Messages.Each(func(message mensageria.Message) {
		payloads = append(payloads, message.Payload)
	})

	collection := dispatcher.session.DB(dispatcher.Database).C(operation.Collection)
	err := collection.Insert(payloads...)
	return err
}
