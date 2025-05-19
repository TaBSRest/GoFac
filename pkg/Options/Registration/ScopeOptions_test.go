package Registration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
	s "github.com/TaBSRest/GoFac/internal/Scope"
)

func TestSetToPerCall_PerformsCorrectly(t *testing.T) {
	assert := assert.New(t)

	option := o.NewRegistrationOption()

	var err error
	assert.NotPanics(
		func() { err = PerCall(option) },
		"Should not have paniced when setting the scope option!",
	)

	assert.Nil(err, "Should not have returned error!")

	assert.Equal(s.PerCall, option.Scope, "Scope must be set to PerCall!")
}

func TestSetToPerContext_PerformsCorrectly(t *testing.T) {
	assert := assert.New(t)

	option := o.NewRegistrationOption()

	var err error
	assert.NotPanics(
		func() { err = PerContext(option) },
		"Should not have paniced when setting the scope option!",
	)

	assert.Nil(err, "Should not have returned error!")

	assert.Equal(s.PerContext, option.Scope, "Scope must be set to PerContext!")
}

func TestSetToAsSingleton_PerformsCorrectly(t *testing.T) {
	assert := assert.New(t)

	option := o.NewRegistrationOption()

	var err error
	assert.NotPanics(
		func() { err = AsSingleton(option) },
		"Should not have paniced when setting the scope option!",
	)

	assert.Nil(err, "Should not have returned error!")

	assert.Equal(s.Singleton, option.Scope, "Scope must be set to Singleton!")
}
