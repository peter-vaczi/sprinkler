package core_test

import (
	"context"
	"testing"

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
