package luciferdrivers

import (
	"github.com/gissleh/lucifer"
	"github.com/gissleh/lucifer/luciferdrivers/hue"
)

// SupportedDrivers gets a list of supported light drivers.
func SupportedDrivers() []string {
	return []string{
		"hue",
	}
}

// New creates a new driver.
func New(kind string) (lucifer.Driver, error) {
	switch kind {
	case "hue":
		return hue.New(), nil
	default:
		return nil, lucifer.ErrUnsupportedDriver
	}
}
