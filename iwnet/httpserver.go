package iwnet

import (
	"log"
	"net/http"
	"strconv"
)

const HTTP_SERVER_PORT = 13000

type HTTPServer struct {
	Name   string
	Logger *log.Logger
}

func NewHTTPServer(logger *log.Logger) *HTTPServer {
	server := &HTTPServer{
		Name:   "HTTPServer",
		Logger: logger,
	}

	return server
}

func (s *HTTPServer) StartAndRunHTTP() {
	http.HandleFunc("/", s.handle_HTTP_request)
	s.Logger.Println(s.Name + " started")
	http.ListenAndServe(":"+strconv.Itoa(HTTP_SERVER_PORT), nil)
}

func (s *HTTPServer) handle_HTTP_request(w http.ResponseWriter, r *http.Request) {
	s.Logger.Println("["+s.Name+"] request:", r.URL.EscapedPath())
}
