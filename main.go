package main

import (
	"github.com/peter.vaczi/sprinklerd/api"
)

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
		event.ResponseChan <- api.HttpResponse{Error: nil, Body: struct{}{}}
	}
}
