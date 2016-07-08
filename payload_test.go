package walgo

import (
	"net/http"
	"net/url"
	"testing"
)

func TestValueMapData(t *testing.T) {
	v := make(url.Values)
	v.Add("key1", "value1")
	v.Add("key2", "value2")

	p := PayloadFromValues(v)

	doMakeTheRequest(t, p)
}

func TestJsonData(t *testing.T) {
	d := struct {
		Foo string
		Bar string
	}{
		"foo",
		"bar",
	}

	p, err := CreateJsonPayload(d)
	if err != nil {
		t.Fatal(err)
	}

	doMakeTheRequest(t, p)
}

func TestMultipartData(t *testing.T) {
	m := &MultipartPayload{}
	m.Add("key1", "value1")
	m.Add("key2", "value2")

	p, err := PayloadFromMultipart(m)
	if err != nil {
		t.Fatal(err)
	}

	doMakeTheRequest(t, p)
}

func TestMultipartDublets(t *testing.T) {
	m := MultipartPayload{}
	err := m.Add("name", "value")
	if err != nil {
		t.Fatal(err)
	}
	err = m.Add("name", "value")
	if err == nil {
		t.Fatal("Dublets should fail!")
	}
}

func doMakeTheRequest(t *testing.T, p Payload) {
	res, err := Post("http://httpbin.org/post", nil, p)
	if err != nil || res.Error() != nil {
		t.Fatal(err)
	}

	if res.Code() != http.StatusOK {
		t.Fatalf("Response code should be %d. Got: %d", http.StatusOK, res.Code())
	}

	if res.Data() == nil || res.String() == "" {
		t.Fatal("No data in response.")
	}
}
