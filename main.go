package main

import (
	"github.com/peter.vaczi/sprinklerd/api"
	"github.com/peter.vaczi/sprinklerd/core"
)

var devs *core.Devices

func init() {
	devs = core.New()
	devs.Add("rotor-1")
	dev, _ := devs.Add("rotor-2")
	dev.TurnOn()
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
		event.ResponseChan <- api.HttpResponse{Error: nil, Body: devs}
	}
}
