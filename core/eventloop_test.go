package core_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/peter-vaczi/sprinklerd/core"
	"github.com/stretchr/testify/assert"
)

var responses = make(chan core.MsgResponse)

func sendReceive(t *testing.T, msg interface{}, err error) interface{} {
	events := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())

	go core.Run(ctx, events)
	var resp core.MsgResponse

	events <- msg
	resp = <-responses
	cancel()
	assert.Equal(t, err, resp.Error)
	return resp.Body
}

func addDevice(t *testing.T, name string, pin int, err error) *core.Device {
	dev := &core.Device{Name: name, Pin: pin}
	msg := core.MsgDeviceAdd{MsgRequest: core.MsgRequest{ResponseChan: responses}, Device: dev}
	body := sendReceive(t, msg, err)
	assert.Empty(t, body)
	return dev
}

func setDevice(t *testing.T, name string, dev *core.Device, err error) *core.Device {
	msg := core.MsgDeviceSet{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: name, Device: dev}
	body := sendReceive(t, msg, err)
	assert.Empty(t, body)
	return dev
}

func getDevice(t *testing.T, name string, err error) *core.Device {
	msg := core.MsgDeviceGet{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: name}
	body := sendReceive(t, msg, err)

	if err != nil {
		assert.Empty(t, body)
		return nil
	}

	if assert.NotEmpty(t, body) {
		if assert.IsType(t, &core.Device{}, body) {
			return body.(*core.Device)
		}
	}
	return nil
}

func delDevice(t *testing.T, name string, err error) {
	msg := core.MsgDeviceDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: name}
	body := sendReceive(t, msg, err)
	assert.Empty(t, body)
}

func listDevices(t *testing.T) *core.Devices {
	msg := core.MsgDeviceList{MsgRequest: core.MsgRequest{ResponseChan: responses}}
	body := sendReceive(t, msg, nil)

	if assert.NotEmpty(t, body) {
		if assert.IsType(t, &core.Devices{}, body) {
			return body.(*core.Devices)
		}
	}
	return nil
}

func addProgram(t *testing.T, name string, err error) *core.Program {
	pr := &core.Program{Name: name}
	msg := core.MsgProgramCreate{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: pr}
	body := sendReceive(t, msg, err)
	assert.Empty(t, body)
	return pr
}

func getProgram(t *testing.T, name string, err error) *core.Program {
	msg := core.MsgProgramGet{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: name}
	body := sendReceive(t, msg, err)

	if err != nil {
		assert.Empty(t, body)
		return nil
	}

	if assert.NotEmpty(t, body) {
		if assert.IsType(t, &core.Program{}, body) {
			return body.(*core.Program)
		}
	}
	return nil
}

func delProgram(t *testing.T, name string, err error) {
	msg := core.MsgProgramDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: name}
	body := sendReceive(t, msg, err)
	assert.Empty(t, body)
}

func listPrograms(t *testing.T) *core.Programs {
	msg := core.MsgProgramList{MsgRequest: core.MsgRequest{ResponseChan: responses}}
	body := sendReceive(t, msg, nil)

	if assert.NotEmpty(t, body) {
		if assert.IsType(t, &core.Programs{}, body) {
			return body.(*core.Programs)
		}
	}
	return nil
}

func prAddDev(t *testing.T, pr string, dev string, dur time.Duration, err error) {
	msg := core.MsgProgramAddDevice{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: pr, Device: dev, Duration: dur}
	sendReceive(t, msg, err)
}

func prDelDev(t *testing.T, pr string, idx int, err error) {
	msg := core.MsgProgramDelDevice{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: pr, Idx: idx}
	sendReceive(t, msg, err)
}

func startProgram(t *testing.T, pr string, err error) {
	msg := core.MsgProgramStart{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: pr}
	sendReceive(t, msg, err)
}

func stopProgram(t *testing.T, pr string, err error) {
	msg := core.MsgProgramStop{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: pr}
	sendReceive(t, msg, err)
}

func TestEventloopAddDelDevice(t *testing.T) {
	addDevice(t, "dev1", 1, nil)

	dev := getDevice(t, "dev1", nil)
	assert.Equal(t, 1, dev.Pin)

	dev.Pin = 42
	setDevice(t, "dev1", dev, nil)

	dev = getDevice(t, "dev1", nil)
	assert.Equal(t, 42, dev.Pin)

	getDevice(t, "unknown-device", core.NotFound)

	devs := listDevices(t)
	assert.Equal(t, 1, len(*devs))
	d, _ := devs.Get("dev1")
	assert.NotNil(t, d)

	delDevice(t, "dev1", nil)
	getDevice(t, "dev1", core.NotFound)
}

