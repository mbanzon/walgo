package walgo

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type hitCountHandler struct {
	hitCount int
}

func (d *hitCountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.hitCount++
}

func TestIPRateLimitHandlerFunc(t *testing.T) {
	hitCount := 0
	r := NewRateLimiter(10, time.Hour, IPRatePolicy{}, nil)
	f := r.LimitHandlerFunc(func(resW http.ResponseWriter, req2 *http.Request) {
		hitCount++
	})
	for i := 0; i < 20; i++ {
		req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()

		f(w, req)

		if i < 10 {
			if w.Code != http.StatusOK {
				t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusOK)
			}
		} else {
			if w.Code != http.StatusTooManyRequests {
				t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusTooManyRequests)
			}
		}
	}

	if hitCount > 10 {
		t.Fatal("Hit count too high:", hitCount)
	}
}

func TestZeroLimitHandlerFunc(t *testing.T) {
	r := NewRateLimiter(0, time.Hour, IPRatePolicy{}, nil)
	f := r.LimitHandlerFunc(nil)
	for i := 0; i < 20; i++ {
		req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()

		f(w, req)

		if w.Code != http.StatusTooManyRequests {
			t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusTooManyRequests)
		}
	}
}

func TestIPRateLimit(t *testing.T) {
	h := &hitCountHandler{}
	r := NewRateLimiter(10, time.Hour, IPRatePolicy{}, h)
	for i := 0; i < 20; i++ {
		req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if i < 10 {
			if w.Code != http.StatusOK {
				t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusOK)
			}
		} else {
			if w.Code != http.StatusTooManyRequests {
				t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusTooManyRequests)
			}
		}
	}

	if h.hitCount > 10 {
		t.Fatal("Hit count too high:", h.hitCount)
	}
}

func TestCookieRateLimit(t *testing.T) {
	h := &hitCountHandler{}
	r := NewRateLimiter(10, time.Hour, CookieRatePolicy{Name: "My-Cookie"}, h)
	for i := 0; i < 21; i++ {
		req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		if i < 20 {
			req.AddCookie(&http.Cookie{Name: "My-Cookie", Value: "gabbagabbahey"})
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if i < 10 {
			if w.Code != http.StatusOK {
				t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusOK)
			}
		} else {
			if w.Code != http.StatusTooManyRequests {
				t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusTooManyRequests)
			}
		}
	}

	if h.hitCount > 10 {
		t.Fatal("Hit count too high:", h.hitCount)
	}
}

func TestRateLimitReset(t *testing.T) {
	h := &hitCountHandler{}
	r := NewRateLimiter(10, 2*time.Second, IPRatePolicy{}, h)
	for j := 0; j < 2; j++ {
		for i := 0; i < 20; i++ {
			req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.RemoteAddr = "127.0.0.1:12345"
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if i < 10 {
				if w.Code != http.StatusOK {
					t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusOK)
				}
			} else {
				if w.Code != http.StatusTooManyRequests {
					t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusTooManyRequests)
				}
			}
		}

		time.Sleep(2 * time.Second)
	}

	if h.hitCount != 20 {
		t.Fatal("Wrong hit count:", h.hitCount)
	}
}

func TestTwoClientsRateLimit(t *testing.T) {
	h := &hitCountHandler{}
	r := NewRateLimiter(10, time.Hour, IPRatePolicy{}, h)
	for i := 0; i < 20; i++ {
		req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}

		if i%2 == 0 {
			req.RemoteAddr = "127.0.0.2:12345"
		} else {
			req.RemoteAddr = "127.0.0.1:12345"
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusOK)
		}
	}

	if h.hitCount < 20 {
		t.Fatal("Hit count too low:", h.hitCount)
	}
}

func TestHeaderRateLimit(t *testing.T) {
	h := &hitCountHandler{}
	r := NewRateLimiter(10, time.Hour, HeaderRatePolity{Name: "X-Real-Ip"}, h)
	for i := 0; i < 21; i++ {
		req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}

		if i < 20 {
			req.Header.Add("X-Real-Ip", "127.0.0.1")
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if i < 10 {
			if w.Code != http.StatusOK {
				t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusOK)
			}
		} else {
			if w.Code != http.StatusTooManyRequests {
				t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusTooManyRequests)
			}
		}
	}

	if h.hitCount > 10 {
		t.Fatal("Hit count too high:", h.hitCount)
	}
}

func TestTokenRateLimit(t *testing.T) {
	h := &hitCountHandler{}
	r := NewRateLimiter(10, time.Hour, TokenRatePolicy{}, h)
	for i := 0; i < 22; i++ {
		req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}

		if i < 20 {
			req.Header.Add("Authorization", "Bearer gabbagabbahey")
		} else {
			if i < 21 {
				req.Header.Add("Authorization", "gabbagabbahey")
			}
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if i < 10 {
			if w.Code != http.StatusOK {
				t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusOK)
			}
		} else {
			if w.Code != http.StatusTooManyRequests {
				t.Fatalf("(%d) Wrong code (%d) expected: %d", i, w.Code, http.StatusTooManyRequests)
			}
		}
	}

	if h.hitCount > 10 {
		t.Fatal("Hit count too high:", h.hitCount)
	}
}

func TestLimitedGet(t *testing.T) {
	requester := NewRateLimitRequester(defaultRequester, 1, time.Hour)

	for i := 0; i < 2; i++ {
		res, err := requester.Get("http://httpbin.org/get", nil)
		if i < 1 {
			if err != nil || res.Error() != nil {
				t.Fatal(err)
			}
		} else {
			if err != RateLimitExceededErr {
				t.Fatal("Allowed to request beyond rate limit:", i, err)
			}
		}
	}
}

