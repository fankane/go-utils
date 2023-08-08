package machine

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

const (
	invalidPercent = -1.0
	invalidBytes   = 0
)

// GetCPUPercent 1秒钟时间段，CPU使用百分比，多个CPU一起计算
func GetCPUPercent() float64 {
	cpus, err := cpu.Percent(time.Second, false)
	if err != nil || len(cpus) == 0 {
		return invalidPercent
	}
	return cpus[0]
}

func GetMemPercent() float64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return invalidPercent
	}
	return memInfo.UsedPercent
}

func GetMemUsed() uint64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return invalidBytes
	}
	return memInfo.Used
}

func GetMemAvailable() uint64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return invalidBytes
	}
	return memInfo.Available
}

func GetMemTotal() uint64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return invalidBytes
	}
	return memInfo.Total
}

func GetDiskUsedPercent() float64 {
	parts, err := disk.Partitions(true)
	if err != nil || len(parts) == 0 {
		return invalidPercent
	}
	diskInfo, err := disk.Usage(parts[0].Mountpoint)
	if err != nil || diskInfo == nil {
		return invalidPercent
	}
	return diskInfo.UsedPercent
}
