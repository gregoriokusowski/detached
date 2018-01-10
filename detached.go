package detached

import "context"

type Detachable interface {
	Config(ctx context.Context) error
	Bootstrap(ctx context.Context) error
	Status(ctx context.Context) error
	Attach(ctx context.Context) error
}
