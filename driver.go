package lucifer

import (
	"context"
)

// A Driver is a client for interacting with a lighting system.
type Driver interface {
	SetupBridge(ctx context.Context, ip string) (Bridge, string, error)
	AddBridge(ctx context.Context, ip, key string) (Bridge, error)
	Bridge(id string) Bridge
	Bridges() []Bridge
}
