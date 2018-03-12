package dispatch

import (
	"m2m-lazypersistence/internal/pkg/mensageria"
	"m2m-lazypersistence/internal/pkg/repo"
	"strings"
)

type action func(*Dispatcher, repo.Operation) error

var actions = map[string]action{
	"INSERT": insert,
	"PUSH":   push,
	"PULL":   pull,
}

func executeAction(dispatcher *Dispatcher, operation repo.Operation) error {

	action := actions[strings.ToUpper(operation.Action)]
	err := action(dispatcher, operation)

	return err
}

func extractPayload(operation repo.Operation) (payloads []interface{}) {

	operation.Messages.Each(func(message mensageria.Message) {
		payloads = append(payloads, message.Payload)
	})

	return
}

func insert(dispatcher *Dispatcher, operation repo.Operation) error {

	payloads := extractPayload(operation)

	collection := dispatcher.session.DB(dispatcher.Database).C(operation.Collection)
	err := collection.Insert(payloads...)
	return err
}

func push(dispatcher *Dispatcher, operation repo.Operation) error {
	return nil
}

func pull(dispatcher *Dispatcher, operation repo.Operation) error {

	return nil
}
