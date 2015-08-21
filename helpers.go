package walgo

import (
	"encoding/json"
	"net/http"
)

func CheckOutputJsonErr(w http.ResponseWriter, v interface{}, err error, code int) {
	CheckErr(w, err, code, func() {
		w.Header().Add("Content-Type", "application/json")
		if e := json.NewEncoder(w).Encode(v); e != nil {
			http.Error(w, "Internal error.", http.StatusInternalServerError)
		}
	})
}

func CheckErr(w http.ResponseWriter, err error, code int, next func()) {
	if err != nil {
		http.Error(w, "", code)
	} else {
		next()
	}
}
