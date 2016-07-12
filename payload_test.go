package walgo

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestValueMapData(t *testing.T) {
	v := make(url.Values)
	v.Add("key1", "value1")
	v.Add("key2", "value2")

	res, err := PostValues("http://httpbin.org/post", nil, v)
	testResponse(t, res, err)

	res, err = PutValues("http://httpbin.org/put", nil, v)
	testResponse(t, res, err)
}

func TestRawData(t *testing.T) {
	data := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	res, err := PostRaw("http://httpbin.org/post", nil, data)
	testResponse(t, res, err)

	res, err = PutRaw("http://httpbin.org/put", nil, data)
	testResponse(t, res, err)
}

func TestJsonData(t *testing.T) {
	d := struct {
		Foo string
		Bar string
	}{
		"foo",
		"bar",
	}

	res, err := PostJson("http://httpbin.org/post", nil, d)
	testResponse(t, res, err)

	res, err = PutJson("http://httpbin.org/put", nil, d)
	testResponse(t, res, err)
}

func TestMultipartData(t *testing.T) {
	m := &MultipartPayload{}
	m.Add("key1", "value1")
	m.Add("key2", "value2")

	fileSynth := &bytes.Buffer{}
	fileSynth.WriteString("test file for multipart data")

	m.AddFile("test_file", FormFile{
		File: ioutil.NopCloser(fileSynth),
		Name: "yeah.txt",
	})

	res, err := PostMultipart("http://httpbin.org/post", nil, m)
	testResponse(t, res, err)

	res, err = PutMultipart("http://httpbin.org/put", nil, m)
	testResponse(t, res, err)
}

func TestMultipartDublets(t *testing.T) {
	m := &MultipartPayload{}
	err := m.Add("name", "value")
	if err != nil {
		t.Fatal(err)
	}
	err = m.Add("name", "value")
	if err == nil {
		t.Fatal("Dublets should fail!")
	}

	fileSynth := &bytes.Buffer{}
	fileSynth.WriteString("test file for multipart data")
	formFile := FormFile{File: ioutil.NopCloser(fileSynth), Name: "test.txt"}

	err = m.AddFile("test_file", formFile)
	if err != nil {
		t.Fatal(err)
	}
	err = m.AddFile("test_file", formFile)
	if err == nil {
		t.Fatal("Dublets should fail!")
	}

}

func testResponse(t *testing.T, res Response, err error) {
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
