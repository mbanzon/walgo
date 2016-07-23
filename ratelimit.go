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
	// NoAuthorizationHeaderErr is the error returned when the TokenRatePolicy
	// can't resolve the client because there is no Authorization header value
	// in the request.
	NoAuthorizationHeaderErr = errors.New("No authorization header value.")

	// NoBearerPrefixErr is the error returned when the TokenRatePolicy
	// can't resolve the client because there is no "Bearer " prefix in the
	// Authorizastion header value in the request.
	NoBearerPrefixErr = errors.New("Authorization has no bearer prefix.")

	// NoHeaderValueErr is the error returned when the HeaderRatePolicy
	// can't get a value for the supplied header.
	NoHeaderValueErr = errors.New("No header value.")

	// RateLimitExceeded is returned from the RateLimitRequester when making
	// requests beyond the specified limit.
	RateLimitExceededErr = errors.New("Rate limit exceeded.")
)

// RateLimitHandler keeps the information about rate limiting and provide
// functions to handle requests or shield http.HandleFunc using the
// limits provided.
type RateLimitHandler struct {
	maxRequests   int
	duration      time.Duration
	requestCounts map[string][]int64
	lock          *sync.Mutex
	handler       http.Handler
	policy        RatePolicy
}

// RatePolicy is the common interface for the different rate limiting policies.
type RatePolicy interface {
	// GetClient returns a string representing the client using the data in
	// the given request by the rules of the implementing policy.
	GetClient(*http.Request) (string, error)
}

// IPRatePolicy rate limits using the clients IP address as client
// identification.
type IPRatePolicy struct{}

// TokenRatePolicy rate limits using the Bearer token set en the
// request Authorization header as client identification.
type TokenRatePolicy struct{}

// HeaderRatePolicy rate limits using the value from the request
// header with the provided name as client identification.
type HeaderRatePolity struct {
	// Name of the header holding the client identification string.
	Name string
}

// CookieRatePolicy rate limits using the value of a cookie as the client
// identification.
type CookieRatePolicy struct {
	// Name of the cookie holding the client identification string.
	Name string
}

// Creates a new rate limiter that accepts maxRequest in the duration given.
// The rate limiter uses the provided policy and redirects requests within the
// limit to the given handler (when itself is used as http.Handler).
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

// GetClient implemenets getting the client id string from the header using
// the HeaderRatePolicy.
func (p HeaderRatePolity) GetClient(r *http.Request) (client string, err error) {
	client = r.Header.Get(p.Name)
	if client == "" {
		return "", NoHeaderValueErr
	}

	return client, nil
}

