package detached

type Detachable interface {
	Bootstrap() error
	Status() error
	Attach() error
}
