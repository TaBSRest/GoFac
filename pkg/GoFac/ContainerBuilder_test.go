package GoFac_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	gf "github.com/TaBSRest/GoFac/pkg/GoFac"
	"github.com/TaBSRest/GoFac/pkg/GoFac/Options"
	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
)

func TestContainer_Constructor_InitializedProperly(t *testing.T) {
	assert := assert.New(t)

	var gofac *gf.ContainerBuilder
	assert.NotPanics(
		func() {
			gofac = gf.NewContainerBuilder()
		},
		"Should not panic when creating new container",
	)
	assert.NotNil(gofac, "Initialized container should not be nil")
}

func TestNewContainerBuilder_DoesNotBuildTwice(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := gf.NewContainerBuilder()
	_, err := containerBuilder.Build()
	assert.True(containerBuilder.IsBuilt())
	assert.Nil(err)

	_, err = containerBuilder.Build()
	assert.NotNil(err)
}

func TestBuild_AbleToResolveContainerAndTheDependencyInIt(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := gf.NewContainerBuilder()
	gf.RegisterConstructor(containerBuilder, ss.NewIndependentStruct)
	gf.RegisterConstructor(
		containerBuilder,
		func(container *gf.Container) (ss.IIndependentStruct, error) {
			return gf.Resolve[*ss.IndependentStruct](container)
		},
		Options.As[ss.IIndependentStruct],
	)
	container, err := containerBuilder.Build()
	assert.True(containerBuilder.IsBuilt())
	if err != nil {
		assert.Fail(err.Error())
	}

	sample1, err := gf.Resolve[*ss.IndependentStruct](container)
	assert.NotNil(sample1)
	if err != nil {
		assert.Fail(err.Error())
	}

	sample2, err := gf.Resolve[ss.IIndependentStruct](container)
	assert.NotNil(sample2)
	if err != nil {
		assert.Fail(err.Error())
	}
}

func TestGetRegistrations_ReturnedValuesAreImmutable(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := gf.NewContainerBuilder()
	gf.RegisterConstructor(containerBuilder, ss.NewIndependentStruct, Options.As[ss.IIndependentStruct])
	gf.RegisterConstructor(containerBuilder, ss.NewA)
	regs, found := gf.GetRegistrations[ss.IIndependentStruct](containerBuilder)
	assert.True(found)
	assert.Equal(2, len(regs))
	regs = regs[:1]

	newCopy, found := gf.GetRegistrations[ss.IIndependentStruct](containerBuilder)
	assert.True(found)
	assert.Equal(2, len(newCopy))
}

