package core_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/peter-vaczi/sprinkler/core"
	"github.com/stretchr/testify/assert"
)

func TestSchedules(t *testing.T) {
	scheds := core.NewSchedules()
	if assert.NotNil(t, scheds) {
		// 1:07 AM every second day
		s1 := &core.Schedule{Name: "sc1", Spec: "7 1 */2 * *"}
		// every moday noon
		s2 := &core.Schedule{Name: "sc2", Spec: "0 12 * * mon"}
		s3 := &core.Schedule{Name: "sc2", Spec: "* * * * *"}
		s4 := &core.Schedule{Name: "sc4", Spec: "* * * * * invalid spec"}

		assert.Nil(t, scheds.Add(s1))
		assert.NotEmpty(t, *scheds)

		assert.NotNil(t, scheds.Add(s1))
		assert.Equal(t, 1, len(*scheds))

		assert.Nil(t, scheds.Add(s2))
		assert.Equal(t, 2, len(*scheds))

		assert.NotNil(t, scheds.Add(s3))
		assert.Equal(t, 2, len(*scheds))

		assert.NotNil(t, scheds.Add(s4))
		assert.Equal(t, 2, len(*scheds))

		s, err := scheds.Get("s")
		assert.Nil(t, s)
		assert.NotNil(t, err)

		s, err = scheds.Get("sc1")
		assert.NotNil(t, s)
		assert.Equal(t, s1, s)
		assert.Nil(t, err)
		next := s.GetNext()
		assert.True(t, next.After(time.Now()))
		assert.True(t, next.Before(time.Now().Add(48*time.Hour)))

		s, err = scheds.Get("sc2")
		next = s.GetNext()
		fmt.Printf(next.String())
		assert.True(t, next.After(time.Now()))
		assert.True(t, next.Before(time.Now().Add(168*time.Hour)))

		s.SetSpec("1 * * * *")
		next = s.GetNext()
		assert.True(t, next.After(time.Now()))
		assert.True(t, next.Before(time.Now().Add(1*time.Hour)))

		assert.NotNil(t, s.SetSpec("* * 7 * * * * invalid-stuff"))
		next = s.GetNext()
		assert.True(t, next.After(time.Now()))
		assert.True(t, next.Before(time.Now().Add(1*time.Hour)))

		s2, _ = scheds.Get("sc2")
		assert.NotNil(t, scheds.Set("invalid-name", s2))
		assert.Nil(t, scheds.Set("sc1", s2))
		s1, _ = scheds.Get("sc1")
		assert.Equal(t, s1.Program, s2.Program)
		assert.Equal(t, s1.Spec, s2.Spec)
		assert.Equal(t, s1.GetNext(), s2.GetNext())

		assert.NotNil(t, scheds.Del("s"))
		assert.Nil(t, scheds.Del("sc1"))
		assert.Nil(t, scheds.Del("sc2"))

		assert.Empty(t, *scheds)
	}
}

func TestScheduleSetProgram(t *testing.T) {
	gpioStub := NewGpioStub()
	core.InitGpio(gpioStub)
	d1 := &core.Device{Name: "dev1"}
	d2 := &core.Device{Name: "dev2"}
	p := &core.Program{Name: "pr1"}
	s := &core.Schedule{Name: "sc1"}

	p.AddDevice(d1, 1*time.Second)
	p.AddDevice(d2, 1*time.Second)
	s.SetProgram(p)
	assert.NotNil(t, s.Program)
	p.DelDevice(0)
	p.DelDevice(0)
}
