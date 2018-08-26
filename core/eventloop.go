package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Data struct {
	Devices   *Devices   `json:"devices"`
	Programs  *Programs  `json:"programs"`
	Schedules *Schedules `json:"schedules"`
}

func NewData() *Data {
	return &Data{Devices: NewDevices(), Programs: NewPrograms(), Schedules: NewSchedules()}
}

var DataFile = "/var/lib/sprinkler.data"

func LoadState() *Data {
	file, err := os.Open(DataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return NewData()
		}
		log.Printf("failed to open data file: %v", err)
		return nil
	}

	data := NewData()
	err = json.NewDecoder(file).Decode(data)
	if err != nil {
		log.Printf("failed to parse data file: %v", err)
		return nil
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
				return nil
			}
		}
	}

	// re-initialize the program pointers
	for _, sc := range *data.Schedules {
		sc.Program, err = data.Programs.Get(sc.ProgramName)
		if err != nil {
			log.Printf("invalid data file, program %s not found", sc.ProgramName)
			return nil
		}
	}
	return data
}

func (d *Data) StoreState() {
	js, err := json.Marshal(d)
	if err != nil {
		log.Printf("failed to convert data to json: %v", err)
	}

	err = ioutil.WriteFile(DataFile, js, 0744)
}
