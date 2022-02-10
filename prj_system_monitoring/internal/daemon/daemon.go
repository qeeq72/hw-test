package daemon

import (
	"context"
	"fmt"
	"sync"

	"github.com/qeeq72/hw-test/prj_system_monitoring/internal/stat/interfaces"
)

type Daemon struct {
	collectors []interfaces.ICollector
	wg         sync.WaitGroup
}

func NewDaemon(collectors []interfaces.ICollector) *Daemon {
	return &Daemon{
		collectors: collectors,
	}
}

func (d *Daemon) Run(ctx context.Context) {
	out := make(chan fmt.Stringer, len(d.collectors))

	for i := range d.collectors {
		d.wg.Add(1)
		go func(collector interfaces.ICollector) {
			defer d.wg.Done()
			collector.Run(ctx, out)
		}(d.collectors[i])
	}

	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case result := <-out:
				fmt.Println(result)
			}
		}
	}()

	d.wg.Wait()
}
