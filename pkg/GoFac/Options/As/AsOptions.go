package Options

import (
	r "reflect"

	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
)

func As[T any](option *o.RegistrationOption) error {
	tType := r.TypeFor[T]()
	option.RegistrationType = append(option.RegistrationType, tType)
	return nil
}
