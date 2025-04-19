package GoFac

import (
	"reflect"

	r "github.com/TaBSRest/GoFac/internal/Registration"
)

type ContainerBuilder struct {
	built bool
	cache map[reflect.Type][]*r.Registration
}

func NewContainerBuilder() *ContainerBuilder {
	return &ContainerBuilder{
		built: false,
		cache: make(map[reflect.Type][]*r.Registration),
	}
}

func GetRegistrations[T any](cb *ContainerBuilder) ([]r.Registration, bool) {
	key := reflect.TypeFor[T]()
	registrationPointer, found := cb.cache[key]
	var registrations []r.Registration = make([]r.Registration, len(registrationPointer))
	for index, ptr := range registrationPointer {
		registrations[index] = *ptr
	}
	return registrations, found
}

func GetRegistrationsFor(cb *ContainerBuilder, registrationType reflect.Type) ([]r.Registration, bool) {
	registrationPointer, found := cb.cache[registrationType]
	var registrations []r.Registration = make([]r.Registration, len(registrationPointer))
	for index, ptr := range registrationPointer {
		registrations[index] = *ptr
	}
	return registrations, found
}

func (cb *ContainerBuilder) Build() *Container {
	singletonCache := make(map[reflect.Type]*reflect.Value)
	return &Container {
		cb,
		singletonCache,
	}
}

func (cb *ContainerBuilder) IsBuilt() bool {
	return cb.built
}
