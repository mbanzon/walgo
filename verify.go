package walgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

const (
	jsonTagName  = "json"
	walgoTagName = "walgo"
	skipValue    = "skip"
	noDefault    = "nodefault"
)

// VarifyBody reads the body from the HTTP request and tries to decode it as
// JSON. It also checks for the presence of all the values in the given
// interface type. If the parsed body matches the interface the next function
// is called.
func VerifyBody(w http.ResponseWriter, r *http.Request, v interface{}, next func()) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	valid, err := verifyData(data, v)
	if err != nil || !valid {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	next()
}

func verifyData(data []byte, v interface{}) (ok bool, err error) {
	if reflect.Indirect(reflect.ValueOf(v)).Type().Kind() == reflect.Struct {
		tmp := make(map[string]interface{})
		err = json.Unmarshal(data, &tmp)
		if err != nil {
			return false, err
		}

		t := reflect.Indirect(reflect.ValueOf(v)).Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			if !skipFieldTagPresent(f) && (!fieldJsonTagPresent(f, tmp) && !fieldNamePresent(f, tmp)) {
				return false, fmt.Errorf("Field not found: %s", f.Name)
			}

			if noDefaultFieldTagPresent(f) {
				if isZero(getFieldValue(f, tmp)) {
					return false, fmt.Errorf("Field not allowed to have default value: %s", f.Name)
				}
			}
		}

		err = json.Unmarshal(data, &v)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, fmt.Errorf("Value is not a struct: %s", reflect.TypeOf(v).Kind().String())
}

func getFieldValue(f reflect.StructField, t map[string]interface{}) reflect.Value {
	name, found := fieldJsonTag(f, t)
	if !found && fieldNamePresent(f, t) {
		return reflect.ValueOf(t[f.Name])
	}

	return reflect.ValueOf(t[name])
}

func fieldNamePresent(f reflect.StructField, t map[string]interface{}) bool {
	_, ok := t[f.Name]
	if !ok {
		ok = jsonFieldNamePresent(f, t)
	}
	return ok
}

func jsonFieldNamePresent(f reflect.StructField, t map[string]interface{}) bool {
	if name, found := fieldJsonTag(f, t); found {
		_, ok := t[name]
		return ok
	}

	return false
}

func fieldJsonTagPresent(f reflect.StructField, t map[string]interface{}) bool {
	_, ok := fieldJsonTag(f, t)
	return ok
}

func fieldJsonTag(f reflect.StructField, t map[string]interface{}) (name string, found bool) {
	jsonTag := f.Tag.Get(jsonTagName)
	if "" != jsonTag {
		name = strings.Split(jsonTag, ",")[0]
		_, found = t[name]
		return name, found
	}

	return "", false
}

func skipFieldTagPresent(f reflect.StructField) bool {
	walgoTag := f.Tag.Get(walgoTagName)
	if skipValue == walgoTag {
		return true
	}

	return false
}

func noDefaultFieldTagPresent(f reflect.StructField) bool {
	walgoTag := f.Tag.Get(walgoTagName)
	if noDefault == walgoTag {
		return true
	}

	return false
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}