func TestEventloopAddDelProgram(t *testing.T) {
	addProgram(t, "pr1", nil)
	pr := getProgram(t, "pr1", nil)
	assert.Equal(t, "pr1", pr.Name)

	getProgram(t, "unknown-program", core.NotFound)

	prgs := listPrograms(t)
	assert.Equal(t, 1, len(*prgs))
	d, _ := prgs.Get("pr1")
	assert.NotNil(t, d)

	delProgram(t, "pr1", nil)
	getProgram(t, "pr1", core.NotFound)
}

func TestEventloopAddDelDeviceToProgram(t *testing.T) {
	d1 := addDevice(t, "dev1", 1, nil)
	d2 := addDevice(t, "dev2", 2, nil)

	addProgram(t, "pr1", nil)

	prAddDev(t, "pr1", "dev1", 5*time.Second, nil)
	prAddDev(t, "pr1", "dev2", 8*time.Second, nil)

	prAddDev(t, "pr1", "dev-whatever", 5*time.Second, core.NotFound)
	prAddDev(t, "pr-whatever", "dev1", 5*time.Second, core.NotFound)

	prgs := listPrograms(t)
	assert.Equal(t, 1, len(*prgs))
	p, _ := prgs.Get("pr1")
	if assert.NotNil(t, p) {
		assert.Equal(t, len(p.Elements), 2)
		assert.Equal(t, p.Elements[0].Device, d1)
		assert.Equal(t, p.Elements[0].Duration, 5*time.Second)
		assert.Equal(t, p.Elements[1].Device, d2)
		assert.Equal(t, p.Elements[1].Duration, 8*time.Second)
	}

	prDelDev(t, "pr1", 2, core.OutOfRange)

	prDelDev(t, "pr1", 0, nil)
	prDelDev(t, "pr1", 0, nil)
	prDelDev(t, "pr1", 0, core.OutOfRange)

	// cleanup
	delProgram(t, "pr1", nil)
	delDevice(t, "dev1", nil)
	delDevice(t, "dev2", nil)

	prDelDev(t, "pr1", 2, core.NotFound)
}

func TestEventloopDelDeviceInUse(t *testing.T) {
	addDevice(t, "dev1", 1, nil)

	addProgram(t, "pr1", nil)

	prAddDev(t, "pr1", "dev1", 5*time.Second, nil)

	delDevice(t, "dev1", core.DeviceInUse)

	prDelDev(t, "pr1", 0, nil)
	delProgram(t, "pr1", nil)
	delDevice(t, "dev1", nil)
}

func TestEventloopStartStopProgram(t *testing.T) {
	d1 := addDevice(t, "dev1", 1, nil)
	d2 := addDevice(t, "dev2", 2, nil)

	addProgram(t, "pr1", nil)

	prAddDev(t, "pr1", "dev1", 1*time.Second, nil)
	prAddDev(t, "pr1", "dev2", 1*time.Second, nil)

	// start pr
	startProgram(t, "pr-unknown", core.NotFound)
	startProgram(t, "pr1", nil)
	time.Sleep(100 * time.Millisecond)
	assert.True(t, d1.IsOn())
	assert.False(t, d2.IsOn())

	// stop pr
	stopProgram(t, "pr-unknown", core.NotFound)
	stopProgram(t, "pr1", nil)
	time.Sleep(100 * time.Millisecond)
	assert.False(t, d1.IsOn())
	assert.False(t, d2.IsOn())

	// cleanup
	delProgram(t, "pr1", nil)
	delDevice(t, "dev1", nil)
	delDevice(t, "dev2", nil)
}

func TestEventloopLoadStore(t *testing.T) {
	// file not found
	core.DataFile = "file-not-found.json"
	core.LoadState()

	devs := listDevices(t)
	assert.Equal(t, 0, len(*devs))

	// permission denied
	core.DataFile = "/etc/shadow"
	core.LoadState()

	devs = listDevices(t)
	assert.Equal(t, 0, len(*devs))

	// program refere to an no-existent device
	core.DataFile = "invalid-data1.json"
	core.LoadState()

	devs = listDevices(t)
	assert.Equal(t, 0, len(*devs))

	// missing closing brace
	core.DataFile = "invalid-data2.json"
	core.LoadState()

	devs = listDevices(t)
	assert.Equal(t, 0, len(*devs))

	// valid data
	core.DataFile = "data_test.json"
	core.LoadState()

	devs = listDevices(t)
	assert.Equal(t, 5, len(*devs))

	core.DataFile = "data_test2.json"
	core.StoreState()

	str1, _ := ioutil.ReadFile("data_test.json")
	str2, _ := ioutil.ReadFile("data_test2.json")

	os.Remove("data_test2.json")
	assert.Equal(t, str1, str2)
}
