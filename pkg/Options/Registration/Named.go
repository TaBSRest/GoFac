package Registration

import o "github.com/TaBSRest/GoFac/internal/RegistrationOption"

func Named(name string) func(option *o.RegistrationOption) error {
	return func(option *o.RegistrationOption) error {
		option.RegistrationName = name
		return nil
	}
}
