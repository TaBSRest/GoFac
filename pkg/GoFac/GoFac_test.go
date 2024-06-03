package gofac

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ss "github.com/pyj4104/GoFac/tests/SampleStructs"
)

func TestContainer_Constructor_InitializedProperly(t *testing.T) {
	assert := assert.New(t)

	var gofac *Container
	assert.NotPanics(
		func() {
			gofac = NewContainer()
		},
		"Should not panic when creating new container",
	)
	assert.NotNil(gofac, "Initialized container should not be nil")
	assert.NotNil(gofac.cache, "The container's cache should not be nil")
}

func TestContainer_AbleToResolveSimpleObject(t *testing.T) {
	assert := assert.New(t)

	container := NewContainer()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IIndependentStruct](container, ss.NewA)
		}, 
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	var result ss.IIndependentStruct
	resultTmp, errTmp := Resolve[ss.IIndependentStruct](container)
	result = *resultTmp
	err = errTmp
	assert.NotPanics(
		func() {
			resultTmp, errTmp := Resolve[ss.IIndependentStruct](container)
			result = *resultTmp
			err = errTmp
		}, 
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result, "Resolved object should not be nil!")
	assert.Nil(err, "Should not have any error!")
}
