package GoFac

import (
	ctx "context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	o "github.com/TaBSRest/GoFac/pkg/GoFac/Options"
	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
)

func TestResolve_AbleToResolveSimpleObject(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(containerBuilder, ss.NewA, o.As[ss.IIndependentStruct])
		},
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	container, err := containerBuilder.Build()
	assert.Nil(err)

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

func TestResolve_AbleToResolveSelf(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(containerBuilder, ss.NewIndependentStruct)
		},
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	container, err := containerBuilder.Build()
	assert.Nil(err)

	var result *ss.IndependentStruct
	assert.NotPanics(
		func() {
			result, err = Resolve[*ss.IndependentStruct](container)
		},
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result, "Resolved object should not be nil!")
	assert.Nil(err, "Should not have any error!")
	assert.Equal("IndependentStruct", result.ReturnNameIndependentStruct(), "Functions should be able to run")
}

func TestResolve_AbleToResolveUnderMultipleInterfaces(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(
				containerBuilder,
				ss.NewIndependentStruct,
				o.As[ss.IIndependentStruct],
				o.As[ss.IIndependentStruct2],
			)
		},
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	container, err := containerBuilder.Build()
	assert.Nil(err)

	var result1 ss.IIndependentStruct
	assert.NotPanics(
		func() {
			result1, err = Resolve[ss.IIndependentStruct](container)
		},
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result1, "Resolved object should not be nil!")
	assert.Nil(err, "Should not have any error!")
	assert.Equal("IndependentStruct", result1.ReturnNameIndependentStruct(), "Functions should be able to run")

	var result2 ss.IIndependentStruct2
	assert.NotPanics(
		func() {
			result2, err = Resolve[ss.IIndependentStruct2](container)
		},
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result2, "Resolved object should not be nil!")
	assert.Nil(err, "Should not have any error!")
	assert.Equal("IndependentStruct", result2.ReturnNameIndependentStruct(), "Functions should be able to run")
}

func TestResolve_ResolvesTwoDifferentInstances_InstancesAreNotRegisteredAsSingleton(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(containerBuilder, ss.NewS, o.As[ss.IIndependentStruct])
		},
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	container, err := containerBuilder.Build()
	assert.Nil(err)

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
			err = RegisterConstructor(containerBuilder, ss.NewS, o.AsSingleton, o.As[ss.IIndependentStruct])
		},
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	container, err := containerBuilder.Build()
	assert.Nil(err)

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

func TestResolve_ResolvesOneInstance_ObjectRegisteredAsSingletonUnderDifferentType(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(
				containerBuilder,
				ss.NewS,
				o.AsSingleton,
				o.As[ss.IIndependentStruct],
				o.As[ss.IIndependentStruct2],
			)
		},
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	container, err := containerBuilder.Build()
	assert.Nil(err)

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

	var result2 ss.IIndependentStruct2
	assert.NotPanics(
		func() {
			result2, err = Resolve[ss.IIndependentStruct2](container)
		},
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result2, "Second object should not be nil!")
	assert.Nil(err, "Second resolved object should not have any error!")
	assert.Equal("Already Ran!", result2.ReturnNameIndependentStruct(), "Functions should be able to run")

	assert.Same(result1, result2, "They must be the same!")
}

func TestResolve_ResolvesOneInstance_ObjectRegisteredAsSingletonAndItAppliesToDependency(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(containerBuilder, ss.NewS, o.As[ss.IIndependentStruct], o.AsSingleton)
			err = RegisterConstructor(containerBuilder, ss.NewStructRelyingOnIndependentStruct, o.As[ss.IStructRelyingOnIndependentStruct])
		},
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	container, err := containerBuilder.Build()
	assert.Nil(err)

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

	var result3 ss.IStructRelyingOnIndependentStruct
	assert.NotPanics(
		func() {
			result3, err = Resolve[ss.IStructRelyingOnIndependentStruct](container)
		},
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result3, "Second object should not be nil!")
	assert.Nil(err, "Second resolved object should not have any error!")
	assert.Contains(result3.ReturnSubStructName(), "Already Ran!", "Functions should be able to run")
}

