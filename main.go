package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"
)

// Command-line flags.
var (
	httpAddr   = flag.String("http", "localhost:8080", "Listen address")
	pollPeriod = flag.Duration("poll", 5*time.Second, "Poll period")
	version    = flag.String("version", "1.4", "Go version")
)

const baseChangeURL = "http://code.google.com/p/go/source/detail?r="

func main() {
	flag.Parse()
	changeURL := fmt.Sprintf("%sgo%s", baseChangeURL, *version)
	http.Handle("/", NewServer(*version, changeURL, *pollPeriod))

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

// Server implements the outyet server.
// It serves the user interface (it's an http.Handler)
// and polls the remote repository for changes.
type Server struct {
	version string
	url     string
	period  time.Duration

	//shared
	mu  sync.RWMutex
	yes bool
}

// NewServer returns an intiated outyet server.
func NewServer(version, url string, period time.Duration) *Server {

	s := &Server{version: version, url: url, period: period}
	go s.poll() // this makes the program concurrent.
	return s
}

// poll polls the change URL for the specified period until the tag exists.
// Then it sets the Server's yes field true and exits
func (s *Server) poll() {
	for !isTagged(s.url) {
		time.Sleep(s.period)
	}
	s.mu.Lock()
	s.yes = true
	s.mu.Unlock()
}

// isTagged makes an HTTP HEAD request to the given URL and reports whether it
// returned a 200 OK response.
func isTagged(url string) bool {
	header, err := http.Head(url)

	if err != nil {
		log.Print(err)
		return false
	}
	return header.StatusCode == http.StatusOK
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	data := struct {
		Url     string
		Version string
		Yes     bool
	}{
		s.url,
		s.version,
		s.yes,
	}

	s.mu.RUnlock()
	if err := tmpl.Execute(w, data); err != nil {
		log.Print(err)
	}
}

// tmpl is the HTML template that drives the user interface.
var tmpl = template.Must(template.New("tmpl").Parse(`
	<!DOCTYPE html><html><body><center>
	<h2>Is Go {{.Version}} out yet?</h2>
	<h1>
		{{if .Yes}}
			<a href="{{.Url}}">YES!</a>
		{{else}}
			No. :-(
		{{end}}
	</h1>
	</center></body></html>
`))
