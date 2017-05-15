package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// API represents the http rest api of sprinklerd
type API interface {
	Close()
}

type HttpResponse struct {
	Error error
	Body  string
}

type HttpMsg struct {
	ResponseChan chan HttpResponse
}

type HttpStatus struct {
	HttpMsg
}

// New returns a new http api instance
func New(eventChan chan interface{}) API {
	srv := &httpServer{
		router:    mux.NewRouter().StrictSlash(false),
		eventChan: eventChan,
	}

	srv.router.HandleFunc("/v1", srv.statusHandler).Methods("GET")

	srv.server = &http.Server{
		Handler:      srv.router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go srv.run()
	return srv
}

type httpServer struct {
	router    *mux.Router
	server    *http.Server
	eventChan chan interface{}
}

func (s *httpServer) run() {
	log.Fatal(s.server.ListenAndServe())
}

func (s *httpServer) Close() {
	s.server.Close()
}

func (s *httpServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	rch := make(chan HttpResponse)
	s.eventChan <- HttpStatus{HttpMsg: HttpMsg{ResponseChan: rch}}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) handleResponse(w http.ResponseWriter, r *http.Request, rch chan HttpResponse) {
	resp := <-rch

	if resp.Error == nil {
		if len(resp.Body) > 0 {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, resp.Body)
		}
	} else {
		http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
	}
}
