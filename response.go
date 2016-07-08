package walgo

import (
	"encoding/json"
	"time"
)

type Response interface {
	Data() (data []byte)
	String() (s string)
	Code() (code int)
	Duration() (duration time.Duration)
	JSON(v interface{}) (err error)
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
