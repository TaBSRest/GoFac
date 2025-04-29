package RegistrationOption

import (
	"reflect"
	s "github.com/TaBSRest/GoFac/internal/Scope"
	gi "github.com/TaBSRest/GoFac/internal/RegistrationOption/GroupInfo"
)

type RegistrationOption struct {
	Scope s.LifetimeScope
	RegistrationName string
	RegistrationGroup *gi.GroupInfo
	RegistrationType []reflect.Type
}

func NewRegistrationOption() *RegistrationOption {
	var registrationType []reflect.Type
	return &RegistrationOption{
		Scope: s.PerCall,
		RegistrationType: registrationType,
	}
}
