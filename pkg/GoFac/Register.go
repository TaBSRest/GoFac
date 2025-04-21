package GoFac

import (
	"errors"

	h "github.com/TaBSRest/GoFac/internal/Helpers"
	r "github.com/TaBSRest/GoFac/internal/Registration"
	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
)

func RegisterConstructor(
	container *ContainerBuilder,
	factory any,
	configFunctions ...func(*o.RegistrationOption) error,
) error {
	if container.IsBuilt() {
		return h.MakeError("ContainerBuilder.RegisterConstructor", "Cannot register constructors after the container is built!")
	}

	registrar, err := r.NewRegistration(factory, configFunctions...)
	if err != nil {
		return errors.Join(
			h.MakeError("GoFac.RegisterConstructor", "Could not register T:"),
			err,
		)
	}

	for _, key := range registrar.Options.RegistrationType {
		if _, found := container.cache[key]; !found {
			container.cache[key] = []*r.Registration{}
		}
		container.cache[key] = append(container.cache[key], registrar)
	}

	return nil
}

