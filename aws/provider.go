package aws

import "github.com/gregoriokusowski/detached"

func New() detached.Detachable {
	return &Aws{}
}

type Aws struct {
}
