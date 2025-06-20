package GoFac

import (
	ctx "context"
	"fmt"
	"reflect"

	i "github.com/TaBSRest/GoFac/interfaces"
	c "github.com/TaBSRest/GoFac/internal/Construction"
	h "github.com/TaBSRest/GoFac/internal/Helpers"
	r "github.com/TaBSRest/GoFac/internal/Registration"
	s "github.com/TaBSRest/GoFac/internal/Scope"
	te "github.com/TaBSRest/GoFac/internal/TaBSError"
	"github.com/TaBSRest/GoFac/pkg/Container"
)

type singletonCreationResult struct {
	value *reflect.Value
	err   error
}

func Resolve[T any](container i.Container, context ctx.Context) (T, error) {
	var base T

	tInfo := reflect.TypeFor[T]()
	if h.IsArrayOrSlice(tInfo) {
		return base, te.New("Use ResolveMultiple for resolving an array or a slice")
	}

	dependency, err := resolve(container, context, tInfo)
	if err != nil {
		return base, te.New(fmt.Sprintf("Error resolving %s!", tInfo.String())).Join(err)
	}

	dependencyT, ok := dependency.Interface().(T)
	if !ok {
		return base, te.New("Could not cast to the given type! Please check the registration!")
	}
	return dependencyT, nil
}

func ResolveNamed[T any](container i.Container, context ctx.Context, name string) (T, error) {
	var base T

	tInfo := reflect.TypeFor[T]()
	if h.IsArrayOrSlice(tInfo) {
		return base, te.New("Use ResolveMultiple for resolving an array or a slice")
	}

	registration, err := container.GetNamedRegistration(name)
	if err != nil {
		return base, te.New(fmt.Sprintf("Error resolving registration with name %s", name)).Join(err)
	}

	instance, err := resolveOne(container, context, tInfo, registration)
	if err != nil {
		return base, te.New(fmt.Sprintf("Error resolving %s!", tInfo.String())).Join(err)
	}

	instanceT, ok := instance.Interface().(T)
	if !ok {
		return base, te.New("Could not cast to the given type! Please check the registration!")
	}

	return instanceT, nil
}

func ResolveMultiple[T any](container i.Container, context ctx.Context) ([]T, error) {
	var base []T

	if h.IsArrayOrSlice(reflect.TypeFor[T]()) {
		return base, te.New("Do not pass in the array of T. Just pass in T is enough")
	}

	tInfo := reflect.TypeFor[[]T]()
	resolution, err := resolve(container, context, tInfo)
	if err != nil {
		return base, te.New(fmt.Sprintf("Error resolving %s!", tInfo.String())).Join(err)
	}

	if !h.IsArrayOrSlice(resolution.Type()) {
		return base, te.New("Resulting dependency is not array, though it needed to be! Something went horribly wrong!")
	}

	var resolutions []T
	for i := range resolution.Len() {
		resolution, ok := (*resolution).Index(i).Interface().(reflect.Value).Interface().(T)
		if !ok {
			return base, te.New(fmt.Sprintf("One of the dependency could not be casted as %s", tInfo.Elem().Name()))
		}
		resolutions = append(resolutions, resolution)
	}
	return resolutions, nil
}

func ResolveGroup[T any](container i.Container, context ctx.Context, groupName string) ([]T, error) {
	var base []T

	if h.IsArrayOrSlice(reflect.TypeFor[T]()) {
		return base, te.New("Do not pass in the array of T. Just pass in T is enough")
	}

	tInfo := reflect.TypeFor[[]T]()
	registrations, found := container.GetGroupedRegistrations(groupName)
	if found != nil {
		return nil, te.New(fmt.Sprintf("%s not registered!", tInfo.Elem().String()))
	}

	resolution, err := resolveMultiple(container, context, tInfo, registrations)
	if err != nil {
		return base, te.New(fmt.Sprintf("Error resolving %s!", tInfo.String())).Join(err)
	}

	var resolutions []T
	for i := range len(resolution) {
		resolution, ok := resolution[i].Interface().(T)
		if !ok {
			return base, te.New(fmt.Sprintf("One of the dependency could not be casted as %s", tInfo.Elem().Name()))
		}
		resolutions = append(resolutions, resolution)
	}
	return resolutions, nil
}

