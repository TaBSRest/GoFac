package GoFac

import (
	r "github.com/TaBSRest/GoFac/internal/Registration"
	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
	te "github.com/TaBSRest/GoFac/internal/TaBSError"
)

func RegisterConstructor(
	container *ContainerBuilder,
	factory any,
	configFunctions ...func(*o.RegistrationOption) error,
) error {
	if container.IsBuilt() {
		return te.New("Cannot register constructors after the container is built!")
	}

	registrar, err := r.NewRegistration(factory, configFunctions...)
	if err != nil {
		return te.New("Could not register T").Join(err)
	}

	for _, key := range registrar.Options.RegistrationType {
		if _, found := container.cache[key]; !found {
			container.cache[key] = []*r.Registration{}
		}
		container.cache[key] = append(container.cache[key], registrar)
	}

	return nil
}
