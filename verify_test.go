package walgo

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type verificationType struct {
	Foo   string
	Bar   int    `json:"bar"`
	NoFoo string `walgo:"skip"`
}

func TestVerifyBody(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteString(`{"Foo":"Foobar", "bar":42}`)
	r, err := http.NewRequest(http.MethodPost, "", buffer)
	if err != nil {
		t.Fatal(err)
	}

	w := &httptest.ResponseRecorder{}
	var v verificationType

	gotIn := false

	VerifyBody(w, r, &v, func() {
		gotIn = true
	})

	if !gotIn {
		t.Fatal("Didn't get in - verification failed!")
	}

	if w.Code != 0 {
		t.Fatal("Response code shouldn't be set:", w.Code)
	}
}

func TestVerifyBodyFail(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteString(`{"NotFoo":"Foobar", "NotBar":42}`)
	r, err := http.NewRequest(http.MethodPost, "", buffer)
	if err != nil {
		t.Fatal(err)
	}

	w := &httptest.ResponseRecorder{}
	var v verificationType

	gotIn := false

	VerifyBody(w, r, &v, func() {
		gotIn = true
	})

	if gotIn {
		t.Fatal("Got in - verification should not succeed!")
	}

	if http.StatusBadRequest != w.Code {
		t.Fatal("Bad response code:", w.Code)
	}
}

func TestVerifyBodyInvalidType(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteString(`{}`)
	r, err := http.NewRequest(http.MethodPost, "", buffer)
	if err != nil {
		t.Fatal(err)
	}

	w := &httptest.ResponseRecorder{}
	var v string

	gotIn := false

	VerifyBody(w, r, &v, func() {
		gotIn = true
	})

	if gotIn {
		t.Fatal("Got in - verification should not succeed!")
	}

	if http.StatusBadRequest != w.Code {
		t.Fatal("Bad response code:", w.Code)
	}
}
