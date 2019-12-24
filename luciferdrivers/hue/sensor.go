package hue

import (
	"context"
	hue "github.com/collinux/gohue"
	"github.com/gissleh/lucifer"
	"time"
)

type sensor struct {
	gh hue.Sensor

	prevButtonTime  time.Time
	prevButtonState uint16
}

func (sensor *sensor) ID() string {
	return sensor.gh.UniqueID
}

func (sensor *sensor) Name() string {
	return sensor.gh.Name
}

func (sensor *sensor) SetName(name string) error {
	return lucifer.ErrUnsupportedOperation
}

func (sensor *sensor) IsButton() bool {
	return sensor.gh.Type == "ZLLSwitch"
}

func (sensor *sensor) IsDaylight() bool {
	return sensor.gh.Type == "Daylight"
}

func (sensor *sensor) State() (lucifer.SensorState, error) {
	err := sensor.gh.Refresh()
	if err != nil {
		return lucifer.SensorState{}, err
	}

	ghState := sensor.gh.State

	var daylight *bool
	var buttonEvents []lucifer.SensorStateButtonEvent

	if ghState.ButtonEvent != 0 && time.Since(*ghState.LastUpdated.Time) < time.Second*2 {
		differentTime := !sensor.prevButtonTime.IsZero() && !ghState.LastUpdated.Time.Equal(sensor.prevButtonTime)
		if differentTime || ghState.ButtonEvent != sensor.prevButtonState {
			prevButton := sensor.prevButtonState / 1000
			prevState := sensor.prevButtonState % 1000
			currButton := ghState.ButtonEvent / 1000
			currState := ghState.ButtonEvent % 1000

			switch currState {
			case 0: // Press
				buttonEvents = append(buttonEvents, lucifer.SensorStateButtonEvent{
					Kind:   lucifer.ButtonEventPress,
					Button: int(currButton),
				})
			case 1: // Hold
				buttonEvents = append(buttonEvents, lucifer.SensorStateButtonEvent{
					Kind:   lucifer.ButtonEventHold,
					Button: int(currButton),
				})
			case 2: // Release (short)
				if prevState != 0 || prevButton != currButton {
					buttonEvents = append(buttonEvents, lucifer.SensorStateButtonEvent{
						Kind:   lucifer.ButtonEventPress,
						Button: int(currButton),
					})
				}
			case 3: // Release (long)
				buttonEvents = append(buttonEvents, lucifer.SensorStateButtonEvent{
					Kind:   lucifer.ButtonEventRelease,
					Button: int(currButton),
				})
			}
		}

		sensor.prevButtonTime = *ghState.LastUpdated.Time
		sensor.prevButtonState = ghState.ButtonEvent
	} else {
		daylightValue := sensor.gh.State.Daylight
		daylight = &daylightValue
	}

	return lucifer.SensorState{
		Time:         *sensor.gh.State.LastUpdated.Time,
		Daylight:     daylight,
		ButtonEvents: buttonEvents,
	}, nil
}

func (sensor *sensor) Forget() error {
	panic("implement me")
}

func (sensor *sensor) ButtonEvents(ctx context.Context) <-chan lucifer.SensorStateButtonEvent {
	channel := make(chan lucifer.SensorStateButtonEvent, 16)

	go func() {
		unchangedCount := 0

		defer close(channel)

		for {
			state, err := sensor.State()
			if err != nil {
				return
			}

			if len(state.ButtonEvents) > 0 {
				unchangedCount = 0
				for _, event := range state.ButtonEvents {
					select {
					case channel <- event:
					default:
					}
				}
			} else {
				unchangedCount++
			}

			var waitTime time.Duration
			if unchangedCount > 500 {
				waitTime = time.Second / 2
			} else {
				waitTime = time.Second / 50
			}

			select {
			case <-time.After(waitTime):
			case <-ctx.Done():
				return
			}
		}
	}()

	return channel
}
