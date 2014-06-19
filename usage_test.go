package usage

import (
	"testing"
	"time"

	"log"

	. "launchpad.net/gocheck"
)

func Test(t *testing.T) {
	TestingT(t)
}

type UsageSuite struct {
	processMonitor *ProcessMonitor
}

var _ = Suite(&UsageSuite{})

func (us *UsageSuite) SetUpTest(c *C) {
	us.processMonitor = CreateProcessMonitor()
}

func (us *UsageSuite) TestCpuUsage(c *C) {
	time.Sleep(1 * time.Second)
	cpuUsage := us.processMonitor.GetCpuUsage()
	c.Assert(cpuUsage.Total, Equals, 0.0)
}

func (us *UsageSuite) TestMemoryUsage(c *C) {
	memUsage := us.processMonitor.GetMemoryUsage()
	c.Assert(memUsage.Size, Not(Equals), 0)
}

func (us *UsageSuite) TestUnixTimeMs(c *C) {
	unixTimeMs := UnixTimeMs()
	log.Printf("%v", unixTimeMs)
	c.Assert(unixTimeMs, Not(Equals), 0)
}
