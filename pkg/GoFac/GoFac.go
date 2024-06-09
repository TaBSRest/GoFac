package gofac

import (
	"errors"
	"fmt"
	"reflect"

	h "github.com/TaBS/GoFac/internal/Helpers"
	r "github.com/TaBS/GoFac/internal/Registrar"
	o "github.com/TaBS/GoFac/internal/RegistrationOption"
)

type Container struct {
	cache map[string][]*r.Registration
}

func NewContainer() *Container {
	return &Container{
		cache: make(map[string][]*r.Registration),
	}
}

func RegisterConstructor[T interface{}](
	container *Container,
	factory interface{},
	configFunctions ...func(*o.RegistrationOption) error,
) error {
	registrar, err := r.NewRegistration(factory, reflect.TypeFor[T](), configFunctions...)
	if err != nil {
		return errors.Join(
			h.MakeError("GoFac.RegisterConstructor", "Could not register T:"),
			err,
		)
	}

	name := h.GetName[T]()

	if _, found := container.cache[name]; !found {
		container.cache[name] = []*r.Registration{}
	}

	container.cache[name] = append(container.cache[name], registrar)

	return nil
}

func Resolve[T interface{}](container *Container) (T, error) {
	var base T

	tInfo := reflect.TypeFor[T]()
	if isArrayOrSlice(tInfo) {
		return base, h.MakeError("GoFac.Resolve", "Use ResolveMultiple for resolving an array or a slice")
	}

	dependency, err := resolve(container, tInfo)
	if err != nil {
		return base, err
	}

	dependencyT, ok := (*dependency).Interface().(T)
	if !ok {
		return base, h.MakeError("GoFac.Resolve", "Could not cast to the given type! Please check the registration!")
	}
	return dependencyT, nil
}

func ResolveMultiple[T interface{}](container *Container) ([]T, error) {
	var base []T

	if isArrayOrSlice(reflect.TypeFor[T]()) {
		return base, h.MakeError("GoFac.Resolve", "Do not pass in the array of T. Just pass in T is enough.")
	}

	tInfo := reflect.TypeFor[[]T]()
	dependency, err := resolve(container, tInfo)
	if err != nil {
		return base, err
	}

	if !isValueArrayOrSlice(*dependency) {
		return base, h.MakeError("GoFac.Resolve", "Resulting dependency is not array, though it needed to be! Something went horribly wrong!")
	}

	var resolutions []T
	for i := 0; i < (*dependency).Len(); i++ {
		resolution, ok := (*dependency).Index(i).Interface().(reflect.Value).Interface().(T)
		if !ok {
			return base, h.MakeError("GoFac.Resolve", "One of the dependency could not be casted as " + tInfo.Elem().Name())
		}
		resolutions = append(resolutions, resolution)
	}
	return resolutions, nil
}

func resolve(container *Container, tInfo reflect.Type) (*reflect.Value, error) {
	if isArrayOrSlice(tInfo) {
		tmpResolution, err := resolveMultiple(container, tInfo.Elem())
		resolution := reflect.ValueOf(dereferencePointedArr(tmpResolution))
		return &resolution, err
	} else {
		return resolveOne(container, tInfo)
	}
}

func resolveOne(container *Container, tInfo reflect.Type) (*reflect.Value, error) {
	name := h.GetNameFromType(tInfo)
	registrations, found := container.cache[name]
	if !found {
		return nil, h.MakeError("GoFac.Resolve", fmt.Sprintf("%s is not registered!", name))
	}
	registration := registrations[len(registrations)-1]

	constructor := registration.Constructor
	dependencies, err := getDependencies(container, name, constructor)
	if err != nil {
		return nil, err
	}

	return resolveConstructor(constructor, name, dependencies)
}

func resolveMultiple(container *Container, tInfo reflect.Type) ([]*reflect.Value, error) {
	name := h.GetNameFromType(tInfo)
	registrations, found := container.cache[name]
	if !found {
		return nil, h.MakeError("GoFac.Resolve", fmt.Sprintf("%s is not registered!", name))
	}

	var reflections []*reflect.Value = make([]*reflect.Value, 0)
	for _, registration := range registrations {
		constructor := registration.Constructor
		dependencies, err := getDependencies(container, name, constructor)
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

func getDependencies(container *Container, originalConstructorName string, constructor r.Constructor) ([]*reflect.Value, error) {
	dependencies := make([]*reflect.Value, constructor.Info.NumIn())
	for i := 0; i < constructor.Info.NumIn(); i++ {
		dependencyInfo := constructor.Info.In(i)
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

func resolveConstructor(constructor r.Constructor, name string, dependencies []*reflect.Value) (*reflect.Value, error) {
	types := make([]reflect.Type, constructor.Info.NumIn())
	for i := 0; i < constructor.Info.NumIn(); i++ {
		types[i] = constructor.Info.In(i)
	}
	castedDependencies, err := castInput(dereferencePointedArr(dependencies), types)
	if err != nil {
		return nil, err
	}

	value := constructor.Call.Call(castedDependencies)
	if constructor.Info.NumOut() == 2 && !value[1].IsNil() {
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

func isArrayOrSlice(info reflect.Type) bool {
	return info.Kind() == reflect.Slice || info.Kind() == reflect.Array
}

func isValueArrayOrSlice(value reflect.Value) bool {
	return value.Kind() == reflect.Slice || value.Kind() == reflect.Array
}

func dereferencePointedArr(pointedArr []*reflect.Value) ([]reflect.Value) {
	arr := make([]reflect.Value, len(pointedArr))
	for i, val := range pointedArr {
		arr[i] = *val
	}
	return arr
}

func castInput(uncastedInputs []reflect.Value, target []reflect.Type) ([]reflect.Value, error) {
	if len(uncastedInputs) != len(target) {
		return nil, h.MakeError("GoFac.Resolve", "Something went horribly wrong! The number of inputs to a constructor does not match with the dependencies retrieved!")
	}

	var castedInputs []reflect.Value = make([]reflect.Value, len(uncastedInputs))
	for i, uncastedInput := range uncastedInputs {
		if isValueArrayOrSlice(uncastedInput) {
			elementaryType := target[i].Elem()
			castedInput := reflect.MakeSlice(reflect.SliceOf(elementaryType), 0, 10)
			for _, input := range uncastedInput.Interface().([]reflect.Value) {
				if !input.CanConvert(elementaryType) {
					return nil, h.MakeError("GoFac.Resolve", "Cannot convert " + input.Type().Name() + " to " + elementaryType.Name())
				}
				castedInput = reflect.Append(castedInput, input.Convert(elementaryType))
				castedInputs[i] = castedInput
			}
		} else {
			if !uncastedInput.CanConvert(target[i]) {
				return nil, h.MakeError("GoFac.Resolve", "Cannot convert " + uncastedInput.Type().Name() + " to " + target[i].Name())
			}
			castedInputs[i] = uncastedInput.Convert(target[i])
		}
	}

	return castedInputs, nil
}
