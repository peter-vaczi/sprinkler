package core_test

import (
	"testing"

	"github.com/peter.vaczi/sprinklerd/core"
	"github.com/stretchr/testify/assert"
)

func TestDevice(t *testing.T) {
	gpioStub := NewGpioStub()
	core.InitGpio(gpioStub)
	dev := core.Device{Name: "dev1", Pin: 42}
	dev.Init()

	p, _ := gpioStub.pins[42]
	assert.NotNil(t, p)
	assert.Equal(t, 42, p.pin)
	assert.True(t, p.output)
	assert.False(t, p.high)

	dev.TurnOn()
	assert.True(t, p.high)

	dev.SetPin(12)
	p, _ = gpioStub.pins[12]
	assert.NotNil(t, p)
	assert.Equal(t, 12, p.pin)
	assert.False(t, p.high)

	dev.TurnOn()
	assert.True(t, p.high)
	assert.True(t, dev.IsOn())

	dev.TurnOff()
	assert.False(t, p.high)
	assert.False(t, dev.IsOn())

	dev.SetState(13, true)
	p, _ = gpioStub.pins[13]
	assert.NotNil(t, p)
	assert.True(t, p.high)
	assert.True(t, dev.IsOn())

	dev.SetOnIsLow(true)
	p, _ = gpioStub.pins[13]
	dev.TurnOff()
	assert.False(t, dev.IsOn())
	assert.True(t, p.high)

	dev.TurnOn()
	assert.True(t, dev.IsOn())
	assert.False(t, p.high)
}

func TestDevices(t *testing.T) {
	devs := core.NewDevices()
	if assert.NotNil(t, devs) {
		d1 := &core.Device{Name: "dev1"}
		d2 := &core.Device{Name: "dev2"}
		assert.Nil(t, devs.Add(d1))
		assert.NotEmpty(t, *devs)

		assert.Equal(t, core.AlreadyExists, devs.Add(d1))
		assert.Equal(t, 1, len(*devs))

		assert.Nil(t, devs.Add(d2))
		assert.Equal(t, 2, len(*devs))

		d, err := devs.Get("d")
		assert.Nil(t, d)
		assert.Equal(t, core.NotFound, err)

		d, err = devs.Get("dev1")
		assert.NotNil(t, d)
		assert.Equal(t, d1, d)
		assert.Nil(t, err)

		d.Pin = 42
		assert.Nil(t, devs.Set("dev1", d))
		d, _ = devs.Get("dev1")
		assert.Equal(t, 42, d.Pin)

		assert.Equal(t, core.NotFound, devs.Set("d", d))

		assert.Equal(t, core.NotFound, devs.Del("d"))
		assert.Nil(t, devs.Del("dev1"))
		assert.Nil(t, devs.Del("dev2"))

		assert.Empty(t, *devs)
	}
}
