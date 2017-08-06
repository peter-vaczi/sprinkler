package core

import (
	"bytes"
	"encoding/json"
	"log"
)

type Data struct {
	Devices  *Devices  `json:"devices"`
	Programs *Programs `json:"programs"`
}

var data Data

func init() {
	data = Data{Devices: NewDevices(),
		Programs: NewPrograms(),
	}
}

func LoadState() {
}

func StoreState() {
	js, err := json.Marshal(data)
	if err == nil {
		log.Printf("json data: %v", bytes.NewBuffer(js))
	}
}

func Run(mainEvents chan interface{}) {
	for {
		select {
		case event := <-mainEvents:
			handleEvent(event)
		}
	}
}

func handleEvent(event interface{}) {

	switch event := event.(type) {
	case MsgDeviceList:
		event.ResponseChan <- MsgResponse{Error: nil, Body: data.Devices}
	case MsgDeviceAdd:
		err := data.Devices.Add(event.Device)
		event.ResponseChan <- MsgResponse{Error: err}
	case MsgDeviceGet:
		dev, err := data.Devices.Get(event.Name)
		event.ResponseChan <- MsgResponse{Error: err, Body: dev}
	case MsgDeviceDel:
		err := data.Devices.Del(event.Name)
		event.ResponseChan <- MsgResponse{Error: err}
	case MsgDeviceSet:
		err := data.Devices.Set(event.Name, event.Device)
		event.ResponseChan <- MsgResponse{Error: err}
	case MsgProgramList:
		event.ResponseChan <- MsgResponse{Error: nil, Body: data.Programs}
	case MsgProgramCreate:
		err := data.Programs.Add(event.Program)
		event.ResponseChan <- MsgResponse{Error: err}
	case MsgProgramGet:
		prg, err := data.Programs.Get(event.Name)
		event.ResponseChan <- MsgResponse{Error: err, Body: prg}
	case MsgProgramDel:
		err := data.Programs.Del(event.Name)
		event.ResponseChan <- MsgResponse{Error: err}
	case MsgProgramStart:
		prg, err := data.Programs.Get(event.Name)
		if err != nil {
			event.ResponseChan <- MsgResponse{Error: err}
			return
		}
		err = prg.Start()
		event.ResponseChan <- MsgResponse{Error: err}
	case MsgProgramStop:
		prg, err := data.Programs.Get(event.Name)
		if err != nil {
			event.ResponseChan <- MsgResponse{Error: err}
			return
		}
		err = prg.Stop()
		event.ResponseChan <- MsgResponse{Error: err}
	case MsgProgramAddDevice:
		prg, err := data.Programs.Get(event.Program)
		if err != nil {
			event.ResponseChan <- MsgResponse{Error: err}
			return
		}
		dev, err := data.Devices.Get(event.Device)
		if err != nil {
			event.ResponseChan <- MsgResponse{Error: err}
			return
		}
		err = prg.AddDevice(dev, event.Duration)
		event.ResponseChan <- MsgResponse{Error: err}
	case MsgProgramDelDevice:
		prg, err := data.Programs.Get(event.Program)
		if err != nil {
			event.ResponseChan <- MsgResponse{Error: err}
			return
		}
		err = prg.DelDevice(event.Idx)
		event.ResponseChan <- MsgResponse{Error: err}
	}
}
