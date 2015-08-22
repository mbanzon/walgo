package walgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

func RequireBody(w http.ResponseWriter, r *http.Request, v interface{}, next func()) {
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		t := reflect.TypeOf(v)
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			tmp := make(map[string]interface{})
			err = json.Unmarshal(data, &tmp)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}

			// TODO - check if all v fields is present
		}

		next()
	} else {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
