package hue

import (
	hue "github.com/collinux/gohue"
	"github.com/gissleh/lucifer"
	"strings"
)

type light struct {
	gh hue.Light
}

func (light *light) ID() string {
	return light.gh.UniqueID
}

func (light *light) Name() string {
	return light.gh.Name
}

func (light *light) SetName(name string) error {
	return light.gh.SetState(hue.LightState{
		Name: name,
	})
}

func (light *light) SetState(state lucifer.LightState) error {
	h, s, _ := state.Color.HSV()
	ghState := light.gh.State
	newState := hue.LightState{}
	changed := false

	if ghState.On != state.Power {
		changed = true
	}
	newState.On = state.Power

	brightness := uint8(state.Brightness * 254)
	if ghState.Bri != brightness {
		changed = true
	}
	newState.Bri = brightness

	if state.Color.K > 0 {
		h16 := uint16(h * (65536 / 360))
		if h16 == 0 {
			h16 = 1
		}
		if h16 != ghState.Hue {
			changed = true
		}
		newState.Hue = h16

		s8 := uint8(s * 254)
		if s8 == 0 {
			s8 = 1
		}
		if s8 != ghState.Saturation {
			changed = true
		}
		newState.Sat = s8
	} else {
		if ghState.CT != state.Color.K {
			newState.CT = uint16(state.Color.K)
			changed = true
		}
	}

	if !changed {
		return nil
	}

	if newState.On == false {
		return light.gh.SetState(hue.LightState{On: false})
	}

	return light.gh.SetState(newState)
}

func (light *light) State() (lucifer.LightState, error) {
	ghState := light.gh.State

	cm := strings.ToLower(ghState.ColorMode)
	color := lucifer.Color{}
	if cm == "ct" {
		color.SetKelvin(ghState.CT)
	} else {
		color.SetHSV(
			float64(ghState.Hue)/(65536/360),
			float64(ghState.Saturation)/254,
			float64(ghState.Bri)/254,
		)
	}

	return lucifer.LightState{
		Power:      ghState.On,
		Brightness: float64(ghState.Bri) / 254,
		Color:      color,
	}, nil
}

func (light *light) Forget() error {
	return light.gh.Delete()
}
