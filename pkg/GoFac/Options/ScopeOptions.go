package options

import (
	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
	s "github.com/TaBSRest/GoFac/internal/Scope"
)

func PerCall(option *o.RegistrationOption) error {
	option.Scope = s.PerCall
	return nil
}

func PerRequest(option *o.RegistrationOption) error {
	option.Scope = s.PerRequest
	return nil
}

func PerScope(option *o.RegistrationOption) error {
	option.Scope = s.PerScope
	return nil
}

func AsSingleton(option *o.RegistrationOption) error {
	option.Scope = s.Singleton
	return nil
}
