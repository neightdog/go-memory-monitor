package monitor

import (
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUService struct {
	Usage  uint64 // percent used
	Status int32  // 1 if above 50%, 0 otherwise
}

// NewCPUService creates a new CPUService.
func NewCPUService() *CPUService {
	return &CPUService{}
}

// Start begins monitoring CPU usage in a background goroutine.
func (c *CPUService) Start() {
	go func() {
		for {
			percentages, err := cpu.Percent(0, false)
			if err != nil || len(percentages) == 0 {
				atomic.StoreUint64(&c.Usage, 0)
				atomic.StoreInt32(&c.Status, 0)
				time.Sleep(2 * time.Second)
				continue
			}
			usage := percentages[0]
			atomic.StoreUint64(&c.Usage, uint64(usage))
			if usage > 50 {
				atomic.StoreInt32(&c.Status, 1)
			} else {
				atomic.StoreInt32(&c.Status, 0)
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

// GetUsage returns the current CPU usage percentage.
func (c *CPUService) GetUsage() uint64 {
	return atomic.LoadUint64(&c.Usage)
}

// GetStatus returns 1 if usage is above 50%, else 0.
func (c *CPUService) GetStatus() int32 {
	return atomic.LoadInt32(&c.Status)
}
