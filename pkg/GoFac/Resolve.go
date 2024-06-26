package GoFac

import (
	"errors"
	"fmt"
	"reflect"

	c "github.com/TaBSRest/GoFac/internal/Construction"
	h "github.com/TaBSRest/GoFac/internal/Helpers"
	s "github.com/TaBSRest/GoFac/internal/Scope"
)

func Resolve[T interface{}](container *Container) (T, error) {
	var base T

	tInfo := reflect.TypeFor[T]()
	if h.IsArrayOrSlice(tInfo) {
		return base, h.MakeError("GoFac.Resolve", "Use ResolveMultiple for resolving an array or a slice")
	}

	dependency, err := container.resolve(tInfo)
	if err != nil {
		return base, err
	}

	dependencyT, ok := dependency.Interface().(T)
	if !ok {
		return base, h.MakeError("GoFac.Resolve", "Could not cast to the given type! Please check the registration!")
	}
	return dependencyT, nil
}

func ResolveMultiple[T interface{}](container *Container) ([]T, error) {
	var base []T

	if h.IsArrayOrSlice(reflect.TypeFor[T]()) {
		return base, h.MakeError("GoFac.Resolve", "Do not pass in the array of T. Just pass in T is enough.")
	}

	tInfo := reflect.TypeFor[[]T]()
	dependency, err := container.resolve(tInfo)
	if err != nil {
		return base, err
	}

	if !h.IsValueArrayOrSlice(*dependency) {
		return base, h.MakeError("GoFac.Resolve", "Resulting dependency is not array, though it needed to be! Something went horribly wrong!")
	}

	var resolutions []T
	for i := 0; i < (*dependency).Len(); i++ {
		resolution, ok := (*dependency).Index(i).Interface().(reflect.Value).Interface().(T)
		if !ok {
			return base, h.MakeError("GoFac.Resolve", "One of the dependency could not be casted as "+tInfo.Elem().Name())
		}
		resolutions = append(resolutions, resolution)
	}
	return resolutions, nil
}

func (container *Container) resolve(tInfo reflect.Type) (*reflect.Value, error) {
	if h.IsArrayOrSlice(tInfo) {
		tmpResolution, err := container.resolveMultiple(tInfo.Elem())
		resolution := reflect.ValueOf(h.DereferencePointedArr(tmpResolution))
		return &resolution, err
	} else {
		return container.resolveOne(tInfo)
	}
}

func (container *Container) resolveOne(tInfo reflect.Type) (*reflect.Value, error) {
	name := h.GetNameFromType(tInfo)
	ptr, found := container.SingletonCache[name]

	if found {
		instance := ptr.Elem()
		return &instance, nil
	}

	registrations, found := (*container).GetRegistrations(name)
	if !found {
		return nil, h.MakeError("GoFac.Resolve", fmt.Sprintf("%s is not registered!", name))
	}
	registration := registrations[len(registrations)-1]

	constructor := registration.Construction
	dependencies, err := container.getDependencies(name, constructor)
	if err != nil {
		return nil, err
	}

	instance, err := resolveConstructor(constructor, name, dependencies)
	if err == nil && registration.Options.Scope == s.Singleton {
		ptr := reflect.New(instance.Type())
		ptr.Elem().Set(*instance)
		container.SingletonCache[name] = &ptr
	}

	return instance, err
}

func (container *Container) resolveMultiple(tInfo reflect.Type) ([]*reflect.Value, error) {
	name := h.GetNameFromType(tInfo)
	registrations, found := (*container).GetRegistrations(name)
	if !found {
		return nil, h.MakeError("GoFac.Resolve", fmt.Sprintf("%s is not registered!", name))
	}

	var reflections []*reflect.Value = make([]*reflect.Value, 0)
	for _, registration := range registrations {
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

	return reflections, nil
}

func (container *Container) getDependencies(originalConstructorName string, construction c.Construction) ([]*reflect.Value, error) {
	dependencies := make([]*reflect.Value, construction.Info.NumIn())
	for i := 0; i < construction.Info.NumIn(); i++ {
		dependencyInfo := construction.Info.In(i)
		dependency, err := container.resolve(dependencyInfo)
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
	for i := 0; i < construction.Info.NumIn(); i++ {
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
