package GoFac

import (
	"errors"
	"fmt"
	"reflect"

	c "github.com/TaBSRest/GoFac/internal/Construction"
	h "github.com/TaBSRest/GoFac/internal/Helpers"
	r "github.com/TaBSRest/GoFac/internal/Registration"
	s "github.com/TaBSRest/GoFac/internal/Scope"
)

func Resolve[T any](container *Container) (T, error) {
	var base T

	tInfo := reflect.TypeFor[T]()
	if h.IsArrayOrSlice(tInfo) {
		return base, h.MakeError("GoFac.Resolve", "Use ResolveMultiple for resolving an array or a slice")
	}

	dependency, err := resolve(container, tInfo)
	if err != nil {
		return base, err
	}

	dependencyT, ok := dependency.Interface().(T)
	if !ok {
		return base, h.MakeError("GoFac.Resolve", "Could not cast to the given type! Please check the registration!")
	}
	return dependencyT, nil
}

func ResolveMultiple[T any](container *Container) ([]T, error) {
	var base []T

	if h.IsArrayOrSlice(reflect.TypeFor[T]()) {
		return base, h.MakeError("GoFac.Resolve", "Do not pass in the array of T. Just pass in T is enough.")
	}

	tInfo := reflect.TypeFor[[]T]()
	dependency, err := resolve(container, tInfo)
	if err != nil {
		return base, err
	}

	if !h.IsValueArrayOrSlice(*dependency) {
		return base, h.MakeError("GoFac.Resolve", "Resulting dependency is not array, though it needed to be! Something went horribly wrong!")
	}

	var resolutions []T
	for i := range (*dependency).Len() {
		resolution, ok := (*dependency).Index(i).Interface().(reflect.Value).Interface().(T)
		if !ok {
			return base, h.MakeError("GoFac.Resolve", "One of the dependency could not be casted as "+tInfo.Elem().Name())
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
	typeName := h.GetNameFor(tInfo)

	registrations, found := GetRegistrationsFor(container.ContainerBuilder, tInfo)
	if !found {
		return nil, h.MakeError("GoFac.Resolve", fmt.Sprintf("%s is not registered!", h.GetNameFor(tInfo)))
	}
	registration := registrations[len(registrations)-1]
	if item := container.resolveSingleton(registration); item != nil {
		return item, nil
	}

	constructor := registration.Construction
	dependencies, err := container.getDependencies(typeName, constructor)
	if err != nil {
		return nil, err
	}

	instance, err := resolveConstructor(constructor, typeName, dependencies)
	if err == nil && registration.Options.Scope == s.Singleton {
		container.SingletonCache[registration] = instance
	}

	return instance, err
}

func (container *Container) resolveMultiple(tInfo reflect.Type) ([]*reflect.Value, error) {
	name := h.GetNameFor(tInfo)
	registrations, found := GetRegistrationsFor(container.ContainerBuilder, tInfo)
	if !found {
		return nil, h.MakeError("GoFac.Resolve", fmt.Sprintf("%s is not registered!", name))
	}

	var reflections []*reflect.Value
	for _, registration := range registrations {
		if item := container.resolveSingleton(registration); item != nil {
			reflections = append(reflections, item)
		} else {
			constructor := registration.Construction
			dependencies, err := container.getDependencies(name, constructor)
			if err != nil {
				return nil, err
			}

			reflection, err := resolveConstructor(constructor, name, dependencies)
			if err != nil {
				return nil, errors.Join(
					h.MakeError("GoFac.Resolve", fmt.Sprintf("Error resolving %s", constructor.Info.Name())),
					err,
				)
			}
			reflections = append(reflections, reflection)
		}
	}

	return reflections, nil
}

func (container *Container) getDependencies(originalConstructorName string, construction c.Construction) ([]*reflect.Value, error) {
	dependencies := make([]*reflect.Value, construction.Info.NumIn())
	for i := range construction.Info.NumIn() {
		dependencyInfo := construction.Info.In(i)
		dependency, err := resolve(container, dependencyInfo)
		if err != nil {
			return nil, errors.Join(
				h.MakeError(
					"GoFac.Resolve",
					"Could not resolve "+originalConstructorName+":",
				),
				err,
			)
		}

		dependencies[i] = dependency
	}
	return dependencies, nil
}

func resolveConstructor(construction c.Construction, name string, dependencies []*reflect.Value) (*reflect.Value, error) {
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
		return nil, errors.Join(
			h.MakeError(
				"GoFac.Resolve",
				"Constructor of "+name+" threw an error:",
			),
			value[1].Interface().(error),
		)
	}

	return &value[0], nil
}

func (container *Container) resolveSingleton(registration *r.Registration) (*reflect.Value) {
	instance, found := container.SingletonCache[registration]
	if found {
		return instance
	}
	return nil
}
