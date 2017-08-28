package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/peter.vaczi/sprinklerd/core"
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

func TestEventloopAddDelDevice(t *testing.T) {
	var msg interface{}
	var body interface{}

	// add one device
	dev := &core.Device{Name: "dev1", Pin: 1}
	msg = core.MsgDeviceAdd{MsgRequest: core.MsgRequest{ResponseChan: responses}, Device: dev}
	body = sendReceive(t, msg, nil)
	assert.Empty(t, body)

	// read device dev1
	msg = core.MsgDeviceGet{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "dev1"}
	body = sendReceive(t, msg, nil)
	if assert.NotEmpty(t, body) {
		if assert.IsType(t, &core.Device{}, body) {
			assert.Equal(t, 1, body.(*core.Device).Pin)
		}
	}

	// read a non existing device
	msg = core.MsgDeviceGet{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "unknown-device"}
	body = sendReceive(t, msg, core.NotFound)
	assert.Empty(t, body)

	// read all devices (list of one)
	msg = core.MsgDeviceList{MsgRequest: core.MsgRequest{ResponseChan: responses}}
	body = sendReceive(t, msg, nil)
	if assert.NotEmpty(t, body) {
		if assert.IsType(t, &core.Devices{}, body) {
			d, _ := body.(*core.Devices).Get("dev1")
			assert.NotNil(t, d)
		}
	}

	// delete dev1
	msg = core.MsgDeviceDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "dev1"}
	body = sendReceive(t, msg, nil)
	assert.Empty(t, body)

	// read dev1 (should fail, since we deleted it)
	msg = core.MsgDeviceGet{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "dev1"}
	body = sendReceive(t, msg, core.NotFound)
	assert.Empty(t, body)
}

func TestEventloopAddDelProgram(t *testing.T) {
	var msg interface{}
	var body interface{}

	// create one program
	pr := &core.Program{Name: "pr1"}
	msg = core.MsgProgramCreate{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: pr}
	body = sendReceive(t, msg, nil)
	assert.Empty(t, body)

	// read program pr1
	msg = core.MsgProgramGet{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "pr1"}
	body = sendReceive(t, msg, nil)
	if assert.NotEmpty(t, body) {
		if assert.IsType(t, &core.Program{}, body) {
			assert.Equal(t, "pr1", body.(*core.Program).Name)
		}
	}

	// read a non existing program
	msg = core.MsgProgramGet{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "unknown-program"}
	body = sendReceive(t, msg, core.NotFound)
	assert.Empty(t, body)

	// read all prorams (list of one)
	msg = core.MsgProgramList{MsgRequest: core.MsgRequest{ResponseChan: responses}}
	body = sendReceive(t, msg, nil)
	if assert.NotEmpty(t, body) {
		if assert.IsType(t, &core.Programs{}, body) {
			d, _ := body.(*core.Programs).Get("pr1")
			assert.NotNil(t, d)
		}
	}

	// delete pr1
	msg = core.MsgProgramDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "pr1"}
	body = sendReceive(t, msg, nil)
	assert.Empty(t, body)

	// read pr1 (should fail, since we deleted it)
	msg = core.MsgProgramGet{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "pr1"}
	body = sendReceive(t, msg, core.NotFound)
	assert.Empty(t, body)
}

func TestEventloopAddDelDeviceToProgram(t *testing.T) {
	var msg interface{}

	// add one device
	d1 := &core.Device{Name: "dev1", Pin: 1}
	msg = core.MsgDeviceAdd{MsgRequest: core.MsgRequest{ResponseChan: responses}, Device: d1}
	sendReceive(t, msg, nil)

	d2 := &core.Device{Name: "dev2", Pin: 2}
	msg = core.MsgDeviceAdd{MsgRequest: core.MsgRequest{ResponseChan: responses}, Device: d2}
	sendReceive(t, msg, nil)

	// create one program
	pr := &core.Program{Name: "pr1"}
	msg = core.MsgProgramCreate{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: pr}
	sendReceive(t, msg, nil)

	// add dev1 to pr
	msg = core.MsgProgramAddDevice{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: "pr1", Device: "dev1", Duration: 5 * time.Second}
	sendReceive(t, msg, nil)

	// add dev2 to pr
	msg = core.MsgProgramAddDevice{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: "pr1", Device: "dev2", Duration: 8 * time.Second}
	sendReceive(t, msg, nil)

	// add unknown dev to pr
	msg = core.MsgProgramAddDevice{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: "pr1", Device: "dev-whatever", Duration: 5 * time.Second}
	sendReceive(t, msg, core.NotFound)

	// add dev1 to unknown pr
	msg = core.MsgProgramAddDevice{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: "pr-whatever", Device: "dev1", Duration: 5 * time.Second}
	sendReceive(t, msg, core.NotFound)

	msg = core.MsgProgramList{MsgRequest: core.MsgRequest{ResponseChan: responses}}
	body := sendReceive(t, msg, nil)
	if assert.NotEmpty(t, body) {
		if assert.IsType(t, &core.Programs{}, body) {
			p, _ := body.(*core.Programs).Get("pr1")
			if assert.NotNil(t, p) {
				assert.Equal(t, len(p.Elements), 2)
				assert.Equal(t, p.Elements[0].Device, d1)
				assert.Equal(t, p.Elements[0].Duration, 5*time.Second)
				assert.Equal(t, p.Elements[1].Device, d2)
				assert.Equal(t, p.Elements[1].Duration, 8*time.Second)
			}
		}
	}

	// delete unknown index
	msg = core.MsgProgramDelDevice{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: "pr1", Idx: 2}
	sendReceive(t, msg, core.OutOfRange)

	// delete all
	msg = core.MsgProgramDelDevice{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: "pr1", Idx: 0}
	sendReceive(t, msg, nil)
	sendReceive(t, msg, nil)
	sendReceive(t, msg, core.OutOfRange)

	// cleanup
	msg = core.MsgProgramDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "pr1"}
	sendReceive(t, msg, nil)
	msg = core.MsgDeviceDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "dev1"}
	sendReceive(t, msg, nil)
	msg = core.MsgDeviceDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "dev2"}
	sendReceive(t, msg, nil)
}

