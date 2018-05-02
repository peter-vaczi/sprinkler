package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/peter-vaczi/sprinkler/core"
	"github.com/peter-vaczi/sprinkler/utils"
)

// API represents the http rest api of sprinkler
type API interface {
	Run()
	ServeHTTP(http.ResponseWriter, *http.Request)
	Close()
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
	srv.router.HandleFunc("/v1/programs/{name}/start", srv.startProgram).Methods("POST")
	srv.router.HandleFunc("/v1/programs/{name}/stop", srv.stopProgram).Methods("POST")
	srv.router.HandleFunc("/v1/programs/{name}/devices", srv.addDeviceToProgram).Methods("POST")
	srv.router.HandleFunc("/v1/programs/{name}/devices/{idx}", srv.delDeviceFromProgram).Methods("DELETE")

	srv.server = &http.Server{
		Handler:      srv.router,
		Addr:         ipPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv
}

type httpServer struct {
	router    *mux.Router
	server    *http.Server
	eventChan chan interface{}
}

func (s *httpServer) Run() {
	log.Fatal(s.server.ListenAndServe())
}

func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *httpServer) Close() {
	s.server.Close()
}

func (s *httpServer) listDevices(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse)
	s.eventChan <- core.MsgDeviceList{MsgRequest: core.MsgRequest{ResponseChan: rch}}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) addDevice(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse)

	dev := &core.Device{}
	err := json.NewDecoder(r.Body).Decode(dev)
	if err == nil {
		s.eventChan <- core.MsgDeviceAdd{MsgRequest: core.MsgRequest{ResponseChan: rch}, Device: dev}
	} else {
		rch <- core.MsgResponse{Error: err}
	}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) getDevice(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse)
	vars := mux.Vars(r)
	name := vars["name"]
	s.eventChan <- core.MsgDeviceGet{MsgRequest: core.MsgRequest{ResponseChan: rch}, Name: name}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) delDevice(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse)
	vars := mux.Vars(r)
	name := vars["name"]
	s.eventChan <- core.MsgDeviceDel{MsgRequest: core.MsgRequest{ResponseChan: rch}, Name: name}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) setDevice(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse)
	vars := mux.Vars(r)
	name := vars["name"]

	dev := &core.Device{}
	log.Printf("before parse")
	err := json.NewDecoder(r.Body).Decode(dev)
	log.Printf("after parse %v", err)
	if err == nil {
		s.eventChan <- core.MsgDeviceSet{MsgRequest: core.MsgRequest{ResponseChan: rch}, Name: name, Device: dev}
	} else {
		rch <- core.MsgResponse{Error: err}
	}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) listPrograms(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse)
	s.eventChan <- core.MsgProgramList{MsgRequest: core.MsgRequest{ResponseChan: rch}}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) createProgram(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse)

	prg := &core.Program{}
	err := json.NewDecoder(r.Body).Decode(prg)
	if err == nil {
		s.eventChan <- core.MsgProgramCreate{MsgRequest: core.MsgRequest{ResponseChan: rch}, Program: prg}
	} else {
		rch <- core.MsgResponse{Error: err}
	}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) getProgram(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse)
	vars := mux.Vars(r)
	name := vars["name"]
	s.eventChan <- core.MsgProgramGet{MsgRequest: core.MsgRequest{ResponseChan: rch}, Name: name}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) delProgram(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse)
	vars := mux.Vars(r)
	name := vars["name"]
	s.eventChan <- core.MsgProgramDel{MsgRequest: core.MsgRequest{ResponseChan: rch}, Name: name}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) startProgram(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse)
	vars := mux.Vars(r)
	name := vars["name"]
	s.eventChan <- core.MsgProgramStart{MsgRequest: core.MsgRequest{ResponseChan: rch}, Name: name}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) stopProgram(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse)
	vars := mux.Vars(r)
	name := vars["name"]
	s.eventChan <- core.MsgProgramStop{MsgRequest: core.MsgRequest{ResponseChan: rch}, Name: name}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) addDeviceToProgram(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse, 1)
	vars := mux.Vars(r)
	name := vars["name"]

	data := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&data)
	dur, _ := time.ParseDuration(data["duration"])

	if err == nil {
		s.eventChan <- core.MsgProgramAddDevice{MsgRequest: core.MsgRequest{ResponseChan: rch}, Program: name, Device: data["device"], Duration: dur}
	} else {
		rch <- core.MsgResponse{Error: err}
	}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) delDeviceFromProgram(w http.ResponseWriter, r *http.Request) {
	rch := make(chan core.MsgResponse)
	vars := mux.Vars(r)
	name := vars["name"]
	idx, _ := strconv.Atoi(vars["idx"])

	s.eventChan <- core.MsgProgramDelDevice{MsgRequest: core.MsgRequest{ResponseChan: rch}, Program: name, Idx: idx}
	s.handleResponse(w, r, rch)
}

func (s *httpServer) handleResponse(w http.ResponseWriter, r *http.Request, rch chan core.MsgResponse) {
	resp := <-rch

	switch resp.Error {
	case nil:
		log.Printf("%s %s -> %s ", r.Method, r.URL, "200 OK")
		if resp.Body != nil {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, utils.EncodeJson(resp.Body))
		}
	case core.NotFound:
		log.Printf("%s %s -> %s ", r.Method, r.URL, resp.Error.Error())
		http.Error(w, resp.Error.Error(), http.StatusNotFound)
	case core.OutOfRange:
		log.Printf("%s %s -> %s ", r.Method, r.URL, resp.Error.Error())
		http.Error(w, resp.Error.Error(), http.StatusNotFound)
	case core.AlreadyExists:
		log.Printf("%s %s -> %s ", r.Method, r.URL, resp.Error.Error())
		http.Error(w, resp.Error.Error(), http.StatusNotAcceptable)
	case core.DeviceInUse:
		log.Printf("%s %s -> %s ", r.Method, r.URL, resp.Error.Error())
		http.Error(w, resp.Error.Error(), http.StatusNotAcceptable)
	default:
		log.Printf("%s %s -> %s ", r.Method, r.URL, resp.Error.Error())
		http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
	}
}
