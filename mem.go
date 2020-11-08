package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

type VirtualMemoryStat struct {
	// Total amount of RAM on this system
	Total uint64 `json:"total"`

	// RAM available for programs to allocate
	//
	// This value is computed from the kernel specific values.
	Available uint64 `json:"available"`

	// RAM used by programs
	//
	// This value is computed from the kernel specific values.
	Used uint64 `json:"used"`

	// Percentage of RAM used by programs
	//
	// This value is computed from the kernel specific values.
	UsedPercent float64 `json:"usedPercent"`

	// This is the kernel's notion of free memory; RAM chips whose bits nobody
	// cares about the value of right now. For a human consumable number,
	// Available is what you really want.
	Free uint64 `json:"free"`

	// OS X / BSD specific numbers:
	// http://www.macyourself.com/2010/02/17/what-is-free-wired-active-and-inactive-system-memory-ram/
	Active   uint64 `json:"active"`
	Inactive uint64 `json:"inactive"`
	Wired    uint64 `json:"wired"`

	// FreeBSD specific numbers:
	// https://reviews.freebsd.org/D8467
	Laundry uint64 `json:"laundry"`

	// Linux specific numbers
	// https://www.centos.org/docs/5/html/5.1/Deployment_Guide/s2-proc-meminfo.html
	// https://www.kernel.org/doc/Documentation/filesystems/proc.txt
	// https://www.kernel.org/doc/Documentation/vm/overcommit-accounting
	Buffers        uint64 `json:"buffers"`
	Cached         uint64 `json:"cached"`
	WriteBack      uint64 `json:"writeBack"`
	Dirty          uint64 `json:"dirty"`
	WriteBackTmp   uint64 `json:"writeBackTmp"`
	Shared         uint64 `json:"shared"`
	Slab           uint64 `json:"slab"`
	Sreclaimable   uint64 `json:"sreclaimable"`
	Sunreclaim     uint64 `json:"sunreclaim"`
	PageTables     uint64 `json:"pageTables"`
	SwapCached     uint64 `json:"swapCached"`
	CommitLimit    uint64 `json:"commitLimit"`
	CommittedAS    uint64 `json:"committedAS"`
	HighTotal      uint64 `json:"highTotal"`
	HighFree       uint64 `json:"highFree"`
	LowTotal       uint64 `json:"lowTotal"`
	LowFree        uint64 `json:"lowFree"`
	SwapTotal      uint64 `json:"swapTotal"`
	SwapFree       uint64 `json:"swapFree"`
	Mapped         uint64 `json:"mapped"`
	VmallocTotal   uint64 `json:"vmallocTotal"`
	VmallocUsed    uint64 `json:"vmallocUsed"`
	VmallocChunk   uint64 `json:"vmallocChunk"`
	HugePagesTotal uint64 `json:"hugePagesTotal"`
	HugePagesFree  uint64 `json:"hugePagesFree"`
	HugePageSize   uint64 `json:"hugePageSize"`
}

func (m VirtualMemoryStat) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

func (ps *PSUtils)VirtualMemory() (*VirtualMemoryStat, error) {
	filename := "/proc/meminfo"
	lines, _ := ps.ReadLines(filename)

	// flag if MemAvailable is in /proc/meminfo (kernel 3.14+)
	memavail := false

	ret := &VirtualMemoryStat{}

	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		value = strings.Replace(value, " kB", "", -1)

		t, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return ret, err
		}
		switch key {
		case "MemTotal":
			ret.Total = t * 1024
		case "MemFree":
			ret.Free = t * 1024
		case "MemAvailable":
			memavail = true
			ret.Available = t * 1024
		case "Buffers":
			ret.Buffers = t * 1024
		case "Cached":
			ret.Cached = t * 1024
		case "Active":
			ret.Active = t * 1024
		case "Inactive":
			ret.Inactive = t * 1024
		case "WriteBack":
			ret.WriteBack = t * 1024
		case "WriteBackTmp":
			ret.WriteBackTmp = t * 1024
		case "Dirty":
			ret.Dirty = t * 1024
		case "Shmem":
			ret.Shared = t * 1024
		case "Slab":
			ret.Slab = t * 1024
		case "Sreclaimable":
			ret.Sreclaimable = t * 1024
		case "Sunreclaim":
			ret.Sunreclaim = t * 1024
		case "PageTables":
			ret.PageTables = t * 1024
		case "SwapCached":
			ret.SwapCached = t * 1024
		case "CommitLimit":
			ret.CommitLimit = t * 1024
		case "Committed_AS":
			ret.CommittedAS = t * 1024
		case "HighTotal":
			ret.HighTotal = t * 1024
		case "HighFree":
			ret.HighFree = t * 1024
		case "LowTotal":
			ret.LowTotal = t * 1024
		case "LowFree":
			ret.LowFree = t * 1024
		case "SwapTotal":
			ret.SwapTotal = t * 1024
		case "SwapFree":
			ret.SwapFree = t * 1024
		case "Mapped":
			ret.Mapped = t * 1024
		case "VmallocTotal":
			ret.VmallocTotal = t * 1024
		case "VmallocUsed":
			ret.VmallocUsed = t * 1024
		case "VmallocChunk":
			ret.VmallocChunk = t * 1024
		case "HugePages_Total":
			ret.HugePagesTotal = t
		case "HugePages_Free":
			ret.HugePagesFree = t
		case "Hugepagesize":
			ret.HugePageSize = t * 1024
		}
	}

	ret.Cached += ret.Sreclaimable

	if !memavail {
		ret.Available = ret.Cached + ret.Free
	}

	ret.Used = ret.Total - ret.Free - ret.Buffers - ret.Cached
	ret.UsedPercent = float64(ret.Used) / float64(ret.Total) * 100.0

	return ret, nil
}

