package dispatch

import (
	"encoding/json"
	"m2m-lazypersistence/internal/pkg/mensageria"
	"m2m-lazypersistence/internal/pkg/repo"
	"m2m-lazypersistence/internal/pkg/util"
	"reflect"
	"strings"

	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

type action func(*mgo.Collection, repo.Operation) error

var actions = map[string]action{
	"INSERT":    insert,
	"PUSH":      push,
	"PULL":      pull,
	"INCREMENT": increment,
}

func executeAction(dispatcher *Dispatcher, operation repo.Operation) error {

	action := actions[strings.ToUpper(operation.Action)]
	collection := selectCollection(dispatcher, operation)
	err := action(collection, operation)

	return err
}

func extractPayload(operation repo.Operation) (payloads []interface{}) {

	operation.Messages.Each(func(message mensageria.Message) {
		payloads = append(payloads, message.Payload)
	})

	return
}

func insert(collection *mgo.Collection, operation repo.Operation) error {
	payloads := extractPayload(operation)
	err := collection.Insert(payloads...)
	return err
}

func push(collection *mgo.Collection, operation repo.Operation) error {
	payloads := extractPayload(operation)

	pushToArray := bson.M{
		"$push": bson.M{
			operation.Field: bson.M{
				"$each": payloads,
			},
		},
	}

	err := resolveUpdate(collection, operation.Create, operation.ID, operation.Condition, pushToArray)

	return err
}

func pull(collection *mgo.Collection, operation repo.Operation) error {
	payloads := extractPayload(operation)

	pullToArray := bson.M{
		"$pull": bson.M{
			operation.Field: bson.M{
				"$in": payloads,
			},
		},
	}

	err := resolveUpdate(collection, operation.Create, operation.ID, operation.Condition, pullToArray)

	return err
}

func increment(collection *mgo.Collection, operation repo.Operation) error {

	increments := bson.M{}

	operation.Messages.Each(func(message mensageria.Message) {

		entries, ok := message.Payload.(map[string]interface{})

		if !ok {
			return
		}

		for entryKey, entryValue := range entries {

			value := entryValue.(float64)

			v, ok := increments[entryKey]

			if !ok {
				if util.CheckDecimal(value) {
					increments[entryKey] = value
				} else {
					increments[entryKey] = int(value)
				}

				continue
			}

			var vConv float64

			if reflect.TypeOf(v) == reflect.TypeOf(1.0) {
				vConv = v.(float64)
			} else {
				vConv = float64(v.(int))
			}

			if util.CheckDecimal(vConv) {
				increments[entryKey] = vConv + value
				continue
			}

			if util.CheckDecimal(value) {
				increments[entryKey] = float64(vConv) + float64(value)
			} else {
				increments[entryKey] = int(vConv) + int(value)
			}
		}
	})

	incrementToField := bson.M{
		"$inc": increments,
	}

	err := resolveUpdate(collection, operation.Create, operation.ID, operation.Condition, incrementToField)

	return err
}

func selectCollection(dispatcher *Dispatcher, operation repo.Operation) (collection *mgo.Collection) {

	var database string

	if operation.Database != "" {
		database = operation.Database
	} else {
		database = dispatcher.Database
	}

	collection = dispatcher.session.DB(database).C(operation.Collection)

	return
}

func checkCondition(conditionValue string, condition *map[string]interface{}) error {

	err := json.Unmarshal([]byte(conditionValue), condition)

	return err
}

func resolveUpdate(collection *mgo.Collection, created bool, id string, conditionValue string, actionFields bson.M) error {

	var idDocument interface{}

	if id != "" {
		if bson.IsObjectIdHex(id) {
			idDocument = bson.ObjectIdHex(id)
		} else {
			idDocument = id
		}

		if created {
			_, err := collection.UpsertId(idDocument, actionFields)
			return err
		}

		return collection.UpdateId(idDocument, actionFields)
	}

	var condition map[string]interface{}

	if conditionValue != "" {
		err := checkCondition(conditionValue, &condition)

		if err != nil {
			return err
		}

		fields := bson.M{}

		for key, value := range condition {

			_, ok := value.(string)

			if ok {
				fields[key] = value
			} else {
				number, ok := value.(float64)
				if ok {
					if util.CheckDecimal(number) {
						fields[key] = value
					} else {
						fields[key] = int(number)
					}
				}
			}

		}

		actionFields["$set"] = fields

		if created {
			_, err := collection.Upsert(condition, actionFields)
			return err
		}

		return collection.Update(condition, actionFields)
	}

	return nil
}
