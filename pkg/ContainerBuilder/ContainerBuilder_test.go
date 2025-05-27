package ContainerBuilder_test

import (
	ctx "context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/TaBSRest/GoFac"
	i "github.com/TaBSRest/GoFac/interfaces"
	cb "github.com/TaBSRest/GoFac/pkg/ContainerBuilder"
	RegistrationOptions "github.com/TaBSRest/GoFac/pkg/Options/Registration"
	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
)

func TestNew_DoesNotPanicAndReturnsContainerBuilder(t *testing.T) {
	assert := assert.New(t)

	var containerBuilder *cb.ContainerBuilder
	assert.NotPanics(func() {
		containerBuilder = cb.New()
	}, "Should not panic when constructing a ContainerBuilder",
	)
	assert.NotNil(containerBuilder, "Should return a non-nil ContainerBuilder")
}

func TestNew_ReturnsError_WhenBuildingTwice(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	_, err := containerBuilder.Build()
	assert.True(containerBuilder.IsBuilt())
	assert.Nil(err)

	_, err = containerBuilder.Build()
	assert.NotNil(err, "Should return error when building twice")
}

func TestRegisterConstructor_DoesNotPanic(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	var err error
	assert.NotPanics(func() {
		err = containerBuilder.Register(ss.NewA, RegistrationOptions.As[ss.IIndependentStruct])
	}, "Should not panic when registering a constructor",
	)
	assert.Nil(err)
}

func TestRegisterConstructor_ReturnsError_IfNameIsDuplicated(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	err1 := containerBuilder.Register(
		ss.NewA,
		RegistrationOptions.As[ss.IIndependentStruct],
		RegistrationOptions.Named("MyStruct"),
	)
	assert.Nil(err1)

	err2 := containerBuilder.Register(
		ss.NewA,
		RegistrationOptions.As[ss.IIndependentStruct],
		RegistrationOptions.Named("MyStruct"),
	)
	assert.NotNil(err2, "Expected error due to duplicate name, but got nil")
	if err2 != nil {
		assert.Contains(err2.Error(), "The name is already taken! If the registration is for a group, please use Options.Grouped!")
	}
}

func TestRegisterConstructor_AddsToGroupSuccessfully(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	err := containerBuilder.Register(
		ss.NewA,
		RegistrationOptions.As[ss.IIndependentStruct],
		RegistrationOptions.Grouped[ss.IIndependentStruct]("GroupA"),
	)
	assert.Nil(err)

	grouped, err := containerBuilder.GetGroupedRegistrations("GroupA")
	assert.Nil(err)
	assert.Len(grouped, 1)
}

func TestRegisterConstructor_ReturnsError_IfGroupTypeIsNotConsistent(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	err := containerBuilder.Register(
		ss.NewA,
		RegistrationOptions.As[ss.IIndependentStruct],
		RegistrationOptions.Grouped[ss.IIndependentStruct]("GroupA"),
	)
	assert.Nil(err)

	err = containerBuilder.Register(
		ss.NewB,
		RegistrationOptions.As[ss.IIndependentStruct2],
		RegistrationOptions.Grouped[ss.IIndependentStruct2]("GroupA"),
	)
	assert.NotNil(err)
	assert.Contains(err.Error(), "The type of the group must be consistent for all group members!")
}

func TestRegisterConstructor_AddsRegistrationToPerContextRegistrations_WhenScopeIsPerContext(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	err := containerBuilder.Register(
		ss.NewA,
		RegistrationOptions.As[ss.IIndependentStruct],
		RegistrationOptions.PerContext,
	)
	assert.Nil(err)

	perContextRegistrations := containerBuilder.GetPerContextRegistrations()
	assert.Len(perContextRegistrations, 1,
		"One registration should be in perContextRegistrations",
	)
}

func TestGetRegistrations_ReturnedValuesAreImmutable(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	containerBuilder.Register(ss.NewIndependentStruct, RegistrationOptions.As[ss.IIndependentStruct])
	containerBuilder.Register(ss.NewA)
	regs, found := cb.GetRegistrations[ss.IIndependentStruct](containerBuilder)
	assert.True(found)
	assert.Equal(2, len(regs))

	regs = regs[:1]

	newCopy, found := cb.GetRegistrations[ss.IIndependentStruct](containerBuilder)
	assert.True(found)
	assert.Equal(2, len(newCopy))
}

func TestGetRegistrationsFor_ReturnsExpectedRegistrations(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	err := containerBuilder.Register(
		ss.NewA,
		RegistrationOptions.As[ss.IIndependentStruct],
	)
	assert.Nil(err)

	typ := reflect.TypeFor[ss.IIndependentStruct]()
	registrations, found := containerBuilder.GetRegistrationsFor(typ)
	assert.True(found)
	assert.Len(registrations, 1)
}

func TestBuild_ReturnsContainer(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	container, err := containerBuilder.Build()

	assert.Nil(err)
	assert.NotNil(container)
	assert.True(containerBuilder.IsBuilt())

	self, err := GoFac.Resolve[i.Container](ctx.Background(), container)
	assert.Nil(err)
	assert.Equal(container, self)
}

func TestBuild_AbleToResolveContainerAndTheDependencyInIt(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	containerBuilder.Register(ss.NewIndependentStruct)
	containerBuilder.Register(
		func(container i.Container) (ss.IIndependentStruct, error) {
			return GoFac.Resolve[*ss.IndependentStruct](ctx.Background(), container)
		},
		RegistrationOptions.As[ss.IIndependentStruct],
	)
	container, err := containerBuilder.Build()
	assert.True(containerBuilder.IsBuilt())
	if err != nil {
		assert.Fail(err.Error())
	}

	sample1, err := GoFac.Resolve[*ss.IndependentStruct](ctx.Background(), container)
	assert.NotNil(sample1)
	if err != nil {
		assert.Fail(err.Error())
	}

	sample2, err := GoFac.Resolve[ss.IIndependentStruct](ctx.Background(), container)
	assert.NotNil(sample2)
	if err != nil {
		assert.Fail(err.Error())
	}
}

func TestBuild_SetsBuildOptionInContainer(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	container, err := containerBuilder.Build()

	assert.Nil(err)
	assert.NotNil(container)
	assert.NotNil(container.BuildOption)
}
