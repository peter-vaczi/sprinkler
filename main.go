package main

import (
	"github.com/peter.vaczi/sprinklerd/api"
)

type Device struct {
	Name string
	On   bool
}

var Devices []Device

func init() {
	Devices = append(Devices, Device{Name: "rotor-1"})
	Devices = append(Devices, Device{Name: "rotor-2", On: true})
}

func main() {
	mainEvents := make(chan interface{})
	api.New(mainEvents)

	for {
		select {
		case event := <-mainEvents:
			handleEvent(event)
		}
	}
}

func handleEvent(event interface{}) {

	switch event := event.(type) {
	case api.HttpStatus:
		event.ResponseChan <- api.HttpResponse{Error: nil, Body: Devices}
	}
}
