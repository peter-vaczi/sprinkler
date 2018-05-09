package core

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Data struct {
	Devices   *Devices   `json:"devices"`
	Programs  *Programs  `json:"programs"`
	Schedules *Schedules `json:"schedules"`
}

var data Data

var DataFile = "/var/lib/sprinkler.data"

func init() {
	data = Data{Devices: NewDevices(),
		Programs:  NewPrograms(),
		Schedules: NewSchedules(),
	}
}

func LoadState() {
	file, err := os.Open(DataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Printf("failed to open data file: %v", err)
		return
	}

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		log.Printf("failed to parse data file: %v", err)
		return
	}

	// re-initialize the gpio members
	for _, dev := range *data.Devices {
		dev.SetState(dev.Pin, dev.On)
	}

	// re-initialize the device pointers
	for _, pr := range *data.Programs {
		for _, elem := range pr.Elements {
			elem.Device, err = data.Devices.Get(elem.DeviceName)
			if err != nil {
				log.Printf("invalid data file, device %s not found", elem.DeviceName)
				data = Data{Devices: NewDevices(), Programs: NewPrograms(), Schedules: NewSchedules()}
				return
			}
		}
	}

	// re-initialize the program pointers
	for _, sc := range *data.Schedules {
		sc.Program, err = data.Programs.Get(sc.ProgramName)
		if err != nil {
			log.Printf("invalid data file, program %s not found", sc.ProgramName)
			data = Data{Devices: NewDevices(), Programs: NewPrograms(), Schedules: NewSchedules()}
			return
		}
	}
}

func StoreState() {
	js, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to convert data to json: %v", err)
	}

	err = ioutil.WriteFile(DataFile, js, 0744)
}

func Run(ctx context.Context, mainEvents chan interface{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case event := <-mainEvents:
			handleEvent(event)
		case <-sigChan:
			return
		case <-ctx.Done():
			return
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
		for _, pr := range *data.Programs {
			for _, e := range pr.Elements {
				if e.DeviceName == event.Name {
					event.ResponseChan <- MsgResponse{Error: DeviceInUse}
					return
				}
			}
		}
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
		prg.Start()
		event.ResponseChan <- MsgResponse{Error: err}
	case MsgProgramStop:
		prg, err := data.Programs.Get(event.Name)
		if err != nil {
			event.ResponseChan <- MsgResponse{Error: err}
			return
		}
		prg.Stop()
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
	case MsgScheduleList:
		event.ResponseChan <- MsgResponse{Error: nil, Body: data.Schedules}
	case MsgScheduleCreate:
		err := data.Schedules.Add(event.Schedule)
		event.ResponseChan <- MsgResponse{Error: err}
	case MsgScheduleGet:
		sch, err := data.Schedules.Get(event.Name)
		event.ResponseChan <- MsgResponse{Error: err, Body: sch}
	case MsgScheduleDel:
		err := data.Schedules.Del(event.Name)
		event.ResponseChan <- MsgResponse{Error: err}
	case MsgScheduleSet:
		if 0 < len(event.Schedule.ProgramName) {
			prg, err := data.Programs.Get(event.Schedule.ProgramName)
			if err != nil {
				event.ResponseChan <- MsgResponse{Error: err}
				return
			}
			event.Schedule.Program = prg
		}
		err := data.Schedules.Set(event.Name, event.Schedule)
		event.ResponseChan <- MsgResponse{Error: err}
	}
}