func TestLimitedPost(t *testing.T) {
	requester := NewRateLimitRequester(defaultRequester, 1, time.Hour)

	for i := 0; i < 2; i++ {
		res, err := requester.Post("http://httpbin.org/post", nil)
		if i < 1 {
			if err != nil || res.Error() != nil {
				t.Fatal(err)
			}
		} else {
			if err != RateLimitExceededErr {
				t.Fatal("Allowed to request beyond rate limit:", i, err)
			}
		}
	}
}

func TestLimitedPostRaw(t *testing.T) {
	requester := NewRateLimitRequester(defaultRequester, 1, time.Hour)

	for i := 0; i < 2; i++ {
		res, err := requester.PostRaw("http://httpbin.org/post", nil, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
		if i < 1 {
			if err != nil || res.Error() != nil {
				t.Fatal(err)
			}
		} else {
			if err != RateLimitExceededErr {
				t.Fatal("Allowed to request beyond rate limit:", i, err)
			}
		}
	}
}

func TestLimitedPostMultipart(t *testing.T) {
	requester := NewRateLimitRequester(defaultRequester, 1, time.Hour)

	m := &MultipartPayload{}
	m.Add("foo", "bar")

	for i := 0; i < 2; i++ {
		res, err := requester.PostMultipart("http://httpbin.org/post", nil, m)
		if i < 1 {
			if err != nil || res.Error() != nil {
				t.Fatal(err)
			}
		} else {
			if err != RateLimitExceededErr {
				t.Fatal("Allowed to request beyond rate limit:", i, err)
			}
		}
	}
}

func TestLimitedPostValues(t *testing.T) {
	requester := NewRateLimitRequester(defaultRequester, 1, time.Hour)

	v := url.Values{}
	v.Add("foo", "bar")

	for i := 0; i < 2; i++ {
		res, err := requester.PostValues("http://httpbin.org/put", nil, v)
		if i < 1 {
			if err != nil || res.Error() != nil {
				t.Fatal(err)
			}
		} else {
			if err != RateLimitExceededErr {
				t.Fatal("Allowed to request beyond rate limit:", i, err)
			}
		}
	}
}

func TestLimitedPostJson(t *testing.T) {
	requester := NewRateLimitRequester(defaultRequester, 1, time.Hour)

	payload := struct{ string }{"foobar"}

	for i := 0; i < 2; i++ {
		res, err := requester.PostJson("http://httpbin.org/post", nil, &payload)
		if i < 1 {
			if err != nil || res.Error() != nil {
				t.Fatal(err)
			}
		} else {
			if err != RateLimitExceededErr {
				t.Fatal("Allowed to request beyond rate limit:", i, err)
			}
		}
	}
}

func TestLimitedPut(t *testing.T) {
	requester := NewRateLimitRequester(defaultRequester, 1, time.Hour)

	for i := 0; i < 2; i++ {
		res, err := requester.Put("http://httpbin.org/put", nil)
		if i < 1 {
			if err != nil || res.Error() != nil {
				t.Fatal(err)
			}
		} else {
			if err != RateLimitExceededErr {
				t.Fatal("Allowed to request beyond rate limit:", i, err)
			}
		}
	}
}

func TestLimitedPutRaw(t *testing.T) {
	requester := NewRateLimitRequester(defaultRequester, 1, time.Hour)

	for i := 0; i < 2; i++ {
		res, err := requester.PutRaw("http://httpbin.org/put", nil, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
		if i < 1 {
			if err != nil || res.Error() != nil {
				t.Fatal(err)
			}
		} else {
			if err != RateLimitExceededErr {
				t.Fatal("Allowed to request beyond rate limit:", i, err)
			}
		}
	}
}

func TestLimitedPutMultipart(t *testing.T) {
	requester := NewRateLimitRequester(defaultRequester, 1, time.Hour)

	m := &MultipartPayload{}
	m.Add("foo", "bar")

	for i := 0; i < 2; i++ {
		res, err := requester.PutMultipart("http://httpbin.org/put", nil, m)
		if i < 1 {
			if err != nil || res.Error() != nil {
				t.Fatal(err)
			}
		} else {
			if err != RateLimitExceededErr {
				t.Fatal("Allowed to request beyond rate limit:", i, err)
			}
		}
	}
}

func TestLimitedPutValues(t *testing.T) {
	requester := NewRateLimitRequester(defaultRequester, 1, time.Hour)

	v := url.Values{}
	v.Add("foo", "bar")

	for i := 0; i < 2; i++ {
		res, err := requester.PutValues("http://httpbin.org/put", nil, v)
		if i < 1 {
			if err != nil || res.Error() != nil {
				t.Fatal(err)
			}
		} else {
			if err != RateLimitExceededErr {
				t.Fatal("Allowed to request beyond rate limit:", i, err)
			}
		}
	}
}

func TestLimitedPutJson(t *testing.T) {
	requester := NewRateLimitRequester(defaultRequester, 1, time.Hour)

	payload := struct{ string }{"foobar"}

	for i := 0; i < 2; i++ {
		res, err := requester.PutJson("http://httpbin.org/put", nil, &payload)
		if i < 1 {
			if err != nil || res.Error() != nil {
				t.Fatal(err)
			}
		} else {
			if err != RateLimitExceededErr {
				t.Fatal("Allowed to request beyond rate limit:", i, err)
			}
		}
	}
}
