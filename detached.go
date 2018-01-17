package detached

import "context"

type Detachable interface {
	Config(ctx context.Context) error
	Bootstrap(ctx context.Context) error
	Attach(ctx context.Context) error
	Status(ctx context.Context) error
	// Shutdown? Or automatic after x inactive minutes?
	// Teardown? Implode? How (and should) delete cfm/volumes/etc?
}
