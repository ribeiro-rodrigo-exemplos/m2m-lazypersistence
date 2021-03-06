package repo

import (
	"crypto/md5"
	"encoding/hex"
	"m2m-lazypersistence/internal/pkg/mensageria"
)

//Repository - Responsavel por armazenar as mensagens em memória temporariamente
type Repository struct {
	operations map[string]*Operation
	size       int
}

//Operation - Representa uma operação sobre um conjunto de dados
type Operation struct {
	Database   string
	Collection string
	Action     string
	Field      string
	ID         string
	Condition  string
	Create     bool
	Messages   OperationDataSet
}

//OperationDataSet - Conjunto de operações que deve ser realizada sobre as mensagens
type OperationDataSet []mensageria.Message

//Save - Armazena a mensagem em memoria
func (r *Repository) Save(request mensageria.RequestPersistence) {

	if r.operations == nil {
		r.operations = map[string]*Operation{}
	}

	database, collection, action, id, condition, field, create := extractHeaders(request.Headers)

	if isNotValidRequest(collection, action) {
		request.Message.Discard()
		return
	}

	conditionKey := generateConditionKey(condition)

	key := database + collection + action + id + conditionKey + field
	operation := r.operations[key]

	if operation == nil {
		operation = &Operation{
			Database:   database,
			Collection: collection,
			Action:     action,
			Field:      field,
			ID:         id,
			Condition:  condition,
		}
	}

	operation.Messages = append(operation.Messages, request.Message)
	operation.Create = create

	r.operations[key] = operation
	r.size++
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
	return r.size
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

	r.size = 0
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

func extractHeaders(headers map[string]interface{}) (database, collection, action, id,
	condition, field string, create bool) {

	database = extract(headers["database"])
	collection = extract(headers["collection"])
	action = extract(headers["action"])
	id = extract(headers["id"])
	condition = extract(headers["condition"])
	field = extract(headers["field"])
	create = extractBool(headers["create"])

	return
}

func extract(value interface{}) string {
	v, ok := value.(string)

	if ok {
		return v
	}

	return ""
}

func extractBool(value interface{}) bool {
	v, ok := value.(bool)

	if ok {
		return v
	}

	s, ok := value.(string)

	if ok {
		return s == "true"
	}

	return false
}

func isNotValidRequest(collection, action string) bool {
	return collection == "" || action == ""
}

func generateConditionKey(conditionValue string) string {

	if conditionValue != "" {
		hasher := md5.New()
		hasher.Write([]byte(conditionValue))
		return hex.EncodeToString(hasher.Sum(nil))
	}

	return ""
}
