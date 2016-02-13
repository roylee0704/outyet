package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

/* TestIntegration

1. Setup third-party fake server.
2. Setup own server, with URL from fake server.
3. Simulate GET Request of fake client
4. Prepare a recorder
5. Trigger own server's Handler(Recorder ,GET-req)
6. Check result from Recorder.
*/
func TestIntegration(t *testing.T) {
	status := statusHandler(404)      // always returns 404
	ts := httptest.NewServer(&status) // fake server (third-party)
	defer ts.Close()

	s := NewServer("1.x", ts.URL, 1*time.Millisecond) // polls from fake server

	r, _ := http.NewRequest("GET", "/", nil) // fake client with GET requests
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r) // triggers handler with a recorder

	if b := w.Body.String(); !strings.Contains(b, "No.") {
		t.Errorf("body = %q, wanted no", b)
	}

	status = 200

	time.Sleep(20 * time.Millisecond) // otherwise no time for server to poll again.

	w = httptest.NewRecorder()
	s.ServeHTTP(w, r)

	if b := w.Body.String(); !strings.Contains(b, "YES!") {
		t.Errorf("body = %q, wanted yes!", b)
	}

}

type statusHandler int // this class has only one state

// ServeHTTP always return statuscode from input.
func (s *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(int(*s))
}

// TestIsTagged sends http.HEAD request to fake server, expects true if 200
func TestIsTagged(t *testing.T) {

	status := statusHandler(404)
	ts := httptest.NewServer(&status)
	defer ts.Close()

	if isTagged(ts.URL) {
		t.Error("isTagged returned true, want false")
	}

	status = 200

	if !isTagged(ts.URL) {
		t.Error("isTagged returned false, want true")
	}

}
