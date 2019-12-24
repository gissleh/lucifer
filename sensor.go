package lucifer

import "context"

type Sensor interface {
	// ID gets the sensor's ID.
	ID() string

	// Gets whether the sensor is a button
	IsButton() bool

	// Gets whether the sensor is a daylight sensor
	IsDaylight() bool

	// Name gets the sensor's name.
	Name() string

	// SetName sets the sensor's name
	SetName(name string) error

	// State is the sensor's state.
	State() (SensorState, error)

	// SubscribeButtonEvents subscribes to button events.
	ButtonEvents(ctx context.Context) <-chan SensorStateButtonEvent

	// Forget forgets the sensor.
	Forget() error
}
