package walgo

import (
	"strconv"
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

type ParameterMap map[string]string

func (p ParameterMap) AddString(key, value string) {
	p[key] = value
}

func (p ParameterMap) AddInt(key string, value int) {
	p[key] = strconv.Itoa(value)
}