func TestResolve_CannotResolve_ConstructorThrowsError(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(containerBuilder, ss.NewAReturningError, o.As[ss.IIndependentStruct])
		},
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")

	container, err := containerBuilder.Build()
	assert.Nil(err)

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
		`GoFac/pkg/GoFac.Resolve: Error resolving SampleStructs.IIndependentStruct!
	Inner error: GoFac/pkg/GoFac.runConstructor: Constructor of SampleStructs.IIndependentStruct threw an error
		Inner error: IndependentStruct: Error Forming IndependentStruct!`,
		err.Error(),
		"Error must show that constructor threw an error",
	)
}

func TestResolve_AbleToResolveInterfaceRelyingOnIndependentStruct(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(containerBuilder, ss.NewA, o.As[ss.IIndependentStruct])
			err = RegisterConstructor(
				containerBuilder,
				ss.NewStructRelyingOnIndependentStruct,
				o.As[ss.IStructRelyingOnIndependentStruct],
			)
		},
		"Should not have paniced when registering a constructor!",
	)

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	assert.Nil(err, "No Error should have happened when registering!"+errorMsg)

	container, err := containerBuilder.Build()
	assert.Nil(err)

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

func TestResolve_CannotResolveInterfaceRelyingOnIndependentStruct_DependencyNotRegistered(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(
				containerBuilder,
				ss.NewStructRelyingOnIndependentStruct,
				o.As[ss.IStructRelyingOnIndependentStruct],
			)
		},
		"Should not have paniced when registering a constructor!",
	)

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	assert.Nil(err, "No Error should have happened when registering!"+errorMsg)

	container, err := containerBuilder.Build()
	assert.Nil(err)

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
		`GoFac/pkg/GoFac.Resolve: Error resolving SampleStructs.IStructRelyingOnIndependentStruct!
	Inner error: GoFac/pkg/GoFac.getDependencies: Could not resolve SampleStructs.IStructRelyingOnIndependentStruct:
		Inner error: GoFac/pkg/GoFac.resolveOne: SampleStructs.IIndependentStruct is not registered!`,
		err.Error(),
		"Resolve must specify the cause of failure",
	)
}

func TestResolve_ResolvesStructWithSliceInputSuccessfully(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(
				containerBuilder,
				ss.NewA,
				o.As[ss.IIndependentStruct],
			)
			err = RegisterConstructor(
				containerBuilder,
				ss.NewB,
				o.As[ss.IIndependentStruct],
			)
			err = RegisterConstructor(
				containerBuilder,
				ss.NewStructRelyingOnIndependentStructs,
				o.As[ss.IStructRelyingOnIndependentStructs],
			)
		},
		"Should not have paniced when registering a constructor!",
	)

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	assert.Nil(err, "No Error should have happened when registering!"+errorMsg)

	container, err := containerBuilder.Build()
	assert.Nil(err)

	var result ss.IStructRelyingOnIndependentStructs
	assert.NotPanics(
		func() {
			result, err = Resolve[ss.IStructRelyingOnIndependentStructs](container)
		},
		"Should not have paniced when resolving interface!",
	)

	assert.NotNil(result, "Resolved object should not be nil!")
	assert.Nil(err, "Should not have any error!")
	assert.Equal("StructRelyingOnIndependentStructs", result.ReturnNameStructRelyingOnIndependentStruct(), "Names of the struct is different!")
	assert.Contains(result.ReturnSubStructNames(), "IndependentStruct", "IndependentStruct should have been resolved too!")
	assert.Contains(result.ReturnSubStructNames(), "IndependentStructB", "IndependentStructB should have been resolved too!")
}

