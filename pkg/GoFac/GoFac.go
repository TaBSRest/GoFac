package gofac

import (
	"errors"
	"fmt"
	"reflect"

	h "github.com/pyj4104/GoFac/internal/Helpers"
	r "github.com/pyj4104/GoFac/internal/Registrar"
	o "github.com/pyj4104/GoFac/internal/RegistrationOption"
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

func Resolve[T interface{}](container *Container) (*T, error) {
	tInfo := reflect.TypeFor[T]()
	if isArrayOrSlice(tInfo) {
		tInfo = tInfo.Elem()
	}

	dependency, err := resolveOne(container, tInfo)
	if err != nil {
		return nil, err
	}

	dependencyT, ok := (*dependency).Interface().(T)
	if !ok {
		return nil, h.MakeError("GoFac.Resolve", "Could not cast to the given type! Please check the registration!")
	}

	return &dependencyT, nil
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

	return resolveConstructor(constructor, name, *dependencies)
}

/*
func resolveMany(container *Container, tInfos reflect.Type) ([]*reflect.Value, error) {
	tInfo := tInfos.Elem()
	if isArrayOrSlice(tInfos.Elem()) {
		// TODO: return error for now...
		return nil, h.MakeError("GoFac.resolveMany", "Arry of arry dependencies are not supported yet!")
	}
	tName := h.GetNameFromType(tInfo)
	for i, constructor := range container.cache[tName] {
		dependencies := make([]reflect.Value, constructor.Info.NumIn())
	}
	return nil, nil
}
*/

//dependenciesOfDependency := *instancesToValue(resolveds)
//dependency := registration.Constructor.Call.Call(dependenciesOfDependency)

func getDependencies(container *Container, originalConstructorName string, constructor r.Constructor) (*[]reflect.Value, error) {
	dependencies := make([]reflect.Value, constructor.Info.NumIn())
	for i := 0; i < constructor.Info.NumIn(); i++ {
		dependencyInfo := constructor.Info.In(i)
		dependency, err := resolveOne(container, dependencyInfo)
		if err != nil {
			return nil, errors.Join(
				h.MakeError(
					"GoFac.Resolve",
					"Could not resolve "+originalConstructorName+":",
				),
				err,
			)
		}
		dependencies[i] = *dependency
	}
	return &dependencies, nil
}

func resolveConstructor(constructor r.Constructor, name string, dependencies []reflect.Value) (*reflect.Value, error) {
	value := constructor.Call.Call(dependencies)
	if constructor.Info.NumOut() == 2 && !value[1].IsNil() {
		fmt.Println(value[1].Interface().(error))
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