func resolve(container i.Container, context ctx.Context, tInfo reflect.Type) (*reflect.Value, error) {
	if h.IsArrayOrSlice(tInfo) {
		registrations, found := container.GetRegistrationsFor(tInfo.Elem())
		if !found {
			return nil, te.New(fmt.Sprintf("%s is not registered!", tInfo.String()))
		}
		tmpResolution, err := resolveMultiple(container, context, tInfo.Elem(), registrations)
		if err != nil {
			return nil, err
		}
		resolution := reflect.ValueOf(h.DereferencePointedArr(tmpResolution))

		return &resolution, err
	} else {
		registrations, found := container.GetRegistrationsFor(tInfo)
		if !found {
			return nil, te.New(fmt.Sprintf("%s is not registered!", tInfo.String()))
		}
		registration := registrations[len(registrations)-1]
		return resolveOne(container, context, tInfo, registration)
	}
}

func resolveOne(container i.Container, context ctx.Context, tInfo reflect.Type, registration *r.Registration) (*reflect.Value, error) {
	constructor := registration.Construction
	dependencies, err := getDependencies(container, context, tInfo.String(), constructor)
	if err != nil {
		return nil, err
	}

	instance, err := resolveInstance(container, context, registration, constructor, tInfo.String(), dependencies)

	return instance, err
}

func resolveMultiple(container i.Container, context ctx.Context, tInfo reflect.Type, registrations []*r.Registration) ([]*reflect.Value, error) {
	var reflections []*reflect.Value
	for _, registration := range registrations {
		constructor := registration.Construction // Assuming registration.Construction is the constructor function
		dependencies, err := getDependencies(container, context, tInfo.String(), constructor)
		if err != nil {
			return nil, err
		}

		reflection, err := resolveInstance(container, context, registration, constructor, tInfo.String(), dependencies)
		if err != nil {
			return nil, te.New("Error resolving " + constructor.Info.Name()).Join(err)
		}
		reflections = append(reflections, reflection)
	}

	return reflections, nil
}

func getDependencies(
	container i.Container,
	context ctx.Context,
	originalConstructorName string,
	construction c.Construction,
) ([]*reflect.Value, error) {
	var errs []error
	dependencies := make([]*reflect.Value, construction.Info.NumIn())
	for i := range construction.Info.NumIn() {
		var dependency *reflect.Value
		dependencyInfo := construction.Info.In(i)
		if dependencyInfo == reflect.TypeFor[ctx.Context]() {
			val := reflect.ValueOf(context)
			dependency = &val
		} else {
			var err error
			dependency, err = resolve(container, context, dependencyInfo)
			if err != nil {
				errs = append(errs, err)
			}
		}

		dependencies[i] = dependency
	}
	if len(errs) != 0 {
		return nil, te.New("Could not resolve " + originalConstructorName + ":").JoinMultiple(errs)
	}
	return dependencies, nil
}

func resolveInstance(
	container i.Container, // container should be the first argument
	context ctx.Context,
	registration *r.Registration,
	ctor c.Construction,
	name string,
	dependencies []*reflect.Value,
) (*reflect.Value, error) {
	switch registration.Options.Scope {
	case s.Singleton:
		var val *reflect.Value
		var err error
		registration.SingletonOnce.Do(func() {
			val, err = runConstructor(ctor, name, dependencies)
			result := &singletonCreationResult{
				value: val,
				err:   err,
			}
			container.GetSingletonCache().Store(registration, result)
		}) // Assuming registration.SingletonOnce is a sync.Once
		return resolveSingleton(container, registration)
	case s.PerContext:
		return resolvePerContext(context, registration, ctor, name, dependencies)
	default:
		return runConstructor(ctor, name, dependencies)
	}
}

func resolveSingleton(container i.Container, registration *r.Registration) (*reflect.Value, error) {
	creationResult, found := container.GetSingletonCache().Load(registration)
	if found {
		val, ok := creationResult.(*singletonCreationResult)
		if !ok {
			return nil, te.New("What's stored in the singleton cache is not of type *singletonCreationResult")
		}
		return val.value, val.err
	}
	return nil, te.New("SingletonCache not found!")
}

func resolvePerContext(
	context ctx.Context,
	registration *r.Registration,
	ctor c.Construction,
	name string,
	dependencies []*reflect.Value,
) (*reflect.Value, error) {
	metadata, found := Container.GetMetadataFromContext(context)
	if !found {
		return nil, te.New("The context is not registered to GoFac.")
	}

	contextRegistration, found := metadata[registration]
	if !found {
		return nil, te.New("The registration is not found in the context.")
	}

	if contextRegistration.Instance != nil {
		return contextRegistration.Instance, nil
	}

	var val *reflect.Value
	var err error
	contextRegistration.Once.Do(func() {
		val, err = runConstructor(ctor, name, dependencies)
		if err == nil {
			contextRegistration.Instance = val
		}
	})

	return metadata[registration].Instance, nil
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
