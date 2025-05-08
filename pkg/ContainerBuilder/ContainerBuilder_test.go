package ContainerBuilder_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	i "github.com/TaBSRest/GoFac/interfaces"
	cb "github.com/TaBSRest/GoFac/pkg/ContainerBuilder"
	"github.com/TaBSRest/GoFac"
	"github.com/TaBSRest/GoFac/pkg/Options"
	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
)

func TestContainer_Constructor_InitializedProperly(t *testing.T) {
	assert := assert.New(t)

	var gofac *cb.ContainerBuilder
	assert.NotPanics(
		func() {
			gofac = cb.New()
		},
		"Should not panic when creating new container",
	)
	assert.NotNil(gofac, "Initialized container should not be nil")
}

func TestNewContainerBuilder_DoesNotBuildTwice(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	_, err := containerBuilder.Build()
	assert.True(containerBuilder.IsBuilt())
	assert.Nil(err)

	_, err = containerBuilder.Build()
	assert.NotNil(err)
}

func TestBuild_AbleToResolveContainerAndTheDependencyInIt(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	containerBuilder.Register(ss.NewIndependentStruct)
	containerBuilder.Register(
		func(container i.Container) (ss.IIndependentStruct, error) {
			return GoFac.Resolve[*ss.IndependentStruct](container)
		},
		Options.As[ss.IIndependentStruct],
	)
	container, err := containerBuilder.Build()
	assert.True(containerBuilder.IsBuilt())
	if err != nil {
		assert.Fail(err.Error())
	}

	sample1, err := GoFac.Resolve[*ss.IndependentStruct](container)
	assert.NotNil(sample1)
	if err != nil {
		assert.Fail(err.Error())
	}

	sample2, err := GoFac.Resolve[ss.IIndependentStruct](container)
	assert.NotNil(sample2)
	if err != nil {
		assert.Fail(err.Error())
	}
}

func TestGetRegistrations_ReturnedValuesAreImmutable(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := cb.New()
	containerBuilder.Register(ss.NewIndependentStruct, Options.As[ss.IIndependentStruct])
	containerBuilder.Register(ss.NewA)
	regs, found := cb.GetRegistrations[ss.IIndependentStruct](containerBuilder)
	assert.True(found)
	assert.Equal(2, len(regs))
	regs = regs[:1]

	newCopy, found := cb.GetRegistrations[ss.IIndependentStruct](containerBuilder)
	assert.True(found)
	assert.Equal(2, len(newCopy))
}

