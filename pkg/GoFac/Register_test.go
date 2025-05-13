package GoFac

import (
	"testing"

	"github.com/stretchr/testify/assert"

	AsOptions "github.com/TaBSRest/GoFac/pkg/GoFac/Options/As"
	ScopeOptions "github.com/TaBSRest/GoFac/pkg/GoFac/Options/Scope"
	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(containerBuilder, ss.NewA, AsOptions.As[ss.IIndependentStruct])
		},
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")
}

func TestRegister_AddsRegistrationToPerContextRegistrations_WhenScopeIsPerContext(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()

	err := RegisterConstructor(
		containerBuilder,
		ss.NewA,
		AsOptions.As[ss.IIndependentStruct],
		ScopeOptions.PerContext,
	)
	assert.Nil(err, "No Error should have happened when registering with PerContext scope")

	assert.Len(containerBuilder.perContextRegistrations, 1,
		"One registration should be in perContextRegistrations",
	)
}
