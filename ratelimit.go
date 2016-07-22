package walgo

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	NoAuthorizationHeaderErr = errors.New("No authorization header value.")
	NoBearerPrefixErr        = errors.New("Authorization has no bearer prefix.")
	NoHeaderValueErr         = errors.New("No header value.")
	RateLimitExceededErr     = errors.New("Rate limit exceeded.")
)

type RateLimitHandler struct {
	maxRequests   int
	duration      time.Duration
	requestCounts map[string][]int64
	lock          *sync.Mutex
	handler       http.Handler
	policy        RatePolicy
}

type RatePolicy interface {
	GetClient(*http.Request) (string, error)
}

type IPRatePolicy struct{}

type TokenRatePolicy struct{}

type HeaderRatePolity struct {
	Name string
}

type CookieRatePolicy struct {
	Name string
}

func NewRateLimiter(maxRequests int, duration time.Duration, p RatePolicy, handler http.Handler) (r *RateLimitHandler) {
	return &RateLimitHandler{
		maxRequests:   maxRequests,
		duration:      duration,
		requestCounts: make(map[string][]int64),
		lock:          &sync.Mutex{},
		handler:       handler,
		policy:        p,
	}
}

func (p HeaderRatePolity) GetClient(r *http.Request) (client string, err error) {
	client = r.Header.Get(p.Name)
	if client == "" {
		return "", NoHeaderValueErr
	}

	return client, nil
}

func (p CookieRatePolicy) GetClient(r *http.Request) (client string, err error) {
	cookie, err := r.Cookie(p.Name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

func (p IPRatePolicy) GetClient(r *http.Request) (client string, err error) {
	client, _, err = net.SplitHostPort(r.RemoteAddr)
	return client, err
}

func (p TokenRatePolicy) GetClient(r *http.Request) (client string, err error) {
	authorization := r.Header.Get(authorizationHeader)
	if authorization == "" {
		return "", NoAuthorizationHeaderErr
	}

	if !strings.HasPrefix(authorization, bearerPrefix) {
		return "", NoBearerPrefixErr
	}

	return strings.TrimPrefix(authorization, bearerPrefix), nil
}

func (r *RateLimitHandler) allowed(client string) (allowed bool) {
	r.lock.Lock()

	if counts, ok := r.requestCounts[client]; ok {
		var i int
		var ts int64

		for i, ts = range counts {
			if ts >= time.Now().UnixNano()-int64(r.duration) {
				break
			}
		}

		newCounts := counts[i:]

		if len(newCounts) > r.maxRequests {
			allowed = false
		} else {
			allowed = true
			newCounts = append(newCounts, time.Now().UnixNano())
		}

		r.requestCounts[client] = newCounts
	} else {
		counts := []int64{time.Now().UnixNano()}
		r.requestCounts[client] = counts

		if r.maxRequests > 1 {
			allowed = true
		} else {
			allowed = false
		}
	}

	r.lock.Unlock()
	return allowed
}

func (r *RateLimitHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if client, err := r.policy.GetClient(req); err == nil {
		if r.allowed(client) {
			r.handler.ServeHTTP(w, req)
			return
		}
	}

	http.Error(w, "Too many requests.", http.StatusTooManyRequests)
}

func (r *RateLimitHandler) LimitHandlerFunc(hf http.HandlerFunc) (h http.HandlerFunc) {
	return func(w http.ResponseWriter, req *http.Request) {
		if client, err := r.policy.GetClient(req); err == nil {
			if r.allowed(client) {
				hf(w, req)
				return
			}
		}

		http.Error(w, "Too many requests.", http.StatusTooManyRequests)
	}
}

type RateLimitRequester struct {
	requester     Requester
	limit         int
	duration      time.Duration
	requestCounts []int64
	lock          *sync.Mutex
}

func NewRateLimitRequester(r Requester, limit int, duration time.Duration) (lr Requester) {
	return &RateLimitRequester{
		requester:     r,
		limit:         limit,
		duration:      duration,
		requestCounts: nil,
		lock:          &sync.Mutex{},
	}
}

func (l *RateLimitRequester) allowed() (allowed bool) {
	l.lock.Lock()

	var i int
	var ts int64

	for i, ts = range l.requestCounts {
		if ts >= time.Now().UnixNano()-int64(l.duration) {
			break
		}
	}

	var newCounts []int64
	if len(l.requestCounts) > 0 {
		newCounts = l.requestCounts[i:]
	}

	if len(newCounts) > l.limit {
		allowed = false
	} else {
		allowed = true
		newCounts = append(newCounts, time.Now().UnixNano())
	}

	l.requestCounts = newCounts

	l.lock.Unlock()
	return allowed
}

func (l *RateLimitRequester) Get(url string, p ParameterMap) (res Response, err error) {
	return l.requester.Get(url, p)
}

func (l *RateLimitRequester) Post(url string, p ParameterMap) (r Response, err error) {
	return l.requester.Post(url, p)
}

func (l *RateLimitRequester) PostJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	return l.requester.PostJson(url, p, v)
}

func (l *RateLimitRequester) PostRaw(url string, p ParameterMap, data []byte) (r Response, err error) {
	return l.requester.PostRaw(url, p, data)
}

func (l *RateLimitRequester) PostMultipart(url string, p ParameterMap, m *MultipartPayload) (r Response, err error) {
	return l.requester.PostMultipart(url, p, m)
}

func (l *RateLimitRequester) PostValues(url string, p ParameterMap, v url.Values) (r Response, err error) {
	return l.requester.PostValues(url, p, v)
}

func (l *RateLimitRequester) Put(url string, p ParameterMap) (r Response, err error) {
	return l.requester.Put(url, p)
}

func (l *RateLimitRequester) PutJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	return l.requester.PutJson(url, p, v)
}

func (l *RateLimitRequester) PutRaw(url string, p ParameterMap, data []byte) (r Response, err error) {
	return l.requester.PutRaw(url, p, data)
}

func (l *RateLimitRequester) PutMultipart(url string, p ParameterMap, m *MultipartPayload) (r Response, err error) {
	return l.requester.PutMultipart(url, p, m)
}

func (l *RateLimitRequester) PutValues(url string, p ParameterMap, v url.Values) (r Response, err error) {
	return l.requester.PutValues(url, p, v)
}

func (l *RateLimitRequester) Delete(url string, p ParameterMap) (r Response, err error) {
	return l.requester.Delete(url, p)
}

func (l *RateLimitRequester) makeRequest(url string, p ParameterMap, method string, load *payload) (r Response, err error) {
	if l.allowed() {
		return l.requester.makeRequest(url, p, method, load)
	} else {
		return nil, RateLimitExceededErr
	}
}
