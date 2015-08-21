package walgo

func Get(url string, p ParameterMap) (res Response, err error) {
	return defaultRequester.Get(url, p)
}

func Post(url string, p ParameterMap, l Payload) (r Response, err error) {
	return defaultRequester.Post(url, p, l)
}

func Put(url string, p ParameterMap, l Payload) (r Response, err error) {
	return defaultRequester.Put(url, p, l)
}

func Delete(url string, p ParameterMap) (r Response, err error) {
	return defaultRequester.Delete(url, p)
}

func (f *requesterImpl) Get(url string, p ParameterMap) (res Response, err error) {
	return f.makeRequest(url, p, "GET", nil)
}

func (f *requesterImpl) Post(url string, p ParameterMap, l Payload) (r Response, err error) {
	return f.makeRequest(url, p, "POST", l)
}

func (f *requesterImpl) Put(url string, p ParameterMap, l Payload) (r Response, err error) {
	return f.makeRequest(url, p, "PUT", l)
}

func (f *requesterImpl) Delete(url string, p ParameterMap) (r Response, err error) {
	return f.makeRequest(url, p, "DELETE", nil)
}
