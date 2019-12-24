package lucifer

import "context"

type Bridge interface {
	ID() string
	Name() string
	Light(ctx context.Context, id string) (Light, error)
	Lights(ctx context.Context) ([]Light, error)
	DiscoverLights(ctx context.Context) ([]Light, error)
	Sensor(ctx context.Context, id string) (Sensor, error)
	Sensors(ctx context.Context) ([]Sensor, error)
	DiscoverSensors(ctx context.Context) ([]Sensor, error)
}
