package registrationoption

import (
	"reflect"

	s "github.com/TaBSRest/GoFac/internal/Scope"
)

type RegistrationOption struct {
	Scope s.LifetimeScope
	RegistrationName string
	RegistrationType []reflect.Type
}

func NewRegistrationOption() *RegistrationOption {
	var registrationType []reflect.Type
	return &RegistrationOption{
		Scope: s.PerCall,
		RegistrationType: registrationType,
	}
}
