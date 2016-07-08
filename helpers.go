package walgo

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

const (
	ContentTypeHeader = "Content-Type"     // Header for specifying content type
	JsonContentType   = "application/json" // Content type for JSON payloads
)

// CheckErrOutputJson checks the provided error. If it is nil the provided
// v is encoded as JSON and written to the given response writer. If it is
// not nil the status code 500 (Internal server error) is instead sent.
//
// CheckErr is used to check the error. If JSON encoding fails 500 is sent.
func CheckErrOutputJson(err error, w http.ResponseWriter, v interface{}) {
	CheckErr(w, err, func() {
		w.Header().Add(ContentTypeHeader, JsonContentType)
		if e := json.NewEncoder(w).Encode(v); e != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
		}
	})
}

// CheckErr checks the provided error and if it is not nil it sends the
// status code 500 (Internal server error) to the given ResponseWriter.
//
// There is a special case where the error provided is the sql.ErrNoRows.
// In that case the status code 404 is returned.
func CheckErr(w http.ResponseWriter, err error, next func()) {
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "", http.StatusNotFound)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
	} else {
		next()
	}
}
