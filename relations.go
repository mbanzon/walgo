package walgo

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ArgumentNotStructErr = errors.New("Argument is not a struct.")
)

type Relations struct {
}

type relation struct {
}

func isStruct(v interface{}) bool {
	t := reflect.TypeOf(v)
	return t.Kind() == reflect.Struct
}

func (r *Relations) Relate(from, to interface{}) (err error) {
	return r.RelateVia(from, to, reflect.TypeOf(to).Name()+"Id", "Id")
}

func (r *Relations) RelateVia(from, to interface{}, fieldFrom, fieldTo string) (err error) {
	if !isStruct(from) || !isStruct(to) {
		return ArgumentNotStructErr
	}

	fromT := reflect.TypeOf(from)
	fromField, found := fromT.FieldByName(fieldFrom)
	if !found {
		return fmt.Errorf("Field not (%s) found in struct (%s)", fieldFrom, fromT.Name())
	}

	toT := reflect.TypeOf(to)
	toField, found := toT.FieldByName(fieldTo)
	if !found {
		return fmt.Errorf("Field not (%s) found in struct (%s)", fieldTo, toT.Name())
	}

	if fromField.Type.Kind() != toField.Type.Kind() {
		return fmt.Errorf("To and from field are different types: %s != %s", fromField.Type.Kind().String(), toField.Type.Kind().String())
	}

	return nil
}
