package repo

import (
	"fmt"
	"m2m-lazypersistence/internal/pkg/mensageria"
)

//Repository - Responsavel por armazenar as mensagens em memória temporariamente
type Repository struct {
	operations map[string]*Operation
}

//Operation - Representa uma operação sobre um conjunto de dados
type Operation struct {
	Collection string
	Action     string
	Field      string
	ID         string
	Messages   OperationDataSet
}

//OperationDataSet - Conjunto de operações que deve ser realizada sobre as mensagens
type OperationDataSet []mensageria.Message

//Save - Armazena a mensagem em memoria
func (r *Repository) Save(message mensageria.Message) {

	if r.operations == nil {
		r.operations = map[string]*Operation{}
	}

	messageHeaders := message.Headers
	key := messageHeaders.Collection + messageHeaders.Action + messageHeaders.ID + messageHeaders.Field
	operation := r.operations[key]

	if operation == nil {
		operation = &Operation{
			Collection: messageHeaders.Collection,
			Action:     messageHeaders.Action,
			Field:      messageHeaders.Field,
			ID:         messageHeaders.ID,
		}
	}

	operation.Messages = append(operation.Messages, message)
	r.operations[key] = operation
}

//Each - Itera as operações armazenadas no repository
func (r *Repository) Each(callback func(string, Operation)) {
	for key, operation := range r.operations {
		callback(key, *operation)
	}
}

//Reject - Rejeita todas as operações do repository
func (r *Repository) Reject() {
	r.Each(func(_ string, operation Operation) {
		operation.Messages.Each(func(message mensageria.Message) {
			message.Reject()
		})
	})
}

//Size - Quantidade de operações armazenadas no repository
func (r Repository) Size() int {
	return len(r.operations)
}

//Clone - Clona o repository
func (r *Repository) Clone() Repository {
	operationsClone := make(map[string]*Operation)

	r.Each(func(key string, operation Operation) {
		operationsClone[key] = &operation
	})

	return Repository{operations: operationsClone}
}

//Each - Itera as mensagens do conjunto de operações
func (os OperationDataSet) Each(callback func(mensageria.Message)) {
	for _, message := range os {
		callback(message)
	}
}

//Clear - Esvazia o repositorio
func (r *Repository) Clear() {
	r.Each(func(key string, _ Operation) {
		delete(r.operations, key)
	})
}

//Reject - Rejeita o conjunto de mensagens da operação
func (os *Operation) Reject() {
	os.Messages.Each(func(message mensageria.Message) {
		message.Reject()
	})
}

//Confirm - Confirma o conjunto de mensagens da operação
func (os *Operation) Confirm() {
	os.Messages.Each(func(message mensageria.Message) {
		message.Confirm()
	})
}

//Logger - Loga as mensagens
func (r *Repository) Logger() {
	r.Each(func(_ string, operation Operation) {
		fmt.Println(operation.Collection, "-", operation.Action, "***********")
		operation.Messages.Each(func(it mensageria.Message) {
			fmt.Println(it.Payload)
		})
	})
}
