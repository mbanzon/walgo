package walgo

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	defaultRequester Requester
)

func init() {
	defaultRequester = NewRequester(http.DefaultClient, "walgo", "")
}

type Requester interface {
	Get(url string, p ParameterMap) (res Response, err error)
	Post(url string, p ParameterMap, l Payload) (r Response, err error)
	PostJson(url string, p ParameterMap, v interface{}) (r Response, err error)
	Put(url string, p ParameterMap, l Payload) (r Response, err error)
	PutJson(url string, p ParameterMap, v interface{}) (r Response, err error)
	Delete(url string, p ParameterMap) (r Response, err error)

	makeRequest(url string, p ParameterMap, method string, l Payload) (r Response, err error)
}

type requesterImpl struct {
	client    *http.Client
	userAgent string
	authToken string
}

func NewRequester(c *http.Client, userAgent, authToken string) (r Requester) {
	return &requesterImpl{
		client:    c,
		userAgent: userAgent,
		authToken: authToken,
	}
}

func (f *requesterImpl) makeRequest(url string, p ParameterMap, method string, l Payload) (r Response, err error) {
	u, err := createParameterUrl(url, p)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	code := -1
	var output []byte

	buffer := &bytes.Buffer{}

	if l != nil {
		data := l.getData()
		c, err := buffer.Write(data)
		if c != len(data) || err != nil {
			return nil, errors.New("Error creating data buffer.")
		}
	}

	req, err := http.NewRequest(method, u.String(), buffer)
	if err != nil {
		return nil, err
	}

	if l != nil {
		req.Header.Add("Content-Type", l.getContentType())
	}

	req.Header.Add("User-Agent", f.userAgent)
	if "" != f.authToken {
		req.Header.Add("Authorization", "Bearer "+f.authToken)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.Body != nil {
		code = resp.StatusCode
		output, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	duration := time.Now().Sub(startTime)

	r = responseImpl{
		data:     output,
		code:     code,
		duration: duration,
	}

	return r, err
}
