package walgo

import (
	"net/http"
	"net/http/httptest"
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
