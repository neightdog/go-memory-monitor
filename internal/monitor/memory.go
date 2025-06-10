package monitor

import (
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryService struct {
	Usage  uint64 // percent used
	Status int32  // 1 if above 50%, 0 otherwise
}

// NewMemoryService creates a new MemoryService.
func NewMemoryService() *MemoryService {
	return &MemoryService{}
}

// Start begins monitoring memory usage in a background goroutine.
func (m *MemoryService) Start() {
	go func() {
		for {
			vm, err := mem.VirtualMemory()
			if err != nil {
				atomic.StoreUint64(&m.Usage, 0)
				atomic.StoreInt32(&m.Status, 0)
				time.Sleep(2 * time.Second)
				continue
			}
			percentageUsed := vm.UsedPercent

			atomic.StoreUint64(&m.Usage, uint64(percentageUsed))
			if percentageUsed > 50 {
				atomic.StoreInt32(&m.Status, 1)
			} else {
				atomic.StoreInt32(&m.Status, 0)
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

// GetUsage returns the current memory usage percentage.
func (m *MemoryService) GetUsage() uint64 {
	return atomic.LoadUint64(&m.Usage)
}

// GetStatus returns 1 if usage is above 50%, else 0.
func (m *MemoryService) GetStatus() int32 {
	return atomic.LoadInt32(&m.Status)
}
