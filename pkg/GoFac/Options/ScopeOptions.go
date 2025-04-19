package Options

import (
	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
	s "github.com/TaBSRest/GoFac/internal/Scope"
)

func PerCall(option *o.RegistrationOption) error {
	option.Scope = s.PerCall
	return nil
}

func PerContext(option *o.RegistrationOption) error {
	option.Scope = s.PerContext
	return nil
}

func AsSingleton(option *o.RegistrationOption) error {
	option.Scope = s.Singleton
	return nil
}
