package gofac

import (
	"reflect"

	h "github.com/pyj4104/GoFac/internal/Helpers"
	r "github.com/pyj4104/GoFac/internal/Registrar"
	o "github.com/pyj4104/GoFac/internal/RegistrationOption"
)

type Container struct {
	cache map[string] *[]*r.Registrar
}

func NewContainer() *Container {
	return &Container {
		cache: make(map[string] *[]*r.Registrar),
	}
}

func RegisterConstructor[T interface{}](
	container *Container,
	factory func(...any)any,
	configFunctions ...func(*o.RegistrationOption) error,
) error {
	registrar, err := r.NewRegistrar(reflect.TypeFor[T](), factory, configFunctions...)
	if err != nil {
		return err
	}

	name := h.GetName[T]()

	if _, found := container.cache[name]; !found {
		container.cache[name] = new([]*r.Registrar)
	}

	*container.cache[name] = append(*container.cache[name], registrar)

	return nil
}
