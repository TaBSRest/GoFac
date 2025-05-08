package interfaces

import (
	"reflect"

	r "github.com/TaBSRest/GoFac/internal/Registration"
)

type ContainerBuilder interface {
	GetRegistrationsFor(registrationType reflect.Type) ([]*r.Registration, bool)
	GetNamedRegistration(name string) (*r.Registration, error)
	GetGroupedRegistrations(name string) ([]*r.Registration, error)
}
