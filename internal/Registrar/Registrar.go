package registrar

import (
	"fmt"
	"reflect"

	h "github.com/pyj4104/GoFac/internal/Helpers"
	o "github.com/pyj4104/GoFac/internal/RegistrationOption"
)

type Registrar struct {
	factory func(...any)(any, error)
	typeInfo reflect.Type
	options o.RegistrationOption
}

func NewRegistrar(
	typeReflection reflect.Type,
	factory func(...any)any,
	ConfigurationFunctions ...func(*o.RegistrationOption) error,
) (*Registrar, error) {
	return NewRegistrarWithError(
		typeReflection,
		wrapFactory(factory),
		ConfigurationFunctions...
	)
}

func NewRegistrarWithError(
	typeReflection reflect.Type,
	factory func(...any)(any, error),
	ConfigurationFunctions ...func(*o.RegistrationOption) error,
) (*Registrar, error) {
	if typeReflection == nil {
		return nil, h.MakeError("Registrar.NewRegistrarWithError", "Type Reflection is empty!")
	}
	if factory == nil {
		return nil, h.MakeError("Registrar.NewRegistrarWithError", "Factory is empty!")
	}

	var options *o.RegistrationOption = o.NewRegistrationOption()

	for _, config := range ConfigurationFunctions {
		if err := config(options); err != nil {
			return nil, fmt.Errorf(
				h.GetErrorMessage(
					"Registrar.NewRegistrarWithError",
					fmt.Sprintf(
						"Error registering %s",
						h.GetNameFromType(typeReflection),
					),
				),
				err,
			)
		}
	}

	return &Registrar{
		factory: factory,
		typeInfo: typeReflection,
		options: *options,
	}, nil
}

func wrapFactory(
	runOriginalFactory func(...any)any,
) func(...any)(any, error) {
	return func(args ...any)(any, error) {
		var err error = nil
		defer func() {
			if r := recover().(error); r != nil {
				err = r
			}
		}()

		val := runOriginalFactory()

		return val, err
	}
}
