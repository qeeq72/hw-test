package cpu

import (
	"fmt"
	"time"
)

type Average struct {
	User      float32
	Nice      float32
	System    float32
	Idle      float32
	Load      float32
	TimeStamp time.Time
}

func (a Average) String() string {
	return fmt.Sprintf("User: %.2f, Nice: %.2f, System: %.2f, Idle: %.2f, Load: %.3f - %s", a.User, a.Nice, a.System, a.Idle, a.Load, a.TimeStamp.Format(time.RFC3339Nano))
}

func (c *CPUStat) GetAverage(before, after *Snapshot) *Average {
	dUser := float32(after.User - before.User)
	dNice := float32(after.Nice - before.Nice)
	dSystem := float32(after.System - before.System)
	dIdle := float32(after.Idle - before.Idle)
	return &Average{
		User:      dUser / float32(c.period/time.Second),
		Nice:      dNice / float32(c.period/time.Second),
		System:    dSystem / float32(c.period/time.Second),
		Idle:      dIdle / float32(c.period/time.Second),
		Load:      100.0 * (float32(dUser+dNice+dSystem) / float32(dUser+dNice+dSystem+dIdle)),
		TimeStamp: time.Now(),
	}
}
