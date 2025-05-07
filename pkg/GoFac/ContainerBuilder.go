package GoFac

import (
	"reflect"
	"sync"

	"github.com/TaBSRest/GoFac/internal/BuildOption"
	i "github.com/TaBSRest/GoFac/internal/Interfaces"
	r "github.com/TaBSRest/GoFac/internal/Registration"
	te "github.com/TaBSRest/GoFac/internal/TaBSError"
	ScopeOptions "github.com/TaBSRest/GoFac/pkg/GoFac/Options/Scope"
)

type ContainerBuilder struct {
	built                   bool
	cache                   map[reflect.Type][]*r.Registration
	perContextRegistrations []*r.Registration
}

func NewContainerBuilder() *ContainerBuilder {
	var perContextRegistrations []*r.Registration
	return &ContainerBuilder{
		built:                   false,
		cache:                   make(map[reflect.Type][]*r.Registration),
		perContextRegistrations: perContextRegistrations,
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

func (cb *ContainerBuilder) Build(
	uuidProvider i.UUIDProvider,
	configFunctions ...func(*BuildOption.BuildOption),
) (*Container, error) {
	if cb.built {
		return nil, te.New("This ContainerBuilder is already built!")
	}

	buildOption, err := BuildOption.New(uuidProvider)
	if err != nil {
		return nil, err
	}
	for _, optionFunction := range configFunctions {
		optionFunction(buildOption)
	}

	container := &Container{
		ContainerBuilder: cb,
		BuildOption:      buildOption,
		SingletonCache:   sync.Map{},
	}

	err = RegisterConstructor(
		cb,
		func() *Container { return container },
		ScopeOptions.AsSingleton,
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
