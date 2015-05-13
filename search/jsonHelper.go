package search

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

func getXMLName(d reflect.Value, label string) (string, bool) {
	if d.IsValid() {
		switch d.Type().Kind() {
		case reflect.Struct:
			v, _ := d.Type().FieldByName(label)
			parts := strings.Split(v.Tag.Get("xml"), " ")
			return parts[1], true
		default:
			//fmt.Printf("Not Struct: %v\n", spew.Sdump(d))
		}
	}
	return "", false
}

// kudos to sberry answering this
// http://stackoverflow.com/questions/29689092/is-there-an-easier-way-to-add-a-layer-over-a-json-object-using-golang-json-encod
func wrapJSONInterface(item interface{}) (map[string]interface{}, error) {
	reflectValue := reflect.Indirect(reflect.ValueOf(item))
	if n, ok := getXMLName(reflectValue, "XMLName"); ok {
		if k := reflectValue.FieldByName("Queries"); k.IsValid() {
			for i := 0; i < k.Len(); i++ {
				b, err1 := wrapJSONInterface(k.Index(i).Interface())
				if err1 != nil {
					continue
				}
				k.Index(i).Set(reflect.ValueOf(b))
			}

		}
		return map[string]interface{}{n: item}, nil
	}
	return nil, errors.New("You failed")
}

func wrapJSON(i interface{}) ([]byte, error) {
	wrappedInterface, err := wrapJSONInterface(i)
	if err != nil {
		return nil, err
	}
	return json.Marshal(wrappedInterface)
}

type stringToInterface func(string) interface{}

func unwrapJSON(b []byte, mapper stringToInterface) (interface{}, error) {
	var data map[string]json.RawMessage
	var nestedName string
	var item interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	for name, value := range data {
		nestedName = name
		item = mapper(nestedName)
		json.Unmarshal(value, item)
	}
	reflectStruct := reflect.Indirect(reflect.ValueOf(item))
	typeOfT := reflectStruct.Type()
	if structField, ok := typeOfT.FieldByName("Queries"); ok {
		var queryJSON map[string]map[string]json.RawMessage
		var queriesData []json.RawMessage
		if err := json.Unmarshal(b, &queryJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(queryJSON[nestedName]["queries"], &queriesData); err != nil {
			return nil, err
		}
		childQueries := make([]interface{}, len(queriesData))
		for i, rawJSON := range queriesData {
			child, err := unwrapJSON(rawJSON, mapper)
			if err != nil {
				continue
			}
			childQueries[i] = child
		}
		newItemValue := reflect.New(reflectStruct.Type()).Elem()
		queriesFieldPos := structField.Index[0]
		for i := 0; i < reflectStruct.NumField(); i++ {
			if i == queriesFieldPos {
				newItemValue.Field(i).Set(reflect.ValueOf(childQueries))
			} else {
				newItemValue.Field(i).Set(reflectStruct.Field(i))
			}
		}
		item = newItemValue.Interface()
	}
	return item, nil
}
