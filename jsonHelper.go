package goMarklogicGo

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

func getXMLName(d interface{}, label string) (string, bool) {
	switch reflect.TypeOf(d).Kind() {
	case reflect.Struct:
		v, _ := reflect.TypeOf(d).FieldByName(label)
		parts := strings.Split(v.Tag.Get("xml"), " ")
		return parts[1], true
	}
	return "", false
}

func wrapJSON(item interface{}) ([]byte, error) {
	if n, ok := getXMLName(item, "XMLName"); ok {
		b, err := json.Marshal(map[string]interface{}{n: item})
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return nil, errors.New("You failed")
}
