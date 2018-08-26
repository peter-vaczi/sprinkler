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
func New(daemonSocket string, data *core.Data) API {
	srv := &httpServer{
		router: mux.NewRouter().StrictSlash(false),
		data:   data,
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
	srv.router.HandleFunc("/v1/schedules", srv.listSchedules).Methods("GET")
	srv.router.HandleFunc("/v1/schedules", srv.createSchedule).Methods("POST")
	srv.router.HandleFunc("/v1/schedules/{name}", srv.getSchedule).Methods("GET")
	srv.router.HandleFunc("/v1/schedules/{name}", srv.delSchedule).Methods("DELETE")
	srv.router.HandleFunc("/v1/schedules/{name}", srv.setSchedule).Methods("PUT")

	srv.server = &http.Server{
		Handler:      srv.router,
		Addr:         ipPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv
}

type httpServer struct {
	router *mux.Router
	server *http.Server
	data   *core.Data
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
	s.sendResponse(w, r, nil, s.data.Devices)
}

func (s *httpServer) addDevice(w http.ResponseWriter, r *http.Request) {
	dev := &core.Device{}
	err := json.NewDecoder(r.Body).Decode(dev)
	if err == nil {
		err = s.data.Devices.Add(dev)
	}
	s.sendResponse(w, r, err, nil)
}

func (s *httpServer) getDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	dev, err := s.data.Devices.Get(name)
	s.sendResponse(w, r, err, dev)
}

func (s *httpServer) delDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if s.data.Programs.IsDeviceInUse(name) {
		s.sendResponse(w, r, core.DeviceInUse, nil)
		return
	}
	err := s.data.Devices.Del(name)
	s.sendResponse(w, r, err, nil)
}

func (s *httpServer) setDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	dev := &core.Device{}
	err := json.NewDecoder(r.Body).Decode(dev)
	if err == nil {
		err = s.data.Devices.Set(name, dev)
	}
	s.sendResponse(w, r, err, nil)
}

func (s *httpServer) listPrograms(w http.ResponseWriter, r *http.Request) {
	s.sendResponse(w, r, nil, s.data.Programs)
}

func (s *httpServer) createProgram(w http.ResponseWriter, r *http.Request) {
	prg := &core.Program{}
	err := json.NewDecoder(r.Body).Decode(prg)
	if err == nil {
		err = s.data.Programs.Add(prg)
	}
	s.sendResponse(w, r, err, nil)
}

func (s *httpServer) getProgram(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	prg, err := s.data.Programs.Get(name)
	s.sendResponse(w, r, err, prg)
}

func (s *httpServer) delProgram(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	err := s.data.Programs.Del(name)
	s.sendResponse(w, r, err, nil)
}

func (s *httpServer) startProgram(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	prg, err := s.data.Programs.Get(name)
	if err == nil {
		prg.Start()
	}
	s.sendResponse(w, r, err, nil)
}

func (s *httpServer) stopProgram(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	prg, err := s.data.Programs.Get(name)
	if err == nil {
		prg.Stop()
	}
	s.sendResponse(w, r, err, nil)
}

func (s *httpServer) addDeviceToProgram(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	data := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&data)
	dur, _ := time.ParseDuration(data["duration"])

	if err == nil {
		prg, err := s.data.Programs.Get(name)
		if err != nil {
			s.sendResponse(w, r, err, nil)
			return
		}
		dev, err := s.data.Devices.Get(data["device"])
		if err != nil {
			s.sendResponse(w, r, err, nil)
			return
		}
		err = prg.AddDevice(dev, dur)
	}
	s.sendResponse(w, r, err, nil)
}

func (s *httpServer) delDeviceFromProgram(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	idx, _ := strconv.Atoi(vars["idx"])

	prg, err := s.data.Programs.Get(name)
	if err == nil {
		err = prg.DelDevice(idx)
	}
	s.sendResponse(w, r, err, nil)
}

func (s *httpServer) listSchedules(w http.ResponseWriter, r *http.Request) {
	s.sendResponse(w, r, nil, s.data.Schedules)
}

func (s *httpServer) createSchedule(w http.ResponseWriter, r *http.Request) {
	sch := &core.Schedule{}
	err := json.NewDecoder(r.Body).Decode(sch)
	if err == nil {
		if len(sch.ProgramName) != 0 {
			prg, err := s.data.Programs.Get(sch.ProgramName)
			if err != nil {
				s.sendResponse(w, r, err, nil)
				return
			}
			sch.Program = prg
		}
		err = s.data.Schedules.Add(sch)
	}
	s.sendResponse(w, r, err, nil)
}

func (s *httpServer) getSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	sch, err := s.data.Schedules.Get(name)
	s.sendResponse(w, r, err, sch)
}

func (s *httpServer) delSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	err := s.data.Schedules.Del(name)
	s.sendResponse(w, r, err, nil)
}

func (s *httpServer) setSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	sch := &core.Schedule{}
	err := json.NewDecoder(r.Body).Decode(sch)
	if err == nil {
		if len(sch.ProgramName) != 0 {
			prg, err := s.data.Programs.Get(sch.ProgramName)
			if err != nil {
				s.sendResponse(w, r, err, nil)
				return
			}
			sch.Program = prg
		}
		err = s.data.Schedules.Set(name, sch)
	}
	s.sendResponse(w, r, err, nil)
}

func (s *httpServer) sendResponse(w http.ResponseWriter, r *http.Request, err error, body interface{}) {

	switch err {
	case nil:
		log.Printf("%s %s -> %s ", r.Method, r.URL, "200 OK")
		if body != nil {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, utils.EncodeJson(body))
		}
	case core.NotFound:
		log.Printf("%s %s -> %s ", r.Method, r.URL, err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
	case core.OutOfRange:
		log.Printf("%s %s -> %s ", r.Method, r.URL, err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
	case core.AlreadyExists:
		log.Printf("%s %s -> %s ", r.Method, r.URL, err.Error())
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	case core.DeviceInUse:
		log.Printf("%s %s -> %s ", r.Method, r.URL, err.Error())
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	default:
		log.Printf("%s %s -> %s ", r.Method, r.URL, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
