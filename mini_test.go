package mini

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))
	defer srv.Close()

	resp, err := Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if string(resp.Body) != "hello" {
		t.Fatalf("body = %q, want %q", resp.Body, "hello")
	}
}

func TestPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	resp, err := Post(srv.URL, []byte("data"))
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Body) != "data" {
		t.Fatalf("body = %q, want %q", resp.Body, "data")
	}
}

func TestDefaultHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Token"); got != "abc" {
			t.Errorf("X-Token = %q, want %q", got, "abc")
		}
	}))
	defer srv.Close()

	c := New(WithHeader("X-Token", "abc"))
	if _, err := c.Get(srv.URL); err != nil {
		t.Fatal(err)
	}
}

func TestTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		time.Sleep(200 * time.Millisecond)
	}))
	defer srv.Close()

	c := New(WithTimeout(10 * time.Millisecond))
	if _, err := c.Get(srv.URL); err == nil {
		t.Fatal("expected a timeout error")
	}
}
