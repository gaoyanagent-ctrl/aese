package simtime

import (
	"testing"
	"time"
)

func TestClockFailsClosedOnBackwardTime(t *testing.T) {
	start := time.Date(2026, 7, 8, 10, 0, 0, 0, time.FixedZone("CST", 8*3600))
	c, _ := New(start)
	if err := c.Step(start.Add(-time.Second)); err == nil {
		t.Fatal("backward time accepted")
	}
	if c.Now() != start {
		t.Fatal("clock changed")
	}
	if err := c.RunUntil(start.Add(time.Hour)); err == nil {
		t.Fatal("paused clock advanced")
	}
	c.Resume()
	if err := c.RunUntil(start.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
}
