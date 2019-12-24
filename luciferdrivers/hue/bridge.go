package hue

import (
	"context"
	"errors"
	"fmt"
	hue "github.com/collinux/gohue"
	"github.com/gissleh/lucifer"
)

type bridge struct {
	gh *hue.Bridge
}

func (bridge *bridge) ID() string {
	return bridge.gh.Info.Device.SerialNumber
}

func (bridge *bridge) Name() string {
	return bridge.gh.Info.Device.FriendlyName
}

func (bridge *bridge) Light(ctx context.Context, id string) (lucifer.Light, error) {
	lights, err := bridge.Lights(ctx)
	if err != nil {
		return nil, err
	}

	for _, light := range lights {
		if light.ID() == id {
			return light, nil
		}
	}

	return nil, errors.New("light not found")
}

func (bridge *bridge) Lights(ctx context.Context) ([]lucifer.Light, error) {
	ghLights, err := bridge.gh.GetAllLights()
	if err != nil {
		return nil, err
	}

	lights := make([]lucifer.Light, len(ghLights))
	for i, ghLight := range ghLights {
		lights[i] = &light{gh: ghLight}
	}

	return lights, nil
}

func (bridge *bridge) DiscoverLights(ctx context.Context) ([]lucifer.Light, error) {
	before, err := bridge.Lights(ctx)
	if err != nil {
		return nil, err
	}

	err = bridge.gh.FindNewLights()
	if err != nil {
		return nil, err
	}

	after, err := bridge.Lights(ctx)
	if err != nil {
		return nil, err
	}

	newLights := make([]lucifer.Light, 0, 8)
	for _, alight := range after {
		found := false
		for _, blight := range before {
			if alight.ID() == blight.ID() {
				found = true
			}
		}

		if !found {
			newLights = append(newLights, alight)
		}
	}

	return newLights, nil
}

func (bridge *bridge) Sensor(ctx context.Context, id string) (lucifer.Sensor, error) {
	sensors, err := bridge.Sensors(ctx)
	if err != nil {
		return nil, err
	}

	for _, sensor := range sensors {
		if sensor.ID() == id {
			return sensor, nil
		}
	}

	return nil, errors.New("sensor not found")
}

func (bridge *bridge) Sensors(ctx context.Context) ([]lucifer.Sensor, error) {
	ghSensors, err := bridge.gh.GetAllSensors()
	if err != nil {
		return nil, err
	}

	sensors := make([]lucifer.Sensor, 0, len(ghSensors))
	for _, ghSensor := range ghSensors {
		if ghSensor.UniqueID == "" {
			continue
		}

		sensors = append(sensors, &sensor{gh: ghSensor})
	}

	return sensors, nil
}

func (bridge *bridge) DiscoverSensors(ctx context.Context) ([]lucifer.Sensor, error) {
	uri := fmt.Sprintf("/api/%s/sensors", bridge.gh.Username)
	_, _, err := bridge.gh.Post(uri, nil)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
