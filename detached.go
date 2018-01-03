package detached

import "context"

type Detachable interface {
	Bootstrap(ctx context.Context) error
	Status(ctx context.Context) error
	Attach(ctx context.Context) error
}
