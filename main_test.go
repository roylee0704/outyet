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


Conclusion:
- to test response from Third Party Server.
	1. (SETUP) Fake Third Party Server.
	2. (SETUP) So that you can fake the handler(Response) from Faked third party server
	3. Then you can check actual between exptected RESPONSE.


- to test your own Server.
	1. Run your own Server.
	2. Fake a client REQUEST.
	3. Inspect the RESPONSE.
*/
func TestIntegration(t *testing.T) {
	status := statusHandler(404)      // always returns 404
	ts := httptest.NewServer(&status) // fake server (third-party)
	defer ts.Close()

	// override with my own `sleep` implementation
	sleep := make(chan bool)
	pollSleep = func(time.Duration) {
		sleep <- true // notify that I'm start sleeping.
		sleep <- true // notify that I can leave.
	}
	done := make(chan bool)
	pollDone = func() {
		done <- true
	}

	// this is known as  defer closure
	// consider other test might access, and causing dead-lock
	defer func() {
		pollSleep = time.Sleep
		pollDone = func() {}
	}()

	s := NewServer("1.x", ts.URL, 1*time.Millisecond) // polls from fake server

	<-sleep                                  // goto nextline only after I've received a notification that my Poll has make the first isTagged()
	r, _ := http.NewRequest("GET", "/", nil) // fake client with GET requests
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r) // triggers handler with a recorder

	if b := w.Body.String(); !strings.Contains(b, "No.") {
		t.Errorf("body = %q, wanted no", b)
	}

	status = 200 //update status of faked third party server (there's a race)

	<-sleep // allow my own server to exit and goto next loop and call isTagged() again to receive the 200.

	// now, we no longer needs to sleep in testing, we got channel to synchronize
	//time.Sleep(20 * time.Millisecond) // otherwise no time for server to poll again.

	<-done // make sure only my client only request after the state is updated to YES.
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
