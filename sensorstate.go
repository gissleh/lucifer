package lucifer

import "time"

type SensorStateButtonEventKind string

const (
	ButtonEventPress   SensorStateButtonEventKind = "Press"
	ButtonEventRelease SensorStateButtonEventKind = "Release"
	ButtonEventHold    SensorStateButtonEventKind = "Hold"
)

type SensorStateButtonEvent struct {
	Button int                        `json:"button"`
	Kind   SensorStateButtonEventKind `json:"kind"`
}

type SensorState struct {
	Time         time.Time                `json:"time"`
	Daylight     *bool                    `json:"daylight"`
	ButtonEvents []SensorStateButtonEvent `json:"buttonEvents"`
}
