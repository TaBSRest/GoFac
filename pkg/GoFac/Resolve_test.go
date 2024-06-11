package GoFac

import (
	"testing"

	"github.com/stretchr/testify/assert"

	o "github.com/TaBSRest/GoFac/pkg/GoFac/Options"
	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
)


func TestContainer_AbleToResolveSimpleObject(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IIndependentStruct](containerBuilder, ss.NewA)
		}, 
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	container := containerBuilder.Build()

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

func TestResolve_ResolvesTwoDifferentInstances_InstancesAreNotRegisteredAsSingleton(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IIndependentStruct](containerBuilder, ss.NewS)
		}, 
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	container := containerBuilder.Build()

	var result1 ss.IIndependentStruct
	assert.NotPanics(
		func() {
			result1, err = Resolve[ss.IIndependentStruct](container)
		}, 
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result1, "First object should not be nil!")
	assert.Nil(err, "First resolved object should not have any error!")
	assert.Equal("SingletonStruct", result1.ReturnNameIndependentStruct(), "Functions should be able to run")

	var result2 ss.IIndependentStruct
	assert.NotPanics(
		func() {
			result2, err = Resolve[ss.IIndependentStruct](container)
		}, 
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result2, "Second object should not be nil!")
	assert.Nil(err, "Second resolved object should not have any error!")
	assert.Equal("SingletonStruct", result2.ReturnNameIndependentStruct(), "Functions should be able to run")
}

func TestResolve_ResolvesOneInstance_ObjectRegisteredAsSingleton(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IIndependentStruct](containerBuilder, ss.NewS, o.AsSingleton)
		}, 
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	container := containerBuilder.Build()

	var result1 ss.IIndependentStruct
	assert.NotPanics(
		func() {
			result1, err = Resolve[ss.IIndependentStruct](container)
		}, 
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result1, "First object should not be nil!")
	assert.Nil(err, "First resolved object should not have any error!")
	assert.Equal("SingletonStruct", result1.ReturnNameIndependentStruct(), "Functions should be able to run")

	assert.NotNil(result1, "First object should not be nil!")
	assert.Nil(err, "First resolved object should not have any error!")
	assert.Equal("Already Ran!", result1.ReturnNameIndependentStruct(), "Functions should be able to run")

	var result2 ss.IIndependentStruct
	assert.NotPanics(
		func() {
			result2, err = Resolve[ss.IIndependentStruct](container)
		}, 
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result2, "Second object should not be nil!")
	assert.Nil(err, "Second resolved object should not have any error!")
	assert.Equal("Already Ran!", result2.ReturnNameIndependentStruct(), "Functions should be able to run")

	assert.Same(result1, result2, "They must be the same!")
}

func TestContainer_CannotResolve_ConstructorThrowsError(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IIndependentStruct](containerBuilder, ss.NewAReturningError)
		}, 
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	container := containerBuilder.Build()

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

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IIndependentStruct](containerBuilder, ss.NewA)
			err = RegisterConstructor[ss.IStructRelyingOnIndependentStruct](
				containerBuilder,
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

	container := containerBuilder.Build()

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

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IStructRelyingOnIndependentStruct](
				containerBuilder,
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

	container := containerBuilder.Build()

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

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IIndependentStruct](
				containerBuilder,
				ss.NewA,
			)
			err = RegisterConstructor[ss.IIndependentStruct](
				containerBuilder,
				ss.NewB,
			)
			err = RegisterConstructor[ss.IStructRelyingOnIndependentStruct](
				containerBuilder,
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

	container := containerBuilder.Build()

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

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor[ss.IIndependentStruct](
				containerBuilder,
				ss.NewA,
			)
			err = RegisterConstructor[ss.IIndependentStruct](
				containerBuilder,
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

	container := containerBuilder.Build()

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
