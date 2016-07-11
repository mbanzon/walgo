package walgo

import (
	"net/http"
	"net/url"
	"strconv"
)

type ParameterMap map[string]string

func (p ParameterMap) AddString(key, value string) {
	p[key] = value
}

func (p ParameterMap) AddInt(key string, value int) {
	p[key] = strconv.Itoa(value)
}

func Get(url string, p ParameterMap) (res Response, err error) {
	return defaultRequester.Get(url, p)
}

func Post(url string, p ParameterMap) (r Response, err error) {
	return defaultRequester.Post(url, p)
}

func PostRaw(url string, p ParameterMap, data []byte) (r Response, err error) {
	return defaultRequester.PostRaw(url, p, data)
}

func PostMultipart(url string, p ParameterMap, m *MultipartPayload) (r Response, err error) {
	return defaultRequester.PostMultipart(url, p, m)
}

func PostValues(url string, p ParameterMap, v url.Values) (r Response, err error) {
	return defaultRequester.PostValues(url, p, v)
}

func PostJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	return defaultRequester.PostJson(url, p, v)
}

func Put(url string, p ParameterMap) (r Response, err error) {
	return defaultRequester.Put(url, p)
}

func PutRaw(url string, p ParameterMap, data []byte) (r Response, err error) {
	return defaultRequester.PutRaw(url, p, data)
}

func PutMultipart(url string, p ParameterMap, m *MultipartPayload) (r Response, err error) {
	return defaultRequester.PutMultipart(url, p, m)
}

func PutValues(url string, p ParameterMap, v url.Values) (r Response, err error) {
	return defaultRequester.PutValues(url, p, v)
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
