package walgo

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckErrNoErr(t *testing.T) {
	calledInner := false
	res := &httptest.ResponseRecorder{}
	CheckErr(res, nil, func() {
		calledInner = true
	})
	if !calledInner {
		t.Fatal("Inline function should be called!")
	}

	if res.Code != 0 {
		t.Fatal("Response code should not be set")
	}
}

func TestCheckErrWithErr(t *testing.T) {
	calledInner := false
	res := &httptest.ResponseRecorder{}
	CheckErr(res, errors.New("dummy"), func() {
		calledInner = true
	})
	if calledInner {
		t.Fatal("Inline function should not be called!")
	}

	if res.Code != http.StatusInternalServerError {
		t.Fatal("Response code incorrect:", res.Code)
	}
}

func TestCheckErrWithNoRows(t *testing.T) {
	calledInner := false
	res := &httptest.ResponseRecorder{}
	CheckErr(res, sql.ErrNoRows, func() {
		calledInner = true
	})
	if calledInner {
		t.Fatal("Inline function should not be called!")
	}

	if res.Code != http.StatusNotFound {
		t.Fatal("Response code incorrect:", res.Code)
	}
}

func TestCheckErrOutputJson(t *testing.T) {
	res := &httptest.ResponseRecorder{}
	data := struct {
		Foo string
		Bar int
	}{
		"foobar",
		42,
	}
	CheckErrOutputJson(nil, res, &data)

	if res.Code != http.StatusOK {
		t.Fatal("Response code incorrect:", res.Code)
	}

	if res.Header().Get(contentTypeHeader) != jsonContentType {
		t.Fatal("Wront content type header:", res.Header().Get(contentTypeHeader))
	}
}
