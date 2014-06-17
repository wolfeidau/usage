package usage

import (
	"log"
	"os"
	"time"

	sigar "github.com/cloudfoundry/gosigar"
)

type ProcessMonitor struct {
	lastPtimeUser  uint64
	lastPtimeSys   uint64
	lastPtimeTotal uint64
	lastPReqTime   uint64
}

type CpuUsage struct {
	User  float64
	Sys   float64
	Total float64
}

type MemoryUsage struct {
	Size       uint64
	Resident   uint64
	Share      uint64
	PageFaults uint64
}

// create a process monitor for this process
func CreateProcessMonitor() *ProcessMonitor {
	p := &ProcessMonitor{}

	pid := os.Getpid()

	curPtime := sigar.ProcTime{}
	err := curPtime.Get(pid)

	if err != nil {
		log.Fatalf("Issue loading process monitor - %v", err)
	}

	// seed cpu time
	p.lastPtimeUser = curPtime.User
	p.lastPtimeSys = curPtime.Sys
	p.lastPtimeTotal = curPtime.Total
	p.lastPReqTime = unixTimeMs()

	return p
}

// query the cpu usage, this is calculated for period between requests so if you
// poll ever second you will get % of cpu used per second.
func (p *ProcessMonitor) GetCpuUsage() *CpuUsage {

	pid := os.Getpid()

	curPtime := sigar.ProcTime{}

	err := curPtime.Get(pid)

	if err != nil {
		log.Fatalf("[Error] error retrieving process cpu info %v", err)
		return nil
	}

	timeDelta := unixTimeMs() - p.lastPReqTime

	return &CpuUsage{
		Sys:   calcTime(p.lastPtimeSys, curPtime.Sys, timeDelta),
		User:  calcTime(p.lastPtimeUser, curPtime.User, timeDelta),
		Total: calcTime(p.lastPtimeTotal, curPtime.Total, timeDelta),
	}
}

// query the memory usage of the current process
func (p *ProcessMonitor) GetMemoryUsage() *MemoryUsage {

	pid := os.Getpid()

	curMemory := sigar.ProcMem{}

	err := curMemory.Get(pid)

	if err != nil {
		log.Fatalf("[Error] error retrieving memory info %v", err)
		return nil
	}

	return &MemoryUsage{curMemory.Size, curMemory.Resident, curMemory.Share, curMemory.PageFaults}

}

// covers either zero activity or zero time between requests
func calcTime(usageLast uint64, usageCur uint64, timeDelta uint64) float64 {

	usageDelta := usageCur - usageLast

	if usageDelta == 0 || timeDelta == 0 {
		return 0
	} else {
		return 100 * (float64(timeDelta) / float64(usageDelta))
	}
}

func unixTimeMs() uint64 {
	return uint64(time.Now().UnixNano() % 1e6 / 1e3)
}
