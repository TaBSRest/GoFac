package ContainerBuilder

import (
	"fmt"
	"reflect"
	"sync"

	i "github.com/TaBSRest/GoFac/interfaces"
	r "github.com/TaBSRest/GoFac/internal/Registration"
	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
	te "github.com/TaBSRest/GoFac/internal/TaBSError"
	gf "github.com/TaBSRest/GoFac/pkg/Container"
	g "github.com/TaBSRest/GoFac/pkg/ContainerBuilder/Group"
	"github.com/TaBSRest/GoFac/pkg/Options"
)

type ContainerBuilder struct {
	built bool
	namedRegistrations map[string]*r.Registration
	grouped map[string]*g.Group
	cache map[reflect.Type][]*r.Registration
}

func New() *ContainerBuilder {
	return &ContainerBuilder{
		built: false,
		namedRegistrations: make(map[string]*r.Registration),
		grouped: make(map[string]*g.Group),
		cache: make(map[reflect.Type][]*r.Registration),
	}
}

func (cb *ContainerBuilder) Register(
	factory any,
	configFunctions ...func(*o.RegistrationOption) error,
) error {
	registrar, err := cb.getRegistrar(factory, configFunctions...)
	if err != nil {
		return te.New(err.GetMessage())
	}

	cb.register(registrar)

	registrationName := registrar.Options.RegistrationName
	if registrationName != "" {
		if _, found := cb.namedRegistrations[registrationName]; found {
			return te.New("The name is already taken! If the registration is for a group, please use Options.Grouped!")
		}
		cb.namedRegistrations[registrationName] = registrar
	}

	if registrar.Options.RegistrationGroup != nil {
		groupInfo := registrar.Options.RegistrationGroup
		groupName := groupInfo.Name
		if cb.grouped[groupName] == nil {
			cb.grouped[groupName] = &g.Group{
				Registrations : make([]*r.Registration, 0),
				GroupInfo: groupInfo,
			}
		} else {
			if cb.grouped[groupName].GroupType != groupInfo.GroupType {
				return te.New("The type of the group must be consistence for all group members!")
			}
		}
		cb.grouped[groupName].Registrations = append(cb.grouped[groupName].Registrations, registrar)
	}

	return nil
}

func (cb *ContainerBuilder) getRegistrar(factory any, configFunctions ...func(*o.RegistrationOption) error) (*r.Registration, i.TaBSError) {
	if cb.IsBuilt() {
		return nil, te.New("Cannot register constructors after the container is built!")
	}

	registrar, err := r.NewRegistration(factory, configFunctions...)
	if err != nil {
		return nil, te.New("Could not register T").Join(err)
	}

	return registrar, nil
}

func (cb *ContainerBuilder) register(registrar *r.Registration) {
	for _, key := range registrar.Options.RegistrationType {
		if _, found := cb.cache[key]; !found {
			cb.cache[key] = []*r.Registration{}
		}
		cb.cache[key] = append(cb.cache[key], registrar)
	}
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

func (cb *ContainerBuilder) GetNamedRegistration(name string) (*r.Registration, error) {
	registration, found := cb.namedRegistrations[name]
	if !found {
		return nil, te.New(fmt.Sprintf("Could not found the registration with the name %s", name))
	}
	return registration, nil
}

func (cb *ContainerBuilder) GetGroupedRegistrations(name string) ([]*r.Registration, error) {
	registrations, found := cb.grouped[name]
	if !found || len(registrations.Registrations) == 0 {
		return nil, te.New(fmt.Sprintf("Could not found the registration with the name %s", name))
	}
	return registrations.Registrations, nil
}
