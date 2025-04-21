package GoFac

import (
	"fmt"
	"reflect"

	c "github.com/TaBSRest/GoFac/internal/Construction"
	h "github.com/TaBSRest/GoFac/internal/Helpers"
	r "github.com/TaBSRest/GoFac/internal/Registration"
	s "github.com/TaBSRest/GoFac/internal/Scope"
	te "github.com/TaBSRest/GoFac/internal/TaBSError"
)

type singletonCreationResult struct {
	value *reflect.Value
	err error
}

func Resolve[T any](container *Container) (T, error) {
	var base T

	tInfo := reflect.TypeFor[T]()
	if h.IsArrayOrSlice(tInfo) {
		return base, te.New("Use ResolveMultiple for resolving an array or a slice")
	}

	dependency, err := resolve(container, tInfo)
	if err != nil {
		return base, te.New(fmt.Sprintf("Error resolving %s!", tInfo.String())).Join(err)
	}

	dependencyT, ok := dependency.Interface().(T)
	if !ok {
		return base, te.New("Could not cast to the given type! Please check the registration!")
	}
	return dependencyT, nil
}

func ResolveMultiple[T any](container *Container) ([]T, error) {
	var base []T

	if h.IsArrayOrSlice(reflect.TypeFor[T]()) {
		return base, te.New("Do not pass in the array of T. Just pass in T is enough")
	}

	tInfo := reflect.TypeFor[[]T]()
	dependency, err := resolve(container, tInfo)
	if err != nil {
		return base, te.New(fmt.Sprintf("Error resolving %s!", tInfo.String())).Join(err)
	}

	if !h.IsValueArrayOrSlice(*dependency) {
		return base, te.New("Resulting dependency is not array, though it needed to be! Something went horribly wrong!")
	}

	var resolutions []T
	for i := range (*dependency).Len() {
		resolution, ok := (*dependency).Index(i).Interface().(reflect.Value).Interface().(T)
		if !ok {
			return base, te.New(fmt.Sprintf("One of the dependency could not be casted as %s", tInfo.Elem().Name()))
		}
		resolutions = append(resolutions, resolution)
	}
	return resolutions, nil
}

func resolve(container *Container, tInfo reflect.Type) (*reflect.Value, error) {
	if h.IsArrayOrSlice(tInfo) {
		tmpResolution, err := container.resolveMultiple(tInfo.Elem())
		resolution := reflect.ValueOf(h.DereferencePointedArr(tmpResolution))
		return &resolution, err
	} else {
		return resolveOne(container, tInfo)
	}
}

func resolveOne(container *Container, tInfo reflect.Type) (*reflect.Value, error) {
	registrations, found := GetRegistrationsFor(container.ContainerBuilder, tInfo)
	if !found {
		return nil, te.New(fmt.Sprintf("%s is not registered!", tInfo.String()))
	}
	registration := registrations[len(registrations)-1]

	constructor := registration.Construction
	dependencies, err := container.getDependencies(tInfo.String(), constructor)
	if err != nil {
		return nil, err
	}

	instance, err := resolveInstance(container, registration, constructor, tInfo.String(), dependencies)

	return instance, err
}

func (container *Container) resolveMultiple(tInfo reflect.Type) ([]*reflect.Value, error) {
	registrations, found := GetRegistrationsFor(container.ContainerBuilder, tInfo)
	if !found {
		return nil, te.New(fmt.Sprintf("%s is not registered!", tInfo.Name()))
	}

	var reflections []*reflect.Value
	for _, registration := range registrations {
		constructor := registration.Construction
		dependencies, err := container.getDependencies(tInfo.String(), constructor)
		if err != nil {
			return nil, err
		}

		reflection, err := resolveInstance(container, registration, constructor, tInfo.String(), dependencies)
		if err != nil {
			return nil, te.New("Error resolving "+constructor.Info.Name()).Join(err)
		}
		reflections = append(reflections, reflection)
	}

	return reflections, nil
}

func (container *Container) getDependencies(
	originalConstructorName string,
	construction c.Construction,
) ([]*reflect.Value, error) {
	dependencies := make([]*reflect.Value, construction.Info.NumIn())
	for i := range construction.Info.NumIn() {
		dependencyInfo := construction.Info.In(i)
		dependency, err := resolve(container, dependencyInfo)
		if err != nil {
			return nil, te.New("Could not resolve " + originalConstructorName + ":").Join(err)
		}

		dependencies[i] = dependency
	}
	return dependencies, nil
}

func resolveInstance(
	container *Container,
	registration *r.Registration,
	ctor c.Construction,
	name string,
	dependencies []*reflect.Value,
) (*reflect.Value, error) {
	if registration.Options.Scope == s.Singleton {
		var val *reflect.Value
		var err error
		registration.SingletonOnce.Do(func() {
			val, err = runConstructor(ctor, name, dependencies)
			result := &singletonCreationResult{
				value: val,
				err: err,
			}
			container.SingletonCache.Store(registration, result)
		})
		return container.resolveSingleton(registration)
	}

	return runConstructor(ctor, name, dependencies)
}

func runConstructor(construction c.Construction, name string, dependencies []*reflect.Value) (*reflect.Value, error) {
	types := make([]reflect.Type, construction.Info.NumIn())
	for i := range construction.Info.NumIn() {
		types[i] = construction.Info.In(i)
	}
	castedDependencies, err := h.CastInput(h.DereferencePointedArr(dependencies), types)
	if err != nil {
		return nil, err
	}

	value := construction.Value.Call(castedDependencies)
	if construction.Info.NumOut() == 2 && !value[1].IsNil() {
		return nil, te.New(fmt.Sprintf("Constructor of %s threw an error", name)).Join(value[1].Interface().(error))
	}

	return &value[0], nil
}

func (container *Container) resolveSingleton(registration *r.Registration) (*reflect.Value, error) {
	creationResult, found := container.SingletonCache.Load(registration)
	if found {
		val, ok := creationResult.(*singletonCreationResult)
		if !ok {
			return nil, te.New("What's stored in the singleton cache is not of type *singletonCreationResult")
		}
		return val.value, val.err
	}
	return nil, nil
}
