package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	client *Client
	ctx    context.Context
)

func init() {
	client = New(nil)
	ctx = context.TODO()
}

func TestNew(t *testing.T) {
	client := New(nil)

	if client.BaseURL == nil || client.BaseURL.String() != defaultBaseURL {
		t.Fatalf("want base url: %v, got: %v", defaultBaseURL, client.BaseURL)
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := New(nil)
	req, err := c.NewRequest(http.MethodGet, "/foo")
	testErrNil(t, err)
	if want, got := defaultBaseURL+"/foo", req.URL.String(); want != got {
		t.Fatalf("want request url: %s, got: %s", want, got)
	}
}

func TestClient_Do(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	c := New(nil)
	u, _ := url.Parse(server.URL)
	c.BaseURL = u

	method := http.MethodGet

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Fatalf("want request: %s, got: %s", method, r.Method)
		}
		fmt.Fprint(w, "\"bar\"")
	})

	req, err := c.NewRequest(method, "/foo")
	testErrNil(t, err)
	var got string
	err = c.Do(context.TODO(), req, &got)
	testErrNil(t, err)

	if want := "bar"; want != got {
		t.Fatalf("wrong reply, want: %s, got: %s", want, got)
	}
}

func testErrNil(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func testMethod(t *testing.T, r *http.Request, method string) {
	t.Helper()

	if r.Method != method {
		t.Fatalf("want request: %v, got: %v", method, r.Method)
	}
}
