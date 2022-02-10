package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/qeeq72/hw-test/prj_system_monitoring/internal/daemon"
	"github.com/qeeq72/hw-test/prj_system_monitoring/internal/stat/cpu"
	"github.com/qeeq72/hw-test/prj_system_monitoring/internal/stat/interfaces"
)

func main() {
	cpuStat := cpu.NewCPUStat("/proc/stat", 3, 10)
	col := []interfaces.ICollector{cpuStat}

	ctx, cancel := context.WithCancel(context.Background())
	d := daemon.NewDaemon(col)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		d.Run(ctx)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	cancel()
	wg.Wait()
}
