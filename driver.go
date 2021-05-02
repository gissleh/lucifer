package lucifer

import (
	"context"
)

// A Driver is a client for interacting with a lighting system.
type Driver interface {
	SetupBridge(ctx context.Context, addr string) (Bridge, string, error)
	AddBridge(ctx context.Context, addr, key string) (Bridge, error)
	RemoveBridge(ctx context.Context, id string) error
	Bridge(id string) Bridge
	Bridges() []Bridge
}
