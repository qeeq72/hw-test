package cpu

import (
	"context"
	"fmt"
	"time"
)

type CPUStat struct {
	path   string
	buffer chan *Snapshot
	period time.Duration
	delay  time.Duration
}

func NewCPUStat(path string, depth, rate int) *CPUStat {
	return &CPUStat{
		path:   path,
		buffer: make(chan *Snapshot, depth/rate+1),
		period: time.Duration(rate) * time.Second,
		delay:  time.Duration(depth) * time.Second,
	}
}

func (c *CPUStat) Run(ctx context.Context, out chan fmt.Stringer) {
	snapshot, err := c.Collect()
	if err != nil {
		return
	}
	c.buffer <- snapshot

	timer := time.NewTimer(c.delay)
	tickerPush := time.NewTicker(c.period)
	tickerPull := time.NewTicker(999 * time.Hour)
	defer timer.Stop()
	defer tickerPush.Stop()
	defer tickerPull.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			tickerPull.Reset(c.period)
			snapBefore := <-c.buffer
			snapAfter, err := c.Collect()
			if err != nil {
				return
			}
			out <- c.GetAverage(snapBefore, snapAfter)
		case <-tickerPush.C:
			snapshot, err := c.Collect()
			if err != nil {
				return
			}
			c.buffer <- snapshot
		case <-tickerPull.C:
			snapBefore := <-c.buffer
			snapAfter, err := c.Collect()
			if err != nil {
				return
			}
			out <- c.GetAverage(snapBefore, snapAfter)
		}
	}
}
