package registrar

import (
	"errors"
	"fmt"
	"reflect"

	h "github.com/pyj4104/GoFac/internal/Helpers"
	o "github.com/pyj4104/GoFac/internal/RegistrationOption"
)

type Registration struct {
	Constructor Constructor
	TypeInfo reflect.Type
	Options  o.RegistrationOption
}

type Constructor struct {
	Info reflect.Type
	Call reflect.Value
}

func NewRegistration(
	constructor interface{},
	typeInfo reflect.Type,
	ConfigurationFunctions ...func(*o.RegistrationOption) error,
) (*Registration, error) {
	if err := constructorErrorChecks(constructor, typeInfo); err != nil {
		return nil, err
	}

	var options *o.RegistrationOption = o.NewRegistrationOption()

	aggregatedErrors := []error{
		h.MakeError(
			"Registration.NewRegistration",
			fmt.Sprintf(
				"Error registering %s",
				h.GetNameFromType(typeInfo),
			),
		),
	}

	for _, config := range ConfigurationFunctions {
		if err := config(options); err != nil {
			aggregatedErrors = append(aggregatedErrors, err)
		}
	}
	if len(aggregatedErrors) > 1 {
		return nil, errors.Join(
			aggregatedErrors...,
		)
	}

	return &Registration{
		Constructor:  Constructor {
			Info: reflect.TypeOf(constructor),
			Call: reflect.ValueOf(constructor),
		},
		TypeInfo: typeInfo,
		Options:  *options,
	}, nil
}

func constructorErrorChecks(
	constructor interface{},
	typeInfo reflect.Type,
) (error) {
	constructorTypeInfo := reflect.TypeOf(constructor)
	if constructor == nil {
		return h.MakeError("Registration.NewRegistration", "Constructor is nil!")
	}
	if constructorTypeInfo.Kind() != reflect.Func {
		return h.MakeError("Registration.NewRegistration", "Constructor is not a function!")
	}
	if constructorTypeInfo.NumOut() < 1 {
		return h.MakeError("Registration.NewRegistration", "Constructor must return something!")
	}
	if typeInfo == nil {
		return h.MakeError("Registration.NewRegistration", "TypeInfo is nil!")
	}
	if typeInfo.Kind() != reflect.Interface {
		return h.MakeError("Registration.NewRegistration", "Must register interface!")
	}
	if constructorTypeInfo.Out(0) != typeInfo {
		return h.MakeError("Registration.NewRegistration", "Constructor's first return value must be of the same typeInfo!")
	}

	return nil
}
