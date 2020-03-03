package gometakube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

var (
	client *Client
	ctx    context.Context
	mux    *http.ServeMux
	server *httptest.Server
)

func setup() {
	client = New()
	ctx = context.TODO()
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func TestNew(t *testing.T) {
	client := New()

	if client.BaseURL == nil || client.BaseURL.String() != defaultBaseURL {
		t.Fatalf("want base url: %v, got: %v", defaultBaseURL, client.BaseURL)
	}
}

func TestClient_NewRequest(t *testing.T) {
	c := New()
	req, err := c.NewRequest(http.MethodGet, "/foo", nil)
	testErrNil(t, err)

	if want, got := defaultBaseURL+"/foo", req.URL.String(); want != got {
		t.Fatalf("want request url: %s, got: %s", want, got)
	}

	want := "test-data"
	req, err = c.NewRequest(http.MethodPost, "post-data", &want)
	testErrNil(t, err)
	got := new(string)
	err = json.NewDecoder(req.Body).Decode(&got)
	testErrNil(t, err)
	if *got != "test-data" {
		t.Fatalf("request body encoding error, want: %v, got: %v", want, got)
	}
}

func TestClient_Do(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, "\"bar\"")
	})

	req, err := client.NewRequest(http.MethodGet, "/foo", nil)
	testErrNil(t, err)

	var got string

	_, err = client.Do(ctx, req, &got)
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

func testParseTime(s string) *time.Time {
	ret, _ := time.Parse(time.RFC3339, s)
	return &ret
}
