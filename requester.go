package walgo

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	defaultRequester Requester
)

const (
	userAgentHeader     = "User-Agent"
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	DefaultClientName   = "walgo"
)

func init() {
	defaultRequester = NewRequester(http.DefaultClient, DefaultClientName, "")
}

type Requester interface {
	// Get performs a GET reuqest to the given URL with the given parameters.
	Get(url string, p ParameterMap) (res Response, err error)

	// Post performs a POST request to the given URL with the given
	// parameters and no request body.
	Post(url string, p ParameterMap) (r Response, err error)

	// PostJson performs a POST request to the given URL with the given
	// parameters and the supplied interface type encoded as JSON.
	PostJson(url string, p ParameterMap, v interface{}) (r Response, err error)

	// PostRaw performs a POST request to the given URL with the given
	// parameters and the supplied bytes as the request body.
	PostRaw(url string, p ParameterMap, data []byte) (r Response, err error)

	// PostMultipart performs a POST request to the given URL with the given
	// parameters and the supplied multipart payload encoded as the request
	// body.
	PostMultipart(url string, p ParameterMap, m *MultipartPayload) (r Response, err error)

	// PostValues performs a POST request to the given URL with the given
	// parameters and a body consisting of the supplied values urlencoded.
	PostValues(url string, p ParameterMap, v url.Values) (r Response, err error)

	// Put performs a PUT request to the given URL with the given parameters.
	Put(url string, p ParameterMap) (r Response, err error)

	// PutJson performs a PUT request to the given URL with the given
	// parameters and the supplied interface type encoded as JSON.
	PutJson(url string, p ParameterMap, v interface{}) (r Response, err error)

	// PutRaw performs a PUT request to the given URL with the given
	// parameters and the supplied bytes as the request body.
	PutRaw(url string, p ParameterMap, data []byte) (r Response, err error)

	// PutMultipart performs a PUT request to the given URL with the given
	// parameters and the supplied multipart payload encoded as the request
	// body.
	PutMultipart(url string, p ParameterMap, m *MultipartPayload) (r Response, err error)

	// PutValues performs a PUT request to the given URL with the given
	// parameters and a body consisting of the supplied values urlencoded.
	PutValues(url string, p ParameterMap, v url.Values) (r Response, err error)

	// Delete peforms a DELETE request to the given URL with the given
	// parameters.
	Delete(url string, p ParameterMap) (r Response, err error)

	makeRequest(url string, p ParameterMap, method string, l *payload) (r Response, err error)
}

type requesterImpl struct {
	client    *http.Client
	userAgent string
	authToken string
}

// NewRequester creates a new Requester using the supplied client. Every
// request has the User-Agent header set to the given value and if the
// authentication token differs from "" it is also used as the Authorization
// header value (with the "Bearer "-prefix).
func NewRequester(c *http.Client, userAgent, authToken string) (r Requester) {
	return &requesterImpl{
		client:    c,
		userAgent: userAgent,
		authToken: authToken,
	}
}

func (f *requesterImpl) makeRequest(urlStr string, p ParameterMap, method string, l *payload) (r Response, err error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	query := u.Query()

	if p != nil {
		for k, v := range p {
			query.Add(k, v)
		}
	}

	u.RawQuery = query.Encode()

	startTime := time.Now()
	code := -1
	var output []byte

	buffer := &bytes.Buffer{}

	if l != nil {
		data := l.getData()
		c, err2 := buffer.Write(data)
		if c != len(data) || err2 != nil {
			return nil, errors.New("Error creating data buffer.")
		}
	}

	req, err := http.NewRequest(method, u.String(), buffer)
	if err != nil {
		return nil, err
	}

	if l != nil {
		req.Header.Add(contentTypeHeader, l.getContentType())
	}

	req.Header.Add(userAgentHeader, f.userAgent)
	if "" != f.authToken {
		req.Header.Add(authorizationHeader, bearerPrefix+f.authToken)
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
