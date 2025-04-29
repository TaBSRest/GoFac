package GoFac

import (
	"reflect"
	"sync"

	r "github.com/TaBSRest/GoFac/internal/Registration"
)

type ContainerBuilder struct {
	built                   bool
	cache                   map[reflect.Type][]*r.Registration
	perContextRegistrations map[*r.Registration]struct{}
}

func NewContainerBuilder() *ContainerBuilder {
	return &ContainerBuilder{
		built:                   false,
		cache:                   make(map[reflect.Type][]*r.Registration),
		perContextRegistrations: make(map[*r.Registration]struct{}),
	}
}

func GetRegistrations[T any](cb *ContainerBuilder) ([]*r.Registration, bool) {
	key := reflect.TypeFor[T]()
	registrationPointers, found := cb.cache[key]
	registrations := make([]*r.Registration, len(registrationPointers))
	copy(registrations, registrationPointers)
	return registrations, found
}

func GetRegistrationsFor(cb *ContainerBuilder, registrationType reflect.Type) ([]*r.Registration, bool) {
	registrationPointers, found := cb.cache[registrationType]
	registrations := make([]*r.Registration, len(registrationPointers))
	copy(registrations, registrationPointers)
	return registrations, found
}

func (cb *ContainerBuilder) Build() *Container {
	return &Container{
		ContainerBuilder: cb,
		SingletonCache:   sync.Map{},
	}
}

func (cb *ContainerBuilder) IsBuilt() bool {
	return cb.built
}
