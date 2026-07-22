// Package simtime implements a deterministic virtual clock with no wall-clock dependency.
package simtime

import (
	"fmt"
	"time"
)

type Clock struct {
	now    time.Time
	paused bool
}

func New(at time.Time) (*Clock, error) {
	if at.IsZero() {
		return nil, fmt.Errorf("initial time is required")
	}
	return &Clock{now: at, paused: true}, nil
}
func (c *Clock) Now() time.Time { return c.now }
func (c *Clock) Paused() bool   { return c.paused }
func (c *Clock) Pause()         { c.paused = true }
func (c *Clock) Resume()        { c.paused = false }
func (c *Clock) Step(to time.Time) error {
	if to.Before(c.now) {
		return fmt.Errorf("virtual time cannot move backwards")
	}
	c.now = to
	return nil
}
func (c *Clock) RunUntil(to time.Time) error {
	if c.paused {
		return fmt.Errorf("virtual clock is paused")
	}
	return c.Step(to)
}