// GetClient implemenets getting the client id string from the cookie
// value using the CookieRatePolicy.
func (p CookieRatePolicy) GetClient(r *http.Request) (client string, err error) {
	cookie, err := r.Cookie(p.Name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

// GetClient implements getting the client id string from the request
// remote IP address value.
func (p IPRatePolicy) GetClient(r *http.Request) (client string, err error) {
	client, _, err = net.SplitHostPort(r.RemoteAddr)
	return client, err
}

// GetClient implements getting the client id string from the request
// Authorization headers Bearer token.
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

		if len(newCounts) >= r.maxRequests {
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

// ServeHTTP is implemented to satisfy the http.Handler interface. It checks
// if the request should be allowed through using the policy and if that is
// the case it forwards the call to the internal handler.
func (r *RateLimitHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if client, err := r.policy.GetClient(req); err == nil {
		if r.allowed(client) {
			r.handler.ServeHTTP(w, req)
			return
		}
	}

	http.Error(w, "Too many requests.", http.StatusTooManyRequests)
}

// LimitHandlerFunc takes a http.HandlerFunc and wraps it in a rate limited
// version.
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

// RateLimitRequester is used for managing and limiting outgoing requests.
type RateLimitRequester struct {
	requester     Requester
	limit         int
	duration      time.Duration
	requestCounts []int64
	lock          *sync.Mutex
}

// Creates a new Requester that limits requets according to the given limit
// and duration. Requets inside the limit is forwarded to the internal
// requester.
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

	if len(newCounts) >= l.limit {
		allowed = false
	} else {
		allowed = true
		newCounts = append(newCounts, time.Now().UnixNano())
	}

	l.requestCounts = newCounts

	l.lock.Unlock()
	return allowed
}

// Get forwards the request to the internal Requester if it is within the
// rate limit.
func (l *RateLimitRequester) Get(url string, p ParameterMap) (res Response, err error) {
	if l.allowed() {
		return l.requester.Get(url, p)
	} else {
		return nil, RateLimitExceededErr
	}
}

// Post forwards the request to the internal Requester if it is within the
// rate limit.
func (l *RateLimitRequester) Post(url string, p ParameterMap) (r Response, err error) {
	if l.allowed() {
		return l.requester.Post(url, p)
	} else {
		return nil, RateLimitExceededErr
	}
}

// PostJson forwards the request to the internal Requester if it is within the
// rate limit.
func (l *RateLimitRequester) PostJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	if l.allowed() {
		return l.requester.PostJson(url, p, v)
	} else {
		return nil, RateLimitExceededErr
	}
}

// PostRaw forwards the request to the internal Requester if it is within the
// rate limit.
func (l *RateLimitRequester) PostRaw(url string, p ParameterMap, data []byte) (r Response, err error) {
	if l.allowed() {
		return l.requester.PostRaw(url, p, data)
	} else {
		return nil, RateLimitExceededErr
	}
}

// PostMultipart forwards the request to the internal Requester if it is
// within the rate limit.
func (l *RateLimitRequester) PostMultipart(url string, p ParameterMap, m *MultipartPayload) (r Response, err error) {
	if l.allowed() {
		return l.requester.PostMultipart(url, p, m)
	} else {
		return nil, RateLimitExceededErr
	}
}

// PostValues forwards the request to the internal Requester if it is
// within the rate limit.
func (l *RateLimitRequester) PostValues(url string, p ParameterMap, v url.Values) (r Response, err error) {
	if l.allowed() {
		return l.requester.PostValues(url, p, v)
	} else {
		return nil, RateLimitExceededErr
	}
}

// Put forwards the request to the internal Requester if it is
// within the rate limit.
func (l *RateLimitRequester) Put(url string, p ParameterMap) (r Response, err error) {
	if l.allowed() {
		return l.requester.Put(url, p)
	} else {
		return nil, RateLimitExceededErr
	}
}

// PutJson forwards the request to the internal Requester if it is
// within the rate limit.
func (l *RateLimitRequester) PutJson(url string, p ParameterMap, v interface{}) (r Response, err error) {
	if l.allowed() {
		return l.requester.PutJson(url, p, v)
	} else {
		return nil, RateLimitExceededErr
	}
}

// PutRaw forwards the request to the internal Requester if it is
// within the rate limit.
func (l *RateLimitRequester) PutRaw(url string, p ParameterMap, data []byte) (r Response, err error) {
	if l.allowed() {
		return l.requester.PutRaw(url, p, data)
	} else {
		return nil, RateLimitExceededErr
	}
}

// PutMultipart forwards the request to the internal Requester if it is
// within the rate limit.
func (l *RateLimitRequester) PutMultipart(url string, p ParameterMap, m *MultipartPayload) (r Response, err error) {
	if l.allowed() {
		return l.requester.PutMultipart(url, p, m)
	} else {
		return nil, RateLimitExceededErr
	}
}

// PutValues forwards the request to the internal Requester if it is
// within the rate limit.
func (l *RateLimitRequester) PutValues(url string, p ParameterMap, v url.Values) (r Response, err error) {
	if l.allowed() {
		return l.requester.PutValues(url, p, v)
	} else {
		return nil, RateLimitExceededErr
	}
}

// Delete forwards the request to the internal Requester if it is
// within the rate limit.
func (l *RateLimitRequester) Delete(url string, p ParameterMap) (r Response, err error) {
	if l.allowed() {
		return l.requester.Delete(url, p)
	} else {
		return nil, RateLimitExceededErr
	}
}

func (l *RateLimitRequester) makeRequest(url string, p ParameterMap, method string, load *payload) (r Response, err error) {
	if l.allowed() {
		return l.requester.makeRequest(url, p, method, load)
	} else {
		return nil, RateLimitExceededErr
	}
}
