package walgo

import (
	"fmt"
	"net/http"
	"reflect"
)

func RequireBody(r *http.Request, v interface{}, next func()) {
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		next()
	} else {
		fmt.Println("oh no")
	}
}
