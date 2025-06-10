package monitor

import (
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskService struct {
	Usage  uint64 // percent used
	Status int32  // 1 if above 50%, 0 otherwise
	Path   string
}

// NewDiskService creates a new DiskService for the given path.
func NewDiskService(path string) *DiskService {
	return &DiskService{Path: path}
}

// Start begins monitoring disk usage in a background goroutine.
func (d *DiskService) Start() {
	go func() {
		for {

			vm, err := disk.Usage(d.Path)
			if err != nil {
				atomic.StoreUint64(&d.Usage, 0)
				atomic.StoreInt32(&d.Status, 0)
				time.Sleep(2 * time.Second)
				continue
			}

			percentUsed := vm.UsedPercent

			atomic.StoreUint64(&d.Usage, uint64(percentUsed))
			if percentUsed > 50 {
				atomic.StoreInt32(&d.Status, 1)
			} else {
				atomic.StoreInt32(&d.Status, 0)
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

// GetUsage returns the current disk usage percentage.
func (d *DiskService) GetUsage() uint64 {
	return atomic.LoadUint64(&d.Usage)
}

// GetStatus returns 1 if usage is above 50%, else 0.
func (d *DiskService) GetStatus() int32 {
	return atomic.LoadInt32(&d.Status)
}
