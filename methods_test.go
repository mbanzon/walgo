package walgo

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestGet(t *testing.T) {
	res, err := Get("http://httpbin.org/get", nil)
	if err != nil || res.Error() != nil {
		t.Fatal(err)
	}

	if res.Code() != http.StatusOK {
		t.Fatalf("Response code should be %d. Got: %d", http.StatusOK, res.Code())
	}

	if res.Data() == nil || res.String() == "" {
		t.Fatal("No data in response.")
	}

	if res.Duration() <= 0 {
		t.Fatal("No time elapsed during the request.")
	}

	var tmp map[string]interface{}
	err = res.JSON(&tmp)
	if err != nil {
		t.Fatal(err)
	}

	tmpArgs, ok := tmp["args"]
	if !ok {
		t.Fatalf("Response did not contain args: %#v", tmp)
	}

	args, ok := tmpArgs.(map[string]interface{})
	if !ok {
		t.Fatalf("args has the wrong type")
	}

	if len(args) > 0 {
		t.Fatalf("args should be empty")
	}
}

func TestGetWithParameters(t *testing.T) {
	p := make(ParameterMap)
	p.AddString("parameter1", "value1")
	p.AddString("parameter2", "value2")
	p.AddInt("parameter3", 42)

	res, err := Get("http://httpbin.org/get", p)
	if err != nil || res.Error() != nil {
		t.Fatal(err)
	}

	if res.Code() != http.StatusOK {
		t.Fatalf("Response code should be %d. Got: %d", http.StatusOK, res.Code())
	}

	if res.Data() == nil || res.String() == "" {
		t.Fatal("No data in response.")
	}

	var tmp map[string]interface{}
	err = res.JSON(&tmp)
	if err != nil {
		t.Fatal(err)
	}

	tmpArgs, ok := tmp["args"]
	if !ok {
		t.Fatalf("Response did not contain args: %#v", tmp)
	}

	args, ok := tmpArgs.(map[string]interface{})
	if !ok {
		t.Fatalf("args has the wrong type")
	}

	if len(args) == 0 {
		t.Fatalf("args should not be empty")
	}

	p1, ok := args["parameter1"]
	if !ok {
		t.Fatal("parameter1 not in response")
	}
	if p1 != "value1" {
		t.Fatalf("Wrong value for parameter1, got %v expected value1", p1)
	}

	p2, ok := args["parameter2"]
	if !ok {
		t.Fatal("parameter2 not in response")
	}
	if p2 != "value2" {
		t.Fatalf("Wrong value for parameter2, got %v expected value2", p2)
	}

	p3, ok := args["parameter3"]
	if !ok {
		t.Fatal("parameter3 not in response")
	}
	parsed, err := strconv.Atoi(fmt.Sprint(p3))
	if err != nil {
		t.Fatal("Error when parsing parameter3", err)
	}
	if parsed != 42 {
		t.Fatalf("Wrong value for parameter3, got %v expected 42", parsed)
	}
}

func TestPost(t *testing.T) {
	res, err := Post("http://httpbin.org/post", nil, nil)
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

func TestPut(t *testing.T) {
	res, err := Put("http://httpbin.org/put", nil, nil)
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

func TestDelete(t *testing.T) {
	res, err := Delete("http://httpbin.org/delete", nil)
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

func TestInvalidUrl(t *testing.T) {
	_, err := Get("", nil)
	if err == nil {
		t.Fatal("Invalid URL should return an error.")
	}
}

func TestCustomRequester(t *testing.T) {
	r := NewRequester(http.DefaultClient, "Walgo Test", "test123")
	res, err := r.Get("http://httpbin.org/get", nil)

	if err != nil {
		t.Fatal(err)
	}

	if res.Code() != 200 {
		t.Fatal("Unexpected response code:", res.Code())
	}
}
