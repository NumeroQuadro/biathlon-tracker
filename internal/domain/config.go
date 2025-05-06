package domain

import (
	"fmt"
	"time"
)

type Config struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

func (c *Config) GetStartTime() (time.Time, error) {
	return time.Parse("15:04:05.000", c.Start)
}

func (c *Config) GetStartDelta() (time.Time, error) {
	return time.Parse("15:04:05.000", c.StartDelta)
}

func (c *Config) Validate() error {
	if c.Laps <= 0 {
		return ErrInvalidLaps
	}
	if c.LapLen <= 0 {
		return ErrInvalidLapLen
	}
	if c.PenaltyLen <= 0 {
		return ErrInvalidPenaltyLen
	}
	if c.FiringLines <= 0 {
		return ErrInvalidFiringLines
	}

	if _, err := c.GetStartTime(); err != nil {
		return fmt.Errorf("invalid start time format: %v", err)
	}
	if _, err := c.GetStartDelta(); err != nil {
		return fmt.Errorf("invalid start delta format: %v", err)
	}

	return nil
}
