package GoFac

import (
	"reflect"

	r "github.com/TaBSRest/GoFac/internal/Registration"
)

type ContainerBuilder struct {
	built          bool
	cache          map[string][]*r.Registration
	perContextList []*r.Registration
}

func NewContainerBuilder() *ContainerBuilder {
	var perContextOnces []*r.Registration
	return &ContainerBuilder{
		built:          false,
		cache:          make(map[string][]*r.Registration),
		perContextList: perContextOnces,
	}
}

func (cb *ContainerBuilder) GetRegistrations(name string) ([]r.Registration, bool) {
	registrationPointer, found := cb.cache[name]
	var registrations []r.Registration = make([]r.Registration, len(registrationPointer))
	for index, ptr := range registrationPointer {
		registrations[index] = *ptr
	}
	return registrations, found
}

func (cb *ContainerBuilder) Build() *Container {
	singletonCache := make(map[string]*reflect.Value)
	return &Container{
		cb,
		singletonCache,
	}
}

func (cb *ContainerBuilder) IsBuilt() bool {
	return cb.built
}
