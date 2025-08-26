package main

import (
	"net/http"
	"testing"

	"github.com/PPRAMANIK62/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "OK", body)
}
