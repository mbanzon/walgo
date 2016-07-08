package walgo

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

const (
	ContentTypeHeader = "Content-Type"     // Header for specifying content type
	JsonContentType   = "application/json" // Content type for JSON payloads
)

// Checks the provided error. If it is nil the provided v is encoded as
// JSON and written to the given response writer. If it is not nil
// the status code 500 (Internal server error) is instead sent.
func CheckErrOutputJson(err error, w http.ResponseWriter, v interface{}) {
	CheckErr(w, err, http.StatusInternalServerError, func() {
		w.Header().Add(ContentTypeHeader, JsonContentType)
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
