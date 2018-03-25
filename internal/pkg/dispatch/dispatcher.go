package dispatch

import (
	"log"
	"m2m-lazypersistence/internal/pkg/repo"
	"strconv"

	"gopkg.in/mgo.v2"
)

// Dispatcher - Responsável pelo despacho das mensagens
type Dispatcher struct {
	Host     string
	Port     int
	Database string
	session  *mgo.Session
}

// Dispatch - Salva as mensagens no mongodb
func (d *Dispatcher) Dispatch(repository repo.Repository) {
	log.Println("Despejando repositório no mongodb")

	if d.session == nil {
		err := d.openSession()

		if err != nil {
			log.Println("Erro ao abrir sessão com o mongodb")
			repository.Reject()
			return
		}
	}

	go d.execute(repository)
}

func (d *Dispatcher) execute(repository repo.Repository) {
	repository.Each(func(key string, operation repo.Operation) {

		err := executeAction(d, operation)

		if err != nil && err != mgo.ErrNotFound {
			operation.Reject()
		} else {
			operation.Confirm()
		}
	})
}

func (d *Dispatcher) openSession() error {
	session, err := mgo.Dial(d.Host + ":" + strconv.Itoa(d.Port))
	d.session = session

	return err
}
