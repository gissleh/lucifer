package lucifer

type Light interface {
	// ID gets the light's ID.
	ID() string

	// Name gets the light's name.
	Name() string

	// SetName sets the light's name
	SetName(name string) error

	// State is the light's state.
	State() (LightState, error)

	// SetState syncs the state.
	SetState(state LightState) error

	// Forget forgets the light.
	Forget() error
}
