package main

import (
	"fmt"
	"time"

	"github.com/qeeq72/hw-test/prj_system_monitoring/internal/cpustat"
)

func main() {
	cpuStat := cpustat.NewCPUStat("/proc/stat", 5)

	for {
		select {
		case <-cpuStat.AverageDone:
			fmt.Printf("Load average: %v\n", cpuStat.LoadAverage)
			fmt.Printf("CPU modes average: %+v\n", cpuStat.ModeStatAverage)
		default:
			snapshot, err := cpuStat.MakeSnapshot()
			if err != nil {
				fmt.Print(err)
				return
			}
			fmt.Println(*snapshot)
			time.Sleep(5 * time.Second)
		}
	}
}
