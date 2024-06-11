package gofac

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
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
	assert.NotPanics(
		func() {
			result, err = Resolve[ss.IIndependentStruct](container)
		}, 
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result, "Resolved object should not be nil!")
	assert.Nil(err, "Should not have any error!")
	assert.Equal("IndependentStruct", result.ReturnNameIndependentStruct(), "Functions should be able to run")
}

func TestContainer_CannotResolve_ConstructorThrowsError(t *testing.T) {
	assert := assert.New(t)

	container := NewContainer()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IIndependentStruct](container, ss.NewAReturningError)
		}, 
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	var result ss.IIndependentStruct
	assert.NotPanics(
		func() {
			result, err = Resolve[ss.IIndependentStruct](container)
		}, 
		"Should not have paniced when resolving interface!",
	)

	assert.Nil(result, "Resolved object should not be nil!")
	assert.NotNil(err, "Should not have any error!")
	assert.Equal(
		`GoFac.Resolve: Constructor of github.com/TaBSRest/GoFac/tests/SampleStructs/IIndependentStruct threw an error:
IndependentStruct: Error Forming IndependentStruct!`,
		err.Error(),
		"Error must show that constructor threw an error",
	)
}

func TestContainer_AbleToResolveInterfaceRelyingOnIndependentStruct(t *testing.T) {
	assert := assert.New(t)

	container := NewContainer()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IIndependentStruct](container, ss.NewA)
			err = RegisterConstructor[ss.IStructRelyingOnIndependentStruct](
				container,
				ss.NewStructRelyingOnIndependentStruct,
			)
		}, 
		"Should not have paniced when registering a constructor!",
	)

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	assert.Nil(err, "No Error should have happened when registering!" + errorMsg)

	var result ss.IStructRelyingOnIndependentStruct
	assert.NotPanics(
		func() {
			result, err = Resolve[ss.IStructRelyingOnIndependentStruct](container)
		}, 
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result, "Resolved object should not be nil!")
	assert.Nil(err, "Should not have any error!")
	assert.Equal("StructRelyingOnIndependentStruct", result.ReturnNameStructRelyingOnIndependentStruct(), "Functions should be able to run")
}

func TestContainer_CannotResolveInterfaceRelyingOnIndependentStruct_DependencyNotRegistered(t *testing.T) {
	assert := assert.New(t)

	container := NewContainer()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IStructRelyingOnIndependentStruct](
				container,
				ss.NewStructRelyingOnIndependentStruct,
			)
		}, 
		"Should not have paniced when registering a constructor!",
	)

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	assert.Nil(err, "No Error should have happened when registering!" + errorMsg)

	var result ss.IStructRelyingOnIndependentStruct
	assert.NotPanics(
		func() {
			result, err = Resolve[ss.IStructRelyingOnIndependentStruct](container)
		}, 
		"Should not have paniced when resolving interface!",
	)

	assert.Nil(result, "Resolved object should not be nil!")
	assert.NotNil(err, "Should not have any error!")
	assert.Equal(
		`GoFac.Resolve: Could not resolve github.com/TaBSRest/GoFac/tests/SampleStructs/IStructRelyingOnIndependentStruct:
GoFac.Resolve: github.com/TaBSRest/GoFac/tests/SampleStructs/IIndependentStruct is not registered!`,
		err.Error(),
		"Resolve must specify the cause of failure",
	)
}

func TestContainer_Resolve_ResolvesStructWithSliceInputSuccessfully(t *testing.T) {
	assert := assert.New(t)

	container := NewContainer()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IIndependentStruct](
				container,
				ss.NewA,
			)
			err = RegisterConstructor[ss.IIndependentStruct](
				container,
				ss.NewB,
			)
			err = RegisterConstructor[ss.IStructRelyingOnIndependentStruct](
				container,
				ss.NewStructRelyingOnIndependentStructs,
			)
		}, 
		"Should not have paniced when registering a constructor!",
	)

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	assert.Nil(err, "No Error should have happened when registering!" + errorMsg)

	var result ss.IStructRelyingOnIndependentStruct
	assert.NotPanics(
		func() {
			result, err = Resolve[ss.IStructRelyingOnIndependentStruct](container)
		}, 
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result, "Resolved object should not be nil!")
	assert.Nil(err, "Should not have any error!")
	assert.Equal("StructRelyingOnIndependentStructs", result.ReturnNameStructRelyingOnIndependentStruct(), "Names of the struct is different!")
	assert.Contains(result.(ss.IStructRelyingOnIndependentStructs).ReturnSubStructNames(), "IndependentStruct", "IndependentStruct should have been resolved too!")
	assert.Contains(result.(ss.IStructRelyingOnIndependentStructs).ReturnSubStructNames(), "IndependentStructB", "IndependentStructB should have been resolved too!")
}

func TestContainer_Resolve_ResolvesMultipleSuccessfully(t *testing.T) {
	assert := assert.New(t)

	container := NewContainer()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IIndependentStruct](
				container,
				ss.NewA,
			)
			err = RegisterConstructor[ss.IIndependentStruct](
				container,
				ss.NewB,
			)
		}, 
		"Should not have paniced when registering a constructor!",
	)

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	assert.Nil(err, "No Error should have happened when registering!" + errorMsg)

	var result []ss.IIndependentStruct
	assert.NotPanics(
		func() {
			result, err = ResolveMultiple[ss.IIndependentStruct](container)
		}, 
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result, "Resolved object should not be nil!")
	assert.Equal(2, len(result), "Resolved slice must have 2 items!")
	assert.Nil(err, "Should not have any error!")
	assert.Contains(result[0].ReturnNameIndependentStruct(), "IndependentStruct", "IndependentStruct should have been resolved too!")
	assert.Contains(result[1].ReturnNameIndependentStruct(), "IndependentStructB", "IndependentStructB should have been resolved too!")
}
