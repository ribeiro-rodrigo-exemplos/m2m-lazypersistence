package dispatch

import (
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

	err := collection.UpdateId(bson.ObjectIdHex(operation.ID), pushToArray)

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

	err := collection.UpdateId(bson.ObjectIdHex(operation.ID), pullToArray)

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

	err := collection.UpdateId(bson.ObjectIdHex(operation.ID), incrementToField)

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
