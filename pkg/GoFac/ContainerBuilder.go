package GoFac

import (
	"reflect"
	"sync"

	r "github.com/TaBSRest/GoFac/internal/Registration"
	te "github.com/TaBSRest/GoFac/internal/TaBSError"
	"github.com/TaBSRest/GoFac/pkg/GoFac/Options"
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
	registrations, found := cb.cache[key]
	return registrations, found
}

func GetRegistrationsFor(cb *ContainerBuilder, registrationType reflect.Type) ([]*r.Registration, bool) {
	registrations, found := cb.cache[registrationType]
	return registrations, found
}

func (cb *ContainerBuilder) Build() (*Container, error) {
	if cb.built {
		return nil, te.New("This ContainerBuilder is already built!")
	}

	container := &Container{
		ContainerBuilder: cb,
		SingletonCache:   sync.Map{},
	}

	err := RegisterConstructor(
		cb,
		func() *Container { return container },
		Options.AsSingleton,
	)
	if err != nil {
		return nil, err
	}

	cb.built = true

	return container, nil
}

func (cb *ContainerBuilder) IsBuilt() bool {
	return cb.built
}
