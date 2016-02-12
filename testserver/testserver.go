package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

func main() {
	tryHttpTest()
}

// tryHttpTest sets up a temporary HTTP server that serves a simple 'Hello' response.
func tryHttpTest() {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// by default, it sents status-code 200.
		// or you could form message with http.Error
		fmt.Fprintln(w, "Hello, golang-syd")
	}))

	defer ts.Close()

	// spins up a default http client to issue GET request
	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("response from server[%s] %s", ts.URL, greeting)
}
