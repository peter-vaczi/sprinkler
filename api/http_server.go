package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/peter.vaczi/sprinklerd/core"
	"github.com/peter.vaczi/sprinklerd/utils"
)

// API represents the http rest api of sprinklerd
type API interface {
	Close()
}

type HttpResponse struct {
	Error error
	Body  interface{}
}

type HttpMsg struct {
	ResponseChan chan HttpResponse
}

type HttpDeviceList struct {
	HttpMsg
}

type HttpDeviceAdd struct {
	HttpMsg
	Device *core.Device
}

type HttpDeviceGet struct {
	HttpMsg
	Name string
}

type HttpDeviceDel struct {
	HttpMsg
	Name string
}

type HttpDeviceSet struct {
	HttpMsg
	Name   string
	Device *core.Device
}

type HttpProgramList struct {
	HttpMsg
}

type HttpProgramCreate struct {
	HttpMsg
	Program *core.Program
}

type HttpProgramGet struct {
	HttpMsg
	Name string
}

type HttpProgramDel struct {
	HttpMsg
	Name string
}

// New returns a new http api instance
func New(daemonSocket string, eventChan chan interface{}) API {
	srv := &httpServer{
		router:    mux.NewRouter().StrictSlash(false),
		eventChan: eventChan,
	}

	ipPort := daemonSocket[strings.LastIndex(daemonSocket, "/")+1:]

	srv.router.HandleFunc("/v1", srv.listDevices).Methods("GET")
	srv.router.HandleFunc("/v1/devices", srv.listDevices).Methods("GET")
	srv.router.HandleFunc("/v1/devices", srv.addDevice).Methods("POST")
	srv.router.HandleFunc("/v1/devices/{name}", srv.getDevice).Methods("GET")
	srv.router.HandleFunc("/v1/devices/{name}", srv.delDevice).Methods("DELETE")
	srv.router.HandleFunc("/v1/devices/{name}", srv.setDevice).Methods("PUT")
	srv.router.HandleFunc("/v1/programs", srv.listPrograms).Methods("GET")
	srv.router.HandleFunc("/v1/programs", srv.createProgram).Methods("POST")
	srv.router.HandleFunc("/v1/programs/{name}", srv.getProgram).Methods("GET")
	srv.router.HandleFunc("/v1/programs/{name}", srv.delProgram).Methods("DELETE")

	srv.server = &http.Server{
		Handler:      srv.router,
		Addr:         ipPort,
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

func (s *httpServer) listDevices(w http.ResponseWriter, r *http.Request) {
	rch := make(chan HttpResponse)
	s.eventChan <- HttpDeviceList{HttpMsg: HttpMsg{ResponseChan: rch}}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) addDevice(w http.ResponseWriter, r *http.Request) {
	rch := make(chan HttpResponse)

	dev := &core.Device{}
	err := json.NewDecoder(r.Body).Decode(dev)
	if err == nil {
		s.eventChan <- HttpDeviceAdd{HttpMsg: HttpMsg{ResponseChan: rch}, Device: dev}
	} else {
		rch <- HttpResponse{Error: err}
	}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) getDevice(w http.ResponseWriter, r *http.Request) {
	rch := make(chan HttpResponse)
	vars := mux.Vars(r)
	name := vars["name"]
	s.eventChan <- HttpDeviceGet{HttpMsg: HttpMsg{ResponseChan: rch}, Name: name}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) delDevice(w http.ResponseWriter, r *http.Request) {
	rch := make(chan HttpResponse)
	vars := mux.Vars(r)
	name := vars["name"]
	s.eventChan <- HttpDeviceDel{HttpMsg: HttpMsg{ResponseChan: rch}, Name: name}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) setDevice(w http.ResponseWriter, r *http.Request) {
	rch := make(chan HttpResponse)
	vars := mux.Vars(r)
	name := vars["name"]

	dev := &core.Device{}
	err := json.NewDecoder(r.Body).Decode(dev)
	if err == nil {
		s.eventChan <- HttpDeviceSet{HttpMsg: HttpMsg{ResponseChan: rch}, Name: name, Device: dev}
	} else {
		rch <- HttpResponse{Error: err}
	}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) listPrograms(w http.ResponseWriter, r *http.Request) {
	rch := make(chan HttpResponse)
	s.eventChan <- HttpProgramList{HttpMsg: HttpMsg{ResponseChan: rch}}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) createProgram(w http.ResponseWriter, r *http.Request) {
	rch := make(chan HttpResponse)

	prg := &core.Program{}
	err := json.NewDecoder(r.Body).Decode(prg)
	if err == nil {
		s.eventChan <- HttpProgramCreate{HttpMsg: HttpMsg{ResponseChan: rch}, Program: prg}
	} else {
		rch <- HttpResponse{Error: err}
	}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) getProgram(w http.ResponseWriter, r *http.Request) {
	rch := make(chan HttpResponse)
	vars := mux.Vars(r)
	name := vars["name"]
	s.eventChan <- HttpProgramGet{HttpMsg: HttpMsg{ResponseChan: rch}, Name: name}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) delProgram(w http.ResponseWriter, r *http.Request) {
	rch := make(chan HttpResponse)
	vars := mux.Vars(r)
	name := vars["name"]
	s.eventChan <- HttpProgramDel{HttpMsg: HttpMsg{ResponseChan: rch}, Name: name}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) handleResponse(w http.ResponseWriter, r *http.Request, rch chan HttpResponse) {
	resp := <-rch

	if resp.Error == nil {
		if resp.Body != nil {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, utils.EncodeJson(resp.Body))
		}
	} else {
		http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
	}
}
