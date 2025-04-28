package interfaces

import (
	"reflect"

	r "github.com/TaBSRest/GoFac/internal/Registration"
)

type ContainerBuilder interface {
	GetRegistrationsFor(registrationType reflect.Type) ([]*r.Registration, bool)
}
