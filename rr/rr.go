//rr is a package to test httptest.ResponseRecoder

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

func main() {
	testResponseRecorder()
}

// testResponseRecorder shows that you could control req and resp at ease.
func testResponseRecorder() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}

	req, err := http.NewRequest("GET", "http://example.com/foo", nil)

	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()

	// manually triggers handler, acting as if client make a request
	handler(w, req)
	fmt.Printf("%d - %s", w.Code, w.Body.String())
}
