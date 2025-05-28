package Registration

import (
	"fmt"
	"reflect"
	"sync"

	c "github.com/TaBSRest/GoFac/internal/Construction"
	h "github.com/TaBSRest/GoFac/internal/Helpers"
	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
	te "github.com/TaBSRest/GoFac/internal/TaBSError"
)

type Registration struct {
	Construction  c.Construction
	Options       o.RegistrationOption
	SingletonOnce *sync.Once
}

func NewRegistration(
	constructor any,
	ConfigurationFunctions ...func(*o.RegistrationOption) error,
) (*Registration, error) {
	constructorTypeInfo := reflect.TypeOf(constructor)
	if err := constructorErrorChecks(constructor, constructorTypeInfo); err != nil {
		return nil, err
	}

	var options = o.NewRegistrationOption()
	var errors []error
	for _, config := range ConfigurationFunctions {
		if err := config(options); err != nil {
			errors = append(errors, err)
		}
	}

	if len(options.RegistrationType) == 0 {
		options.RegistrationType = append(options.RegistrationType, constructorTypeInfo.Out(0))
	}

	for _, tInfo := range options.RegistrationType {
		if !constructorTypeInfo.Out(0).ConvertibleTo(tInfo) {
			errors = append(
				errors,
				te.New(fmt.Sprintf("The constructor's first return value must be castible to the %s", tInfo.String())),
			)
		}
	}

	if len(errors) > 0 {
		return nil, te.New("Cannot Register").JoinMultiple(errors)
	}

	construction, err := c.NewConstruction(reflect.TypeOf(constructor), reflect.ValueOf(constructor))
	if err != nil {
		return nil, te.New("Error while registering!").Join(err)
	}

	return &Registration{
		Construction:  construction,
		Options:       *options,
		SingletonOnce: new(sync.Once),
	}, nil
}

func constructorErrorChecks(constructor any, constructorTypeInfo reflect.Type) error {
	if constructor == nil {
		return h.MakeError("Registration.NewRegistration", "Constructor is nil!")
	}
	if constructorTypeInfo.Kind() != reflect.Func {
		return h.MakeError("Registration.NewRegistration", "Constructor is not a function!")
	}
	if constructorTypeInfo.NumOut() < 1 {
		return h.MakeError("Registration.NewRegistration", "Constructor must return something!")
	}

	return nil
}
