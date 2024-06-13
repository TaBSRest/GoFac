package GoFac

import (
	"errors"
	"reflect"

	h "github.com/TaBSRest/GoFac/internal/Helpers"
	r "github.com/TaBSRest/GoFac/internal/Registration"
	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
)

func RegisterConstructor[T interface{}](
	container *ContainerBuilder,
	factory interface{},
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

	name := h.GetName[T]()

	if _, found := container.cache[name]; !found {
		container.cache[name] = []*r.Registration{}
	}

	container.cache[name] = append(container.cache[name], registrar)

	return nil
}

