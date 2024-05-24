package machine

import (
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
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

// GetSelfCPUPercent 获取自身进程使用CPU百分比
func GetSelfCPUPercent() float64 {
	return GetPIDCPUPercent(int32(os.Getpid()))
}

// GetPIDCPUPercent 获取指定进程使用CPU百分比
func GetPIDCPUPercent(pid int32) float64 {
	p, err := process.NewProcess(pid)
	if err != nil {
		return invalidPercent
	}
	res, err := p.CPUPercent()
	if err != nil {
		return invalidPercent
	}
	return res
}

// GetSelfMemory 获取自身进程使用内存数据
func GetSelfMemory() uint64 {
	return GetPIDMemory(int32(os.Getpid()))
}

// GetPIDMemory 获取指定进程使用内存大小
func GetPIDMemory(pid int32) uint64 {
	p, err := process.NewProcess(pid)
	if err != nil {
		return invalidBytes
	}
	res, err := p.MemoryInfo()
	if err != nil {
		return invalidBytes
	}
	return res.RSS
}
