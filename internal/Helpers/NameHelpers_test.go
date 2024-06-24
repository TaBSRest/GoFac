package Helpers

import (
	"reflect"
	"testing"

	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
	"github.com/stretchr/testify/assert"
)

func Test_GetName_CreatesNameWithPackageName(t *testing.T) {
	assert := assert.New(t)

	var interfaceName string
	assert.NotPanics(
		func() { interfaceName = GetName[ss.IIndependentStruct]() },
		"Should not throw when getting the interface name",
	)
	assert.NotEmpty(
		interfaceName,
		"Interface name should not be empty",
	)
	assert.Equal(
		"github.com/TaBSRest/GoFac/tests/SampleStructs/IIndependentStruct",
		interfaceName,
		"Interface name must match",
	)
}

func Test_GetNameFromType_CreatesNameWithPackageName(t *testing.T) {
	assert := assert.New(t)

	var interfaceName string
	assert.NotPanics(
		func() { interfaceName = GetNameFromType(reflect.TypeFor[ss.IIndependentStruct]()) },
		"Should not throw when getting the interface name",
	)
	assert.NotEmpty(
		interfaceName,
		"Interface name should not be empty",
	)
	assert.Equal(
		"github.com/TaBSRest/GoFac/tests/SampleStructs/IIndependentStruct",
		interfaceName,
		"Interface name must match",
	)
}
