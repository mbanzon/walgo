package walgo

import (
	"net/http"
)

func Get(url string, p ParameterMap) (res Response, err error) {
	return defaultRequester.Get(url, p)
}

func Post(url string, p ParameterMap, l Payload) (r Response, err error) {
	return defaultRequester.Post(url, p, l)
}

func PostJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	return defaultRequester.PostJson(url, p, v)
}

func Put(url string, p ParameterMap, l Payload) (r Response, err error) {
	return defaultRequester.Put(url, p, l)
}

func PutJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	return defaultRequester.PutJson(url, p, v)
}

func Delete(url string, p ParameterMap) (r Response, err error) {
	return defaultRequester.Delete(url, p)
}

func (f *requesterImpl) Get(url string, p ParameterMap) (res Response, err error) {
	return f.makeRequest(url, p, http.MethodGet, nil)
}

func (f *requesterImpl) Post(url string, p ParameterMap, l Payload) (r Response, err error) {
	return f.makeRequest(url, p, http.MethodPost, l)
}

func (f *requesterImpl) PostJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	payload, err := CreateJsonPayload(v)
	if err != nil {
		return nil, err
	}

	return f.Post(url, p, payload)
}

func (f *requesterImpl) Put(url string, p ParameterMap, l Payload) (r Response, err error) {
	return f.makeRequest(url, p, http.MethodPut, l)
}

func (f *requesterImpl) PutJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	payload, err := CreateJsonPayload(v)
	if err != nil {
		return nil, err
	}

	return f.Put(url, p, payload)
}

func (f *requesterImpl) Delete(url string, p ParameterMap) (r Response, err error) {
	return f.makeRequest(url, p, http.MethodDelete, nil)
}
