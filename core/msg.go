package core

import "time"

type MsgResponse struct {
	Error error
	Body  interface{}
}

type MsgRequest struct {
	ResponseChan chan MsgResponse
}

type MsgDeviceList struct {
	MsgRequest
}

type MsgDeviceAdd struct {
	MsgRequest
	Device *Device
}

type MsgDeviceGet struct {
	MsgRequest
	Name string
}

type MsgDeviceDel struct {
	MsgRequest
	Name string
}

type MsgDeviceSet struct {
	MsgRequest
	Name   string
	Device *Device
}

type MsgProgramList struct {
	MsgRequest
}

type MsgProgramCreate struct {
	MsgRequest
	Program *Program
}

type MsgProgramGet struct {
	MsgRequest
	Name string
}

type MsgProgramDel struct {
	MsgRequest
	Name string
}

type MsgProgramStart struct {
	MsgRequest
	Name string
}

type MsgProgramStop struct {
	MsgRequest
	Name string
}

type MsgProgramAddDevice struct {
	MsgRequest
	Program  string
	Device   string
	Duration time.Duration
}

type MsgProgramDelDevice struct {
	MsgRequest
	Program string
	Idx     int
}

type MsgScheduleList struct {
	MsgRequest
}

type MsgScheduleCreate struct {
	MsgRequest
	Schedule *Schedule
}

type MsgScheduleGet struct {
	MsgRequest
	Name string
}

type MsgScheduleDel struct {
	MsgRequest
	Name string
}

type MsgScheduleSet struct {
	MsgRequest
	Name     string
	Schedule *Schedule
}