func TestEventloopDelDeviceInUse(t *testing.T) {
	var msg interface{}

	// add one device
	d1 := &core.Device{Name: "dev1", Pin: 1}
	msg = core.MsgDeviceAdd{MsgRequest: core.MsgRequest{ResponseChan: responses}, Device: d1}
	sendReceive(t, msg, nil)

	// create one program
	pr := &core.Program{Name: "pr1"}
	msg = core.MsgProgramCreate{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: pr}
	sendReceive(t, msg, nil)

	// add dev1 to pr
	msg = core.MsgProgramAddDevice{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: "pr1", Device: "dev1", Duration: 5 * time.Second}
	sendReceive(t, msg, nil)

	// delete dev1
	msg = core.MsgDeviceDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "dev1"}
	sendReceive(t, msg, core.DeviceInUse)

	// delete all
	msg = core.MsgProgramDelDevice{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: "pr1", Idx: 0}
	sendReceive(t, msg, nil)
	msg = core.MsgProgramDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "pr1"}
	sendReceive(t, msg, nil)

	// delete dev1
	msg = core.MsgDeviceDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "dev1"}
	sendReceive(t, msg, nil)
}

func TestEventloopStartStopProgram(t *testing.T) {
	var msg interface{}

	// add one device
	d1 := &core.Device{Name: "dev1", Pin: 1}
	msg = core.MsgDeviceAdd{MsgRequest: core.MsgRequest{ResponseChan: responses}, Device: d1}
	sendReceive(t, msg, nil)

	d2 := &core.Device{Name: "dev2", Pin: 2}
	msg = core.MsgDeviceAdd{MsgRequest: core.MsgRequest{ResponseChan: responses}, Device: d2}
	sendReceive(t, msg, nil)

	// create one program
	pr := &core.Program{Name: "pr1"}
	msg = core.MsgProgramCreate{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: pr}
	sendReceive(t, msg, nil)

	// add dev1 to pr
	msg = core.MsgProgramAddDevice{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: "pr1", Device: "dev1", Duration: 1 * time.Second}
	sendReceive(t, msg, nil)

	// add dev2 to pr
	msg = core.MsgProgramAddDevice{MsgRequest: core.MsgRequest{ResponseChan: responses}, Program: "pr1", Device: "dev2", Duration: 1 * time.Second}
	sendReceive(t, msg, nil)

	// start pr
	msg = core.MsgProgramStart{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "pr-unknown"}
	sendReceive(t, msg, core.NotFound)
	msg = core.MsgProgramStart{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "pr1"}
	sendReceive(t, msg, nil)
	time.Sleep(100 * time.Millisecond)
	assert.True(t, d1.IsOn())
	assert.False(t, d2.IsOn())

	// stop pr
	msg = core.MsgProgramStop{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "pr-unknown"}
	sendReceive(t, msg, core.NotFound)
	msg = core.MsgProgramStop{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "pr1"}
	sendReceive(t, msg, nil)
	time.Sleep(100 * time.Millisecond)
	assert.False(t, d1.IsOn())
	assert.False(t, d2.IsOn())

	// cleanup
	msg = core.MsgProgramDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "pr1"}
	sendReceive(t, msg, nil)
	msg = core.MsgDeviceDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "dev1"}
	sendReceive(t, msg, nil)
	msg = core.MsgDeviceDel{MsgRequest: core.MsgRequest{ResponseChan: responses}, Name: "dev2"}
	sendReceive(t, msg, nil)
}
