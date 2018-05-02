package core_test

import (
	"testing"
	"time"

	"github.com/peter-vaczi/sprinkler/core"
	"github.com/stretchr/testify/assert"
)

func TestProgramAddDelDevice(t *testing.T) {
	gpioStub := NewGpioStub()
	core.InitGpio(gpioStub)
	d1 := &core.Device{Name: "dev1"}
	d2 := &core.Device{Name: "dev2"}
	p := &core.Program{Name: "pr1"}

	assert.Empty(t, p.Elements)

	assert.Nil(t, p.AddDevice(d1, 1*time.Second))
	assert.NotEmpty(t, p.Elements)

	assert.Nil(t, p.AddDevice(d2, 1*time.Second))
	assert.Equal(t, 2, len(p.Elements))

	assert.Equal(t, core.OutOfRange, p.DelDevice(2))
	assert.Nil(t, p.DelDevice(0))
	assert.Nil(t, p.DelDevice(0))

	assert.Empty(t, p.Elements)
}

func TestProgramStartStop(t *testing.T) {
	gpioStub := NewGpioStub()
	core.InitGpio(gpioStub)
	d1 := &core.Device{Name: "dev1", Pin: 1}
	d2 := &core.Device{Name: "dev2", Pin: 2}
	d1.Init()
	d2.Init()

	p := &core.Program{Name: "pr1"}
	assert.Nil(t, p.AddDevice(d1, 1*time.Second))
	assert.Nil(t, p.AddDevice(d2, 1*time.Second))

	p.Start()
	time.Sleep(500 * time.Millisecond)
	assert.True(t, p.Elements[0].Device.IsOn())
	time.Sleep(1600 * time.Millisecond)
	assert.False(t, p.Elements[0].Device.IsOn())
	assert.True(t, p.Elements[1].Device.IsOn())

	p.Stop()
	assert.False(t, p.Elements[0].Device.IsOn())
	assert.False(t, p.Elements[0].Device.IsOn())

	assert.Nil(t, p.DelDevice(1))
	p.Start()
	time.Sleep(2100 * time.Millisecond)
	assert.False(t, p.Elements[0].Device.IsOn())
}

func TestPrograms(t *testing.T) {
	progs := core.NewPrograms()
	if assert.NotNil(t, progs) {
		p1 := &core.Program{Name: "pr1"}
		p2 := &core.Program{Name: "pr2"}
		assert.Nil(t, progs.Add(p1))
		assert.NotEmpty(t, *progs)

		assert.NotNil(t, progs.Add(p1))
		assert.Equal(t, 1, len(*progs))

		assert.Nil(t, progs.Add(p2))
		assert.Equal(t, 2, len(*progs))

		p, err := progs.Get("p")
		assert.Nil(t, p)
		assert.NotNil(t, err)

		p, err = progs.Get("pr1")
		assert.NotNil(t, p)
		assert.Equal(t, p1, p)
		assert.Nil(t, err)

		assert.NotNil(t, progs.Del("p"))
		assert.Nil(t, progs.Del("pr1"))
		assert.Nil(t, progs.Del("pr2"))

		assert.Empty(t, *progs)
	}
}
