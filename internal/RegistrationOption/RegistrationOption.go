package RegistrationOption

import (
	"reflect"

	gi "github.com/TaBSRest/GoFac/internal/RegistrationOption/GroupInfo"
	s "github.com/TaBSRest/GoFac/internal/Scope"
)

type RegistrationOption struct {
	Scope             s.LifetimeScope
	RegistrationName  string
	RegistrationGroup *gi.GroupInfo
	RegistrationType  []reflect.Type
}

func NewRegistrationOption() *RegistrationOption {
	var registrationType []reflect.Type
	return &RegistrationOption{
		Scope:            s.PerCall,
		RegistrationType: registrationType,
	}
}
