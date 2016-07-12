package walgo

import (
	"encoding/json"
	"time"
)

type Response interface {
	// Data returns the raw bytes from the body of the response.
	Data() (data []byte)

	// String returns the content of the response body as a string.
	String() (s string)

	// Code returns the response code given by the server.
	Code() (code int)

	// Duraction returns the time it took to make the request and get the
	// response.
	Duration() (duration time.Duration)

	// JSON returns the content of the response body decoded as JSON
	// into the given interface.
	JSON(v interface{}) (err error)

	// Error gives the error that occured during the request - if any.
	Error() (err error)
}

type responseImpl struct {
	data     []byte
	code     int
	duration time.Duration
	err      error
}

func (r responseImpl) Data() (data []byte) {
	return r.data
}

func (r responseImpl) String() (s string) {
	return string(r.data)
}

func (r responseImpl) Code() (code int) {
	return r.code
}

func (r responseImpl) Duration() (duration time.Duration) {
	return r.duration
}

func (r responseImpl) JSON(v interface{}) (err error) {
	err = json.Unmarshal(r.data, v)
	return err
}

func (r responseImpl) Error() (err error) {
	return r.err
}
