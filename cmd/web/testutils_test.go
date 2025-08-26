package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

// returns an instance of our application struct
// containing mock dependencies
func newTestApplication(t *testing.T) *application {
	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog: log.New(io.Discard, "", 0),
	}
}

type testServer struct {
	*httptest.Server
}

// initializes and returns a new test server
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// initialize a new cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// add the cookie jar to the test server client
	// any response cookies will now be stored and sent with subsequent
	// requests when using this client
	ts.Client().Jar = jar

	// disable redirect following for the test server client by setting
	// a custom CheckRedirect function.
	// this function will be called whenever a 3xx status response
	// is recieved by the client, and by always returning a
	// http.ErrUseLastResponse error it forces the client to immediately
	// return the recieved response
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// Implement a get() method on our custom testServer type.
// This makes a GET request to the given url path using the test server client
// and returns the response status code, headers, and body
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
