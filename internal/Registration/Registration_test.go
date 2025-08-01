package Registration

import (
	"errors"
	"net/http"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
	s "github.com/TaBSRest/GoFac/internal/Scope"
	RegistrationOptions "github.com/TaBSRest/GoFac/pkg/Options/Registration"
	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
)

func TestRegistration_NewRegistration_ReturnError(t *testing.T) {
	cases := map[string]struct {
		factory        any
		typeReflection reflect.Type
		configFuncs    []func(*o.RegistrationOption) error
		msg            string
	}{
		"Constructor is nil": {
			factory:     nil,
			configFuncs: nil,
			msg:         "Registration.NewRegistration: Constructor is nil!",
		},
		"Constructor is not a function": {
			factory:     ss.IndependentStruct{},
			configFuncs: nil,
			msg:         "Registration.NewRegistration: Constructor is not a function!",
		},
		"TypeInfo and Constructor mismatch!": {
			factory: func(...any) (ss.IndependentStruct, error) { return ss.IndependentStruct{}, nil },
			configFuncs: []func(*o.RegistrationOption) error{
				func(*o.RegistrationOption) error {
					return errors.New("Error!")
				},
			},
			msg: `GoFac/internal/Registration.NewRegistration: Cannot Register
	Inner error: Error!`,
		},
		"Configuration Function Returns Error": {
			factory: func(...any) (ss.IIndependentStruct, error) { return &ss.IndependentStruct{}, nil },
			configFuncs: []func(*o.RegistrationOption) error{
				func(*o.RegistrationOption) error {
					return errors.New("Error!")
				},
			},
			msg: `GoFac/internal/Registration.NewRegistration: Cannot Register
	Inner error: Error!`,
		},
		"One of Many Configuration Functions Returns Error": {
			factory: func(...any) (ss.IIndependentStruct, error) { return &ss.IndependentStruct{}, nil },
			configFuncs: []func(*o.RegistrationOption) error{
				func(*o.RegistrationOption) error { return nil },
				func(*o.RegistrationOption) error {
					return errors.New("Error!")
				},
			},
			msg: `GoFac/internal/Registration.NewRegistration: Cannot Register
	Inner error: Error!`,
		},
		"Many Configuration Functions Returns Error": {
			factory: func(...any) (ss.IIndependentStruct, error) { return &ss.IndependentStruct{}, nil },
			configFuncs: []func(*o.RegistrationOption) error{
				func(*o.RegistrationOption) error {
					return errors.New("Error1!")
				},
				func(*o.RegistrationOption) error { return nil },
				func(*o.RegistrationOption) error {
					return errors.New("Error2!")
				},
			},
			msg: `GoFac/internal/Registration.NewRegistration: Cannot Register
	Inner error: Error1!
	Inner error: Error2!`,
		},
		"Not Castible": {
			factory: func() (ss.IndependentStruct, error) { return ss.IndependentStruct{}, nil },
			configFuncs: []func(*o.RegistrationOption) error{
				RegistrationOptions.As[http.Handler],
			},
			msg: `GoFac/internal/Registration.NewRegistration: Cannot Register
	Inner error: GoFac/internal/Registration.NewRegistration: The constructor's first return value must be castible to the http.Handler`,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			var val *Registration
			var err error
			assert.NotPanics(
				func() {
					val, err = NewRegistration(
						test.factory,
						test.configFuncs...,
					)
				},
				"Should not panic when building a registrar",
			)

			assert.Nil(val, "NewRegistration should have failed")
			assert.NotNil(err, "NewRegistration should have returned an error")
			assert.Equal(test.msg, err.Error(), "Expected error message must be updated")
		})
	}
}

func TestRegistration_NewRegistration_RegistersCorrectly(t *testing.T) {
	cases := map[string]struct {
		factory        any
		typeReflection reflect.Type
		configFuncs    []func(*o.RegistrationOption) error
		msg            string
		expectedScope  s.LifetimeScope
	}{
		"Configuration Functions are nil": {
			factory:        ss.NewA,
			typeReflection: reflect.TypeFor[ss.IIndependentStruct](),
			configFuncs:    nil,
			msg:            "",
			expectedScope:  s.PerCall,
		},
		"Configuration Functions are empty": {
			factory:        ss.NewA,
			typeReflection: reflect.TypeFor[ss.IIndependentStruct](),
			configFuncs:    *new([]func(*o.RegistrationOption) error),
			msg:            "",
			expectedScope:  s.PerCall,
		},
		"Configuration Functions set the scope to AsSingleton": {
			factory:        ss.NewA,
			typeReflection: reflect.TypeFor[ss.IIndependentStruct](),
			configFuncs: []func(*o.RegistrationOption) error{
				func(regOpt *o.RegistrationOption) error {
					regOpt.Scope = s.Singleton
					return nil
				},
			},
			msg:           "",
			expectedScope: s.Singleton,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			var val *Registration
			var err error
			assert.NotPanics(
				func() {
					val, err = NewRegistration(
						test.factory,
						test.configFuncs...,
					)
				},
				"Should not panic when building a registrar",
			)

			assert.NotNil(val, "NewRegistration should have failed!")
			// Direct func comparison is not supported in Go. However, we can still compare the names
			assert.Equal(
				runtime.FuncForPC(reflect.ValueOf(test.factory).Pointer()).Name(),
				runtime.FuncForPC(val.Construction.Value.Pointer()).Name(),
				"Constructor must be the same!",
			)
			assert.Equal(test.expectedScope, val.Options.Scope, "Scope must be the same!")
			assert.Nil(err, "NewRegistration should not have returned an error")
		})
	}
}
