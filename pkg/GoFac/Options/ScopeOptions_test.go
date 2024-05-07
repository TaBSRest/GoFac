package options

import (
	"testing"

	"github.com/stretchr/testify/assert"

	o "github.com/pyj4104/GoFac/internal/RegistrationOption"
	s "github.com/pyj4104/GoFac/internal/Scope"
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

func TestSetToPerRequest_PerformsCorrectly(t *testing.T) {
	assert := assert.New(t)

	option := o.NewRegistrationOption()

	var err error
	assert.NotPanics(
		func() { err = PerRequest(option) },
		"Should not have paniced when setting the scope option!",
	)

	assert.Nil(err, "Should not have returned error!")

	assert.Equal(s.PerRequest, option.Scope, "Scope must be set to PerRequest!")
}

func TestSetToPerScope_PerformsCorrectly(t *testing.T) {
	assert := assert.New(t)

	option := o.NewRegistrationOption()

	var err error
	assert.NotPanics(
		func() { err = PerScope(option) },
		"Should not have paniced when setting the scope option!",
	)

	assert.Nil(err, "Should not have returned error!")

	assert.Equal(s.PerScope, option.Scope, "Scope must be set to PerScope!")
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

