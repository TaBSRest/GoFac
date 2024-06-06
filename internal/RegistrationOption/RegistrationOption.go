package registrationoption

import (
	s "github.com/TaBS/GoFac/internal/Scope"
)

type RegistrationOption struct {
	Scope s.LifetimeScope
}

func NewRegistrationOption() *RegistrationOption {
	return &RegistrationOption{
		Scope: s.PerCall,
	}
}
