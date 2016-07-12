package walgo

import (
	"net/http"
	"net/url"
	"strconv"
)

// Maps parameters used in query string in HTTP requests.
type ParameterMap map[string]string

// Adds a string parameter to the parameter map.
func (p ParameterMap) AddString(key, value string) {
	p[key] = value
}

// Adds an integer parameter to the parameter map.
func (p ParameterMap) AddInt(key string, value int) {
	p[key] = strconv.Itoa(value)
}

// Get performs the Get function on the default requester.
func Get(url string, p ParameterMap) (res Response, err error) {
	return defaultRequester.Get(url, p)
}

// Post performs the Post functions on the default requester.
func Post(url string, p ParameterMap) (r Response, err error) {
	return defaultRequester.Post(url, p)
}

// PostRaw performs the PostRaw function on the default requester.
func PostRaw(url string, p ParameterMap, data []byte) (r Response, err error) {
	return defaultRequester.PostRaw(url, p, data)
}

// PostMultipart performs the PostMultipart function on the default requester.
func PostMultipart(url string, p ParameterMap, m *MultipartPayload) (r Response, err error) {
	return defaultRequester.PostMultipart(url, p, m)
}

// PostValues performs the PostValues function on the default requester.
func PostValues(url string, p ParameterMap, v url.Values) (r Response, err error) {
	return defaultRequester.PostValues(url, p, v)
}

// PostJson performs the PostJson function on the default requester.
func PostJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	return defaultRequester.PostJson(url, p, v)
}

// Put performs the Put function on the default requester.
func Put(url string, p ParameterMap) (r Response, err error) {
	return defaultRequester.Put(url, p)
}

// PutRaw performs the PutRaw function on the default requester.
func PutRaw(url string, p ParameterMap, data []byte) (r Response, err error) {
	return defaultRequester.PutRaw(url, p, data)
}

// PutMultipart performs the PutMultipart function on the default requester.
func PutMultipart(url string, p ParameterMap, m *MultipartPayload) (r Response, err error) {
	return defaultRequester.PutMultipart(url, p, m)
}

// PutValues performs the PutValues function on the default requester.
func PutValues(url string, p ParameterMap, v url.Values) (r Response, err error) {
	return defaultRequester.PutValues(url, p, v)
}

// PutJson performs the PutJson function on the default requster.
func PutJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	return defaultRequester.PutJson(url, p, v)
}

// Delete performs the Delete function on the default requster.
func Delete(url string, p ParameterMap) (r Response, err error) {
	return defaultRequester.Delete(url, p)
}

func (f *requesterImpl) Get(url string, p ParameterMap) (res Response, err error) {
	return f.makeRequest(url, p, http.MethodGet, nil)
}

func (f *requesterImpl) Post(url string, p ParameterMap) (r Response, err error) {
	return f.makeRequest(url, p, http.MethodPost, nil)
}

func (f *requesterImpl) PostRaw(url string, p ParameterMap, data []byte) (r Response, err error) {
	return f.makeRequest(url, p, http.MethodPost, payloadFromRawData(data))
}

func (f *requesterImpl) PostMultipart(url string, p ParameterMap, m *MultipartPayload) (r Response, err error) {
	payload, err := payloadFromMultipart(m)
	if err != nil {
		return nil, err
	}

	return f.makeRequest(url, p, http.MethodPost, payload)
}

func (f *requesterImpl) PostValues(url string, p ParameterMap, v url.Values) (r Response, err error) {
	return f.makeRequest(url, p, http.MethodPost, payloadFromValues(v))
}

func (f *requesterImpl) PostJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	payload, err := createJsonPayload(v)
	if err != nil {
		return nil, err
	}

	return f.makeRequest(url, p, http.MethodPost, payload)
}

func (f *requesterImpl) Put(url string, p ParameterMap) (r Response, err error) {
	return f.makeRequest(url, p, http.MethodPut, nil)
}

func (f *requesterImpl) PutRaw(url string, p ParameterMap, data []byte) (r Response, err error) {
	return f.makeRequest(url, p, http.MethodPut, payloadFromRawData(data))
}

func (f *requesterImpl) PutMultipart(url string, p ParameterMap, m *MultipartPayload) (r Response, err error) {
	payload, err := payloadFromMultipart(m)
	if err != nil {
		return nil, err
	}

	return f.makeRequest(url, p, http.MethodPut, payload)
}

func (f *requesterImpl) PutValues(url string, p ParameterMap, v url.Values) (r Response, err error) {
	return f.makeRequest(url, p, http.MethodPut, payloadFromValues(v))
}

func (f *requesterImpl) PutJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	payload, err := createJsonPayload(v)
	if err != nil {
		return nil, err
	}

	return f.makeRequest(url, p, http.MethodPut, payload)
}

func (f *requesterImpl) Delete(url string, p ParameterMap) (r Response, err error) {
	return f.makeRequest(url, p, http.MethodDelete, nil)
}
