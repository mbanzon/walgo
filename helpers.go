package walgo

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func CheckErrOutputJson(err error, w http.ResponseWriter, v interface{}) {
	CheckErr(w, err, http.StatusInternalServerError, func() {
		w.Header().Add("Content-Type", "application/json")
		if e := json.NewEncoder(w).Encode(v); e != nil {
			log.Println(e)
			http.Error(w, "Internal error.", http.StatusInternalServerError)
		}
	})
}

func CheckErr(w http.ResponseWriter, err error, code int, next func()) {
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			http.Error(w, "", http.StatusNotFound)
		} else {
			http.Error(w, "", code)
		}
	} else {
		next()
	}
}
