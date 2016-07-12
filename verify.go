package walgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

const (
	jsonTagName  = "json"
	walgoTagName = "walgo"
	skipValue    = "skip"
)

// VarifyBody reads the body from the HTTP request and tries to decode it as
// JSON. It also checks for the presence of all the values in the given
// interface type. If the parsed body matches the interface the next function
// is called.
func VerifyBody(w http.ResponseWriter, r *http.Request, v interface{}, next func()) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	valid, err := verifyData(data, v)
	if err != nil || !valid {
		log.Println(err)
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
		}

		err = json.Unmarshal(data, &v)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, fmt.Errorf("Value is not a struct: %s", reflect.TypeOf(v).Kind().String())
}

func fieldNamePresent(f reflect.StructField, t map[string]interface{}) bool {
	_, ok := t[f.Name]
	return ok
}

func fieldJsonTagPresent(f reflect.StructField, t map[string]interface{}) bool {
	jsonTag := f.Tag.Get(jsonTagName)
	if "" != jsonTag {
		splitTag := strings.Split(jsonTag, ",")

		_, ok := t[splitTag[0]]
		return ok
	}

	return false
}

func skipFieldTagPresent(f reflect.StructField) bool {
	walgoTag := f.Tag.Get(walgoTagName)
	if skipValue == walgoTag {
		return true
	}

	return false
}