func TestContainer_Resolve_ResolvesMultipleSuccessfully(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(
				containerBuilder,
				ss.NewA,
				o.As[ss.IIndependentStruct],
			)
			err = RegisterConstructor(
				containerBuilder,
				ss.NewB,
				o.As[ss.IIndependentStruct],
			)
		},
		"Should not have paniced when registering a constructor!",
	)

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	assert.Nil(err, "No Error should have happened when registering!"+errorMsg)

	container, err := containerBuilder.Build()
	assert.Nil(err)

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

func TestContainer_ResolveMultiple_ResolvesSingleton(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(
				containerBuilder,
				ss.NewB,
				o.As[ss.IIndependentStruct],
			)
			err = RegisterConstructor(
				containerBuilder,
				ss.NewA,
				o.As[ss.IIndependentStruct],
				o.AsSingleton,
			)
		},
		"Should not have paniced when registering a constructor!",
	)
	assert.Nil(err)

	container, err := containerBuilder.Build()
	assert.Nil(err)

	a1, err := Resolve[ss.IIndependentStruct](ctx.Background(), container)
	assert.NotNil(a1)
	assert.Nil(err)

	as, err := ResolveMultiple[ss.IIndependentStruct](ctx.Background(), container)
	assert.NotNil(as)
	assert.NotEmpty(as)
	assert.Nil(err)

	assert.Same(a1, as[1], "They must be the same!")
}

func TestResolve_ResolveSingletonObject_UnderMultithreading(t *testing.T) {
	NUM_WORKERS := 1000
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	if err := RegisterConstructor(
		containerBuilder,
		ss.NewIndependentStruct,
		o.AsSingleton,
		o.As[ss.IIndependentStruct],
	); err != nil {
		assert.Fail(err.Error())
	}
	container, err := containerBuilder.Build()
	assert.Nil(err)

	var wg sync.WaitGroup

	resolutionChannels := make(chan ss.IIndependentStruct, NUM_WORKERS)
	resolutionFunc := func() {
		defer wg.Done()
		if resolution, err := Resolve[ss.IIndependentStruct](ctx.Background(), container); err != nil {
			assert.Fail(err.Error())
		} else {
			resolutionChannels <- resolution
		}
	}

	for range NUM_WORKERS {
		wg.Add(1)
		go resolutionFunc()
	}

	wg.Wait()
	close(resolutionChannels)

	var resolutions []ss.IIndependentStruct
	for resolution := range resolutionChannels {
		resolutions = append(resolutions, resolution)
	}
	firstElem := resolutions[0]
	for i := 1; i < NUM_WORKERS; i++ {
		assert.Same(firstElem, resolutions[i])
	}
}

func TestResolve_CannotResolve_UnregisteredType(t *testing.T) {
	assert := assert.New(t)

	container, err := NewContainerBuilder().Build()
	assert.Nil(err)

	var result ss.IIndependentStruct
	result, err = Resolve[ss.IIndependentStruct](container)

	assert.Nil(result)
	assert.Error(err)
	assert.Equal(
		err.Error(),
		`GoFac/pkg/GoFac.Resolve: Error resolving SampleStructs.IIndependentStruct!
	Inner error: GoFac/pkg/GoFac.resolveOne: SampleStructs.IIndependentStruct is not registered!`,
	)
}

func TestResolveMultiple_ReturnsMultipleSingletons(t *testing.T) {
	assert := assert.New(t)

	cb := NewContainerBuilder()
	_ = RegisterConstructor(cb, ss.NewA, o.As[ss.IIndependentStruct], o.AsSingleton)
	_ = RegisterConstructor(cb, ss.NewB, o.As[ss.IIndependentStruct], o.AsSingleton)

	container, err := cb.Build()
	assert.Nil(err)
	slice, err := ResolveMultiple[ss.IIndependentStruct](container)

	assert.Nil(err)
	assert.Len(slice, 2)
	assert.NotNil(slice[0])
	assert.NotNil(slice[1])
	assert.NotSame(slice[0], slice[1], "They should be different singleton instances")
}
