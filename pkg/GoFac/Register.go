package GoFac

import (
	r "github.com/TaBSRest/GoFac/internal/Registration"
	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
	s "github.com/TaBSRest/GoFac/internal/Scope"
	te "github.com/TaBSRest/GoFac/internal/TaBSError"
)

func RegisterConstructor(
	containerBuilder *ContainerBuilder,
	constructor any,
	configFunctions ...func(*o.RegistrationOption) error,
) error {
	if containerBuilder.IsBuilt() {
		return te.New("Cannot register constructors after the container is built!")
	}

	registrar, err := r.NewRegistration(constructor, configFunctions...)
	if err != nil {
		return te.New("Could not register T").Join(err)
	}

	for _, key := range registrar.Options.RegistrationType {
		if _, found := containerBuilder.cache[key]; !found {
			containerBuilder.cache[key] = []*r.Registration{}
		}
		containerBuilder.cache[key] = append(containerBuilder.cache[key], registrar)
		if registrar.Options.Scope == s.PerContext {
			containerBuilder.perContextOnces = append(containerBuilder.perContextOnces, registrar)
		}
	}

	return nil
}
