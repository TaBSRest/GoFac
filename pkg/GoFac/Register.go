package GoFac

import (
	"errors"
	"reflect"

	h "github.com/TaBSRest/GoFac/internal/Helpers"
	r "github.com/TaBSRest/GoFac/internal/Registration"
	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
	s "github.com/TaBSRest/GoFac/internal/Scope"
)

func RegisterConstructor[T interface{}](
	containerBuilder *ContainerBuilder,
	factory interface{},
	configFunctions ...func(*o.RegistrationOption) error,
) error {
	if containerBuilder.IsBuilt() {
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

	if _, found := containerBuilder.cache[name]; !found {
		containerBuilder.cache[name] = []*r.Registration{}
	}

	containerBuilder.cache[name] = append(containerBuilder.cache[name], registrar)

	if registrar.Options.Scope == s.PerContext {
		containerBuilder.perContextList = append(containerBuilder.perContextList, registrar)
	}

	return nil
}
