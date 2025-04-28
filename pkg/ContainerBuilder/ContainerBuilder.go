package ContainerBuilder

import (
	"reflect"
	"sync"

	i "github.com/TaBSRest/GoFac/interfaces"
	gf "github.com/TaBSRest/GoFac/pkg/Container"
	r "github.com/TaBSRest/GoFac/internal/Registration"
	te "github.com/TaBSRest/GoFac/internal/TaBSError"
	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
	"github.com/TaBSRest/GoFac/pkg/Options"
)

type ContainerBuilder struct {
	built bool
	cache map[reflect.Type][]*r.Registration
}

func New() *ContainerBuilder {
	return &ContainerBuilder{
		built: false,
		cache: make(map[reflect.Type][]*r.Registration),
	}
}

func (cb *ContainerBuilder) Register(
	factory any,
	configFunctions ...func(*o.RegistrationOption) error,
) error {
	if cb.IsBuilt() {
		return te.New("Cannot register constructors after the container is built!")
	}

	registrar, err := r.NewRegistration(factory, configFunctions...)
	if err != nil {
		return te.New("Could not register T").Join(err)
	}

	for _, key := range registrar.Options.RegistrationType {
		if _, found := cb.cache[key]; !found {
			cb.cache[key] = []*r.Registration{}
		}
		cb.cache[key] = append(cb.cache[key], registrar)
	}

	return nil
}

func GetRegistrations[T any](cb *ContainerBuilder) ([]*r.Registration, bool) {
	key := reflect.TypeFor[T]()
	registrations, found := cb.cache[key]
	return registrations, found
}

func (cb *ContainerBuilder) GetRegistrationsFor(registrationType reflect.Type) ([]*r.Registration, bool) {
	registrations, found := cb.cache[registrationType]
	return registrations, found
}

func (cb *ContainerBuilder) Build() (*gf.Container, error) {
	if cb.built {
		return nil, te.New("This ContainerBuilder is already built!")
	}

	container := &gf.Container{
		ContainerBuilder: cb,
		SingletonCache:   new(sync.Map),
	}

	err := cb.Register(
		func() i.Container { return container },
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
