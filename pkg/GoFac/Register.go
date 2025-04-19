package GoFac

import (
	"errors"
	"reflect"

	h "github.com/TaBSRest/GoFac/internal/Helpers"
	r "github.com/TaBSRest/GoFac/internal/Registration"
	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
)

func RegisterConstructor[T any](
	container *ContainerBuilder,
	factory any,
	configFunctions ...func(*o.RegistrationOption) error,
) error {
	if container.IsBuilt() {
		return h.MakeError("ContainerBuilder.RegisterConstructor", "Cannot register constructors after the container is built!")
	}

	registrar, err := r.NewRegistration(factory, reflect.TypeFor[T](), configFunctions...)
	if err != nil {
		return errors.Join(
			h.MakeError("GoFac.RegisterConstructor", "Could not register T:"),
			err,
		)
	}

	key := reflect.TypeFor[T]()

	if _, found := container.cache[key]; !found {
		container.cache[key] = []*r.Registration{}
	}

	container.cache[key] = append(container.cache[key], registrar)

	return nil
}

